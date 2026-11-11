package v6

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/poll"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
)

// applyWorkloads creates or updates every Workload of the manifest with the
// bindings it declares, names resolved to the ids this run has. Bindings to
// Environments the run does not touch (others the manifest has not, or all but
// the --environment one) are kept as they are: a deploy to Staging must not
// change what Production serves.
func (s *DeployStateV6) applyWorkloads(ctx context.Context) error {
	s.Workloads = nil
	for _, wm := range s.Manifest.Workloads {
		ours, err := s.bindings(wm)
		if err != nil {
			return err
		}
		if len(ours) == 0 {
			continue
		}

		if err := s.waitBindable(ctx, ours); err != nil {
			return err
		}

		existing, err := s.findWorkload(ctx, wm.Name)
		if err != nil {
			return err
		}
		req, err := workloadRequest(existing, wm, ours)
		if err != nil {
			return err
		}

		var saved sdk.Workload
		if existing != nil {
			saved, err = s.Clients.Workloads.Update(ctx, existing.GetId(), req)
		} else {
			saved, err = s.Clients.Workloads.Create(ctx, req)
		}
		if err != nil {
			return err
		}

		// Record the id before waiting, so a run that stops here finds the
		// Workload again instead of creating another.
		record := func(w sdk.Workload) {
			upsertNamed(&s.Conf.Workloads, contracts.AzionJsonV6Workload{ID: w.GetId(), Name: wm.Name, Url: workloadURL(w, ours)},
				func(w contracts.AzionJsonV6Workload) string { return w.Name })
			s.writeConf()
		}
		record(saved)
		if saved, err = s.waitWorkload(ctx, saved); err != nil {
			return err
		}
		// The workload domain may only be known once the version is built.
		record(saved)

		s.Workloads = append(s.Workloads, appliedWorkload{Workload: saved, Bindings: ours})
		for _, b := range ours {
			s.say(fmt.Sprintf(msg.V6WorkloadSaved, saved.GetName(), saved.GetId(), s.deploymentName(b.DeploymentId), s.environmentName(b.EnvironmentId)))
		}
	}
	return nil
}

// waitWorkload waits until the version a save of the Workload produced is
// built: its domains serve the new bindings only once it is ready. Without
// this wait the deploy reported success while the domain still answered 404.
func (s *DeployStateV6) waitWorkload(ctx context.Context, w sdk.Workload) (sdk.Workload, error) {
	latest := w
	err := s.waitVersion(ctx, "Workload", w.GetName(), strconv.FormatInt(w.GetId(), 10), w.GetVersionId(),
		version.State(w.GetVersionState()), msg.ErrorV6WorkloadFailed,
		func(ctx context.Context) (version.State, error) {
			cur, err := s.Clients.Workloads.Get(ctx, w.GetId())
			if err == nil {
				latest = cur
			}
			return version.State(cur.GetVersionState()), err
		})
	return latest, err
}

// waitVersion waits for the build of a resource the API builds on its own
// after a save (Environments and Workloads have no build call), until its
// version is ready. failed is the error a build that does not reach ready
// wraps.
//
// A resource that reports no version state is not waited on: there is nothing
// to read the build from. The API is not consistent about case, so neither is
// this.
func (s *DeployStateV6) waitVersion(ctx context.Context, kind, name, id, versionID string, state version.State,
	failed error, read func(context.Context) (version.State, error),
) error {
	norm := func(st version.State) version.State { return version.State(strings.ToLower(string(st))) }
	last := norm(state)
	if last == "" || last == version.Ready {
		return nil
	}
	s.say(fmt.Sprintf(msg.V6WaitingForVersion, kind, name, id, last))

	start := time.Now()
	_, err := poll.Until(ctx, s.pollOptions(), read, func(st version.State) (bool, error) {
		st = norm(st)
		if st == "" {
			return true, nil
		}
		if st != last {
			s.say(fmt.Sprintf(msg.V6BuildTransition, strings.ToLower(kind), name, versionID, last, st, time.Since(start).Round(time.Second)))
			last = st
		}
		return version.BuildDone(st)
	})
	s.recordWait(time.Since(start))

	switch {
	case err == nil:
		s.say(fmt.Sprintf(msg.V6VersionReady, kind, name, id, time.Since(start).Round(time.Second)))
		return nil
	case ctx.Err() != nil:
		return ctx.Err()
	case errors.Is(err, version.ErrBuildFailed):
		return fmt.Errorf("%s %q (%s) is %s: %w", kind, name, id, last, failed)
	case errors.Is(err, poll.ErrTimeout):
		return fmt.Errorf("%s %q (%s) is still %s: %w: %w", kind, name, id, last, failed, err)
	default:
		return fmt.Errorf("could not read %s %q (%s) while it builds: %w", kind, name, id, err)
	}
}

