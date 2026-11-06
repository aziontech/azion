package v6

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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

		s.Workloads = append(s.Workloads, appliedWorkload{Workload: saved, Bindings: ours})
		upsertNamed(&s.Conf.Workloads, contracts.AzionJsonV6Workload{ID: saved.GetId(), Name: wm.Name, Url: workloadURL(saved, ours)},
			func(w contracts.AzionJsonV6Workload) string { return w.Name })
		s.writeConf()
		for _, b := range ours {
			s.say(fmt.Sprintf(msg.V6WorkloadSaved, saved.GetName(), saved.GetId(), s.deploymentName(b.DeploymentId), s.environmentName(b.EnvironmentId)))
		}
	}
	return nil
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
	if wm.Active != nil {
		req.Active = wm.Active
	}
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
