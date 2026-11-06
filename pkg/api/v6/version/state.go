// Package version holds what every versioned resource of API v6 shares: the
// build states, how a wait reads them, the error a failed build turns into and
// the four operations each versioned client exposes.
//
// In v6 a save produces an immutable Draft, and nothing reaches traffic until a
// build takes that Draft through queued and building to ready. A step that
// writes configuration therefore does not finish when the API answers; it
// finishes when the build it started settles. This package is the vocabulary
// those steps share, so the seven versioned resources do not each redefine it.
package version

import (
	"fmt"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

// State is the build state the v6 API reports as version_state (resources)
// or state (releases).
type State string

const (
	Draft     State = "draft"
	Queued    State = "queued"
	Building  State = "building"
	Ready     State = "ready"
	Error     State = "error"
	Canceled  State = "canceled"
	Archiving State = "archiving"
	Archived  State = "archived"
	Deleting  State = "deleting"
	Deleted   State = "deleted"
)

// Known reports whether the CLI knows this state. A state the API adds later
// is unknown until this list learns it, and an unknown state is never treated
// as settled, so it can turn neither into success nor into failure by accident.
func (s State) Known() bool {
	switch s {
	case Draft, Queued, Building, Ready, Error, Canceled,
		Archiving, Archived, Deleting, Deleted:
		return true
	}
	return false
}

// Settled reports whether the state will not change without a new request.
func (s State) Settled() bool {
	switch s {
	case Draft, Ready, Error, Canceled, Archived, Deleted:
		return true
	}
	return false
}

// Failed reports whether a settled state means the build did not produce a
// usable version. Draft is settled but not failed: nothing was built yet.
func (s State) Failed() bool {
	switch s {
	case Error, Canceled, Archived, Deleted:
		return true
	}
	return false
}

// BuildDone is the done function for waiting on a build that this run
// started. It reports true once the build has settled, with ErrBuildFailed
// when it settled without a usable version.
//
// Draft does not end the wait even though it is settled: right after the build
// request the API may still report the version as draft until the build is
// queued, and if it never leaves draft the wait's timeout says so.
func BuildDone(s State) (bool, error) {
	switch {
	case s == Ready:
		return true, nil
	case s.Failed():
		return true, fmt.Errorf("%w: %s", ErrBuildFailed, s)
	case !s.Known():
		logger.Debug("Unknown version state, still waiting", zap.String("state", string(s)))
	}
	return false, nil
}
