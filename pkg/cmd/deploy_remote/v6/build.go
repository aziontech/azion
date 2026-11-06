package v6

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestv6 "github.com/aziontech/azion-cli/pkg/manifest/v6"
	"github.com/aziontech/azion-cli/pkg/poll"
	"go.uber.org/zap"
)

// buildConcurrency bounds how many builds of one level run at once. At one
// read per second each, it keeps a run far below the API's rate limit.
const buildConcurrency = 4

// BuildComment is recorded on every build the CLI starts.
const BuildComment = "azion-cli deploy"

// buildLevel builds every Draft of one level and waits for all of them.
//
// A failure does not cancel the other builds of the level: the step waits for
// all and returns every failure joined, so the user sees at once everything
// that went wrong. The builds that succeeded become Ready versions nothing
// refers to yet, which is harmless.
func (s *DeployStateV6) buildLevel(ctx context.Context, level []manifestv6.Pending) error {
	if len(level) == 0 {
		return nil
	}
	names := make([]string, len(level))
	for i, p := range level {
		names[i] = fmt.Sprintf("%s %q", p.Kind, p.Name)
	}
	s.say(fmt.Sprintf(msg.V6BuildsStarting, len(level), strings.Join(names, ", ")))

	sem := make(chan struct{}, buildConcurrency)
	states := make([]version.State, len(level))
	errs := make([]error, len(level))
	var wg sync.WaitGroup
	for i, p := range level {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				errs[i] = ctx.Err()
				return
			}
			defer func() { <-sem }()
			states[i], errs[i] = s.buildOne(ctx, p)
		}()
	}
	wg.Wait()

	for i, p := range level {
		if states[i] != "" {
			p.SetState(states[i])
		}
	}
	s.writeConf()

	if err := ctx.Err(); err != nil {
		return err
	}
	return errors.Join(errs...)
}

// buildOne starts one build and waits for it to settle. When the wait is
// interrupted, it asks the API to stop the build, so Ctrl-C does not leave
// work running that nobody will use.
func (s *DeployStateV6) buildOne(ctx context.Context, p manifestv6.Pending) (version.State, error) {
	fail := func(state version.State, err error) error {
		return &version.BuildError{Kind: p.Kind, Name: p.Name, ID: p.ID, VersionID: p.VersionID, State: state, Err: err}
	}

	if err := p.Client.Build(ctx, p.ID, p.VersionID, BuildComment); err != nil {
		return "", fail("", err)
	}

	start := time.Now()
	last := version.Draft
	state, err := poll.Until(ctx, s.pollOptions(),
		func(ctx context.Context) (version.State, error) {
			return p.Client.VersionState(ctx, p.ID, p.VersionID)
		},
		func(st version.State) (bool, error) {
			if st != last {
				s.say(fmt.Sprintf(msg.V6BuildTransition, p.Kind, p.Name, p.VersionID, last, st, time.Since(start).Round(time.Second)))
				last = st
			}
			return version.BuildDone(st)
		})
	s.recordWait(time.Since(start))

	if ctx.Err() != nil {
		cctx, cancel := cleanupContext(ctx)
		defer cancel()
		if cerr := p.Client.CancelBuild(cctx, p.ID, p.VersionID); cerr != nil {
			logger.Debug("Could not cancel an interrupted build", zap.String("kind", p.Kind), zap.String("name", p.Name), zap.Error(cerr))
		}
		return state, ctx.Err()
	}
	if err != nil {
		return state, fail(state, err)
	}
	s.say(fmt.Sprintf(msg.V6BuildSucceeded, p.Kind, p.Name, p.VersionID, time.Since(start).Round(time.Second)))
	return state, nil
}

// recordWait adds to the time spent waiting on builds, which the timing
// summary shows apart from the time spent in calls.
func (s *DeployStateV6) recordWait(d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.timing != nil {
		s.timing.BuildWait += d
	}
}
