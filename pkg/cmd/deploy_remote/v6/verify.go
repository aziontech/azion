package v6

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	apiDeployments "github.com/aziontech/azion-cli/pkg/api/v6/deployments"
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	manifestv6 "github.com/aziontech/azion-cli/pkg/manifest/v6"
)

// verifyDeploy is the last check of doc/V6-DEPLOY-PLAN.md section 6: what is
// live is what this run meant to put live. For every Deployment released:
//
//   - the release is ready and has the expected traffic role;
//   - on a single_version Deployment, it is the active release;
//   - each resource the release names (what it serves and its dependencies)
//     is in the hierarchy at the version this run built.
//
// And every Workload saved has the bindings this run set. Any other difference
// fails the deploy with the list of what differs.
func (s *DeployStateV6) verifyDeploy(ctx context.Context) error {
	var diffs []error

	for _, t := range s.Targets {
		r, err := s.getRelease(t)(ctx)
		if err != nil {
			return err
		}
		if r.State != version.Ready {
			diffs = append(diffs, fmt.Errorf("deployment %q: release %s is %s, not ready", t.Name, r.ID, r.State))
		}
		want := expectedRole(t.Deployment)
		if r.TrafficRole != want {
			diffs = append(diffs, fmt.Errorf("deployment %q: release %s has traffic role %s, not %s", t.Name, r.ID, r.TrafficRole, want))
		}

		h, err := s.Clients.Deployments.Hierarchy(ctx, t.Deployment.ID)
		if err != nil {
			return err
		}
		if want == apiDeployments.RoleActive && h.ActiveReleaseID != r.ID {
			diffs = append(diffs, fmt.Errorf("deployment %q: the active release is %q, not %s", t.Name, h.ActiveReleaseID, r.ID))
		}
		diffs = append(diffs, s.compareHierarchy(t.Name, h)...)
	}

	for _, aw := range s.Workloads {
		w, err := s.Clients.Workloads.Get(ctx, aw.Workload.GetId())
		if err != nil {
			return err
		}
		if st := version.State(strings.ToLower(w.GetVersionState())); st != "" && st != version.Ready {
			diffs = append(diffs, fmt.Errorf("workload %q: its version %s is %s, not ready, so its domains do not serve it yet", w.GetName(), w.GetVersionId(), st))
		}
		for _, b := range aw.Bindings {
			if !hasBinding(w, b.EnvironmentId, b.DeploymentId) {
				diffs = append(diffs, fmt.Errorf("workload %q has no binding to Deployment %s in Environment %s", w.GetName(), b.DeploymentId, b.EnvironmentId))
			}
		}
	}

	if len(diffs) > 0 {
		return fmt.Errorf("%w:\n%w", msg.ErrorV6VerifyFailed, errors.Join(diffs...))
	}
	return nil
}

// compareHierarchy checks each resource the release names against what one
// Deployment serves: missing, or at a version other than the one this run
// built, is a difference.
func (s *DeployStateV6) compareHierarchy(deployment string, h apiDeployments.Hierarchy) []error {
	served := make(map[string]string, len(h.Resources))
	for _, e := range h.Resources {
		served[e.ResourceType+"/"+e.ResourceID] = e.Version
	}

	var diffs []error
	for _, level := range [][]manifestv6.Pending{s.Publishable, s.Dependencies} {
		for _, p := range level {
			key := resourceTypes[p.Kind] + "/" + strconv.FormatInt(p.ID, 10)
			got, ok := served[key]
			switch {
			case !ok:
				diffs = append(diffs, fmt.Errorf("deployment %q: %s %q (id %d) is not served", deployment, p.Kind, p.Name, p.ID))
			case got != p.VersionID:
				diffs = append(diffs, fmt.Errorf("deployment %q: %s %q is served at version %s, not %s", deployment, p.Kind, p.Name, got, p.VersionID))
			}
		}
	}
	return diffs
}