// waitBindable waits until every Environment and Deployment the bindings name
// can be bound. The API accepts a binding only to one whose version is queued,
// building or ready, and a Deployment leaves draft only once a release of it
// is activated, so each is read again here, right before the Workload is saved.
func (s *DeployStateV6) waitBindable(ctx context.Context, bindings []sdk.WorkloadBindingRequest) error {
	seenEnv, seenDep := map[string]bool{}, map[string]bool{}
	for _, b := range bindings {
		if id := b.EnvironmentId; !seenEnv[id] {
			seenEnv[id] = true
			if err := s.waitUntilBindable(ctx, "Environment", s.environmentName(id), id, func(ctx context.Context) (version.State, error) {
				e, err := s.Clients.Environments.Get(ctx, id)
				return e.VersionState, err
			}); err != nil {
				return err
			}
		}
		if id := b.DeploymentId; !seenDep[id] {
			seenDep[id] = true
			if err := s.waitUntilBindable(ctx, "Deployment", s.deploymentName(id), id, func(ctx context.Context) (version.State, error) {
				d, err := s.Clients.Deployments.Get(ctx, id)
				return d.VersionState, err
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// waitUntilBindable polls one Environment or Deployment until it is bindable.
// It says so only when it has to wait, and fails when the version failed.
func (s *DeployStateV6) waitUntilBindable(ctx context.Context, kind, name, id string, read func(context.Context) (version.State, error)) error {
	told := false
	state, err := poll.Until(ctx, s.pollOptions(), read, func(st version.State) (bool, error) {
		if st.Bindable() {
			return true, nil
		}
		if st.Failed() {
			return true, fmt.Errorf("%s %q (%s) is %s: %w", kind, name, id, st, msg.ErrorV6NotBindable)
		}
		if !told {
			told = true
			s.say(fmt.Sprintf(msg.V6WaitingToBind, kind, name, id, st))
		}
		return false, nil
	})
	switch {
	case err == nil, errors.Is(err, msg.ErrorV6NotBindable):
		return err
	case errors.Is(err, poll.ErrTimeout):
		return fmt.Errorf("%s %q (%s) is still %s, and a Workload can be bound only to one that is queued, building or ready: %w", kind, name, id, state, err)
	default:
		return fmt.Errorf("could not read %s %q (%s) before binding it to the Workload: %w", kind, name, id, err)
	}
}

// appliedWorkload is a Workload this run saved, with the bindings it set.
type appliedWorkload struct {
	Workload sdk.Workload
	Bindings []sdk.WorkloadBindingRequest
}

// bindings resolves a Workload's manifest bindings that are in scope.
func (s *DeployStateV6) bindings(wm contracts.ManifestV6Workload) ([]sdk.WorkloadBindingRequest, error) {
	var out []sdk.WorkloadBindingRequest
	for _, b := range wm.Bindings {
		if !s.inScope(b.EnvironmentName) {
			continue
		}
		env, ok := s.Environments[b.EnvironmentName]
		if !ok {
			return nil, fmt.Errorf("workload %q: environment %q was not resolved", wm.Name, b.EnvironmentName)
		}
		t := s.target(b.DeploymentName)
		if t == nil {
			return nil, fmt.Errorf("workload %q: deployment %q was not resolved", wm.Name, b.DeploymentName)
		}
		req := sdk.WorkloadBindingRequest{EnvironmentId: env.ID, DeploymentId: t.Deployment.ID}
		for _, d := range b.Domains {
			entry := sdk.PatchedWorkloadBindingRequestDomainsInner{Domain: d.Domain}
			if d.Certificate != nil {
				entry.SetCertificate(*d.Certificate)
			}
			req.Domains = append(req.Domains, entry)
		}
		out = append(out, req)
	}
	return out, nil
}

func (s *DeployStateV6) findWorkload(ctx context.Context, name string) (*sdk.Workload, error) {
	for _, rec := range s.Conf.Workloads {
		if rec.Name != name || rec.ID == 0 {
			continue
		}
		w, err := s.Clients.Workloads.Get(ctx, rec.ID)
		if err == nil {
			return &w, nil
		}
		if !errors.Is(err, utils.ErrorNotFound404) {
			return nil, err
		}
	}
	w, found, err := s.Clients.Workloads.FindByName(ctx, name)
	if err != nil || !found {
		return nil, err
	}
	return &w, nil
}

// workloadRequest builds the full request for a Workload. An update starts
// from what exists, because the API replaces the whole Workload, then applies
// what the manifest sets and swaps in this run's bindings, keeping the ones for
// other Environments.
func workloadRequest(existing *sdk.Workload, wm contracts.ManifestV6Workload, ours []sdk.WorkloadBindingRequest) (sdk.WorkloadRequest, error) {
	req := sdk.WorkloadRequest{}
	touched := map[string]bool{}
	for _, b := range ours {
		touched[b.EnvironmentId] = true
	}
	var kept []sdk.WorkloadBindingRequest
	if existing != nil {
		// The response and request types differ (TLSWorkload against
		// TLSWorkloadRequest, and so on) but share their JSON, so the
		// current Workload becomes the base of the request through it.
		data, err := json.Marshal(existing)
		if err != nil {
			return req, err
		}
		if err := json.Unmarshal(data, &req); err != nil {
			return req, err
		}
		for _, b := range req.Bindings {
			if !touched[b.EnvironmentId] {
				kept = append(kept, b)
			}
		}
	}
	req.Name = wm.Name
	if wm.Infrastructure != nil {
		req.Infrastructure = wm.Infrastructure
	}
	if wm.WorkloadDomainAllowAccess != nil {
		req.WorkloadDomainAllowAccess = wm.WorkloadDomainAllowAccess
	}
	if wm.Tls != nil {
		req.Tls = wm.Tls
	}
	if wm.Protocols != nil {
		req.Protocols = wm.Protocols
	}
	if wm.Mtls != nil {
		req.Mtls = wm.Mtls
	}
	req.Bindings = append(kept, ours...)
	return req, nil
}

// workloadURL is the address shown at the end of the deploy: the first domain
// of this run's bindings, or the Workload's own azion domain.
func workloadURL(w sdk.Workload, ours []sdk.WorkloadBindingRequest) string {
	for _, b := range ours {
		if len(b.Domains) > 0 {
			return "https://" + b.Domains[0].Domain
		}
	}
	if d := w.GetWorkloadDomain(); d != "" {
		return "https://" + d
	}
	return ""
}

// hasBinding reports whether the Workload serves deploymentID in envID.
func hasBinding(w sdk.Workload, envID, deploymentID string) bool {
	for _, b := range w.GetBindings() {
		if b.EnvironmentId == envID && b.DeploymentId == deploymentID {
			return true
		}
	}
	return false
}

func (s *DeployStateV6) target(deploymentName string) *releaseTarget {
	for _, t := range s.Targets {
		if t.Name == deploymentName {
			return t
		}
	}
	return nil
}

func (s *DeployStateV6) deploymentName(id string) string {
	for _, t := range s.Targets {
		if t.Deployment.ID == id {
			return t.Name
		}
	}
	return id
}

func (s *DeployStateV6) environmentName(id string) string {
	for name, e := range s.Environments {
		if e.ID == id {
			return name
		}
	}
	return id
}
