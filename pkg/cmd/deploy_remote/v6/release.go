package v6

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	apiDeployments "github.com/aziontech/azion-cli/pkg/api/v6/deployments"
	apiEnvironments "github.com/aziontech/azion-cli/pkg/api/v6/environments"
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestv6 "github.com/aziontech/azion-cli/pkg/manifest/v6"
	"github.com/aziontech/azion-cli/pkg/poll"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// resourceTypes maps a resource kind to the resource_type a release and the
// resource hierarchy use. The values are the hierarchy's keys; the release's
// accepted values are also still to be confirmed by the API team.
var resourceTypes = map[string]string{
	"application":  "application",
	"firewall":     "firewall",
	"custom page":  "custom_page",
	"connector":    "connector",
	"network list": "network_list",
	"waf":          "waf",
	"function":     "function",
}

// releaseTarget is one Deployment this run releases to.
type releaseTarget struct {
	Name       string
	Deployment apiDeployments.Deployment
	Release    apiDeployments.Release
	Activated  bool
}

// withDefaults fills what a manifest without Environments and Deployments
// leaves out, so such a project still deploys: the Environment named by
// --environment, or Production, which must already exist, and one
// Deployment named after the project, bound by every Workload that has no
// binding of its own.
//
// That Deployment is FLEXIBLE, like the ones the bundler writes. STRICT would
// suit the deploy better: it keeps every release to the same resources and
// lets only their version_id change, which is what a CLI deploy does. But the
// API refuses to build any release after the first one on a STRICT Deployment
// (409 43001, "The field 'resources' is immutable after creation", with the
// application read as empty although the release has it), so a STRICT default
// could deploy only once. Go back to STRICT once the API fixes that check.
// FLEXIBLE does not validate the release's resources at all.
func withDefaults(m *contracts.ManifestV6, envName, projectName string) (mustExist map[string]bool) {
	mustExist = map[string]bool{}
	if len(m.Environments) == 0 {
		m.Environments = []sdk.EnvironmentCreateRequest{{Name: envName}}
		mustExist[envName] = true
	}
	if len(m.Deployments) == 0 {
		m.Deployments = []sdk.DeploymentCreateRequest{{Name: projectName, BindingPolicy: defaultBindingPolicy}}
		for i := range m.Workloads {
			if len(m.Workloads[i].Bindings) == 0 {
				m.Workloads[i].Bindings = []contracts.ManifestV6Binding{{EnvironmentName: m.Environments[0].Name, DeploymentName: projectName}}
			}
		}
	}
	return mustExist
}

// inScope reports whether a binding's Environment is part of this run: all of
// them, or the one --environment names.
func (s *DeployStateV6) inScope(envName string) bool {
	return s.Opts.Environment == "" || strings.EqualFold(envName, s.Opts.Environment)
}

// applyEnvironments makes every Environment the manifest declares exist, and
// records their ids. An Environment is reused by name; one that does not exist
// is created. Its deployment policy cannot change, so a different one in the
// manifest stops the deploy instead of being ignored.
func (s *DeployStateV6) applyEnvironments(ctx context.Context) error {
	s.Environments = map[string]apiEnvironments.Environment{}
	for _, want := range s.Manifest.Environments {
		if !s.inScope(want.Name) && !s.boundToTarget(want.Name) {
			continue
		}
		env, found, err := s.findEnvironment(ctx, want.Name)
		if err != nil {
			return err
		}
		switch {
		case found && want.DeploymentPolicy != "" && env.DeploymentPolicy != want.DeploymentPolicy:
			return fmt.Errorf("environment %q is %s and the manifest asks for %s: %w", env.Name, env.DeploymentPolicy, want.DeploymentPolicy, msg.ErrorV6PolicyImmutable)
		case !found && s.mustExist[want.Name]:
			return fmt.Errorf("%q: %w", want.Name, msg.ErrorV6EnvironmentNotFound)
		case !found:
			env, err = s.Clients.Environments.Create(ctx, want)
			if err != nil {
				return err
			}
			s.say(fmt.Sprintf(msg.V6EnvironmentCreated, env.Name, env.ID))
		case s.inScope(want.Name) && !s.mustExist[want.Name]:
			// The manifest is the Environment's configuration: write it on
			// every deploy, or an Environment created earlier keeps what it
			// was created with.
			env, err = s.Clients.Environments.Update(ctx, env.ID, environmentUpdate(want))
			if err != nil {
				return err
			}
			s.say(fmt.Sprintf(msg.V6EnvironmentUpdated, env.Name, env.ID))
		}
		// Creating and updating both save a version that the API builds on
		// its own; the Environment has its new configuration once it is ready.
		if err := s.waitVersion(ctx, "Environment", env.Name, env.ID, env.VersionID, env.VersionState, msg.ErrorV6EnvironmentFailed,
			func(ctx context.Context) (version.State, error) {
				cur, err := s.Clients.Environments.Get(ctx, env.ID)
				if err == nil {
					env = cur
				}
				return cur.VersionState, err
			}); err != nil {
			return err
		}
		s.Environments[want.Name] = env
		s.say(fmt.Sprintf(msg.V6EnvironmentResolved, env.Name, env.ID, env.DeploymentPolicy))
		upsertNamed(&s.Conf.Environments, contracts.AzionJsonV6Environment{ID: env.ID, Name: want.Name, Policy: env.DeploymentPolicy},
			func(e contracts.AzionJsonV6Environment) string { return e.Name })
	}
	s.writeConf()
	return nil
}

// environmentUpdate is the manifest's configuration of an Environment as an
// update. The name is left out, since the Environment may have been found under
// another capitalisation of it, and so is the deployment policy, which cannot
// change (applyEnvironments refuses a manifest that asks for another one).
func environmentUpdate(want sdk.EnvironmentCreateRequest) sdk.PatchedEnvironmentRequest {
	return sdk.PatchedEnvironmentRequest{
		Description:    want.Description,
		LogVerbosity:   want.LogVerbosity,
		RobotsPolicy:   want.RobotsPolicy,
		Protection:     want.Protection,
		BranchTracking: want.BranchTracking,
	}
}

// findEnvironment looks an Environment up by the id azion.json recorded for
// the name, then by the exact name, then ignoring case, so a manifest's
// "production" finds the account's "Production".
func (s *DeployStateV6) findEnvironment(ctx context.Context, name string) (apiEnvironments.Environment, bool, error) {
	for _, rec := range s.Conf.Environments {
		if rec.Name == name && rec.ID != "" {
			env, err := s.Clients.Environments.Get(ctx, rec.ID)
			if err == nil {
				return env, true, nil
			}
			if !errors.Is(err, utils.ErrorNotFound404) {
				return env, false, err
			}
		}
	}
	env, found, err := s.Clients.Environments.FindByName(ctx, name)
	if err != nil || found {
		return env, found, err
	}
	for _, alt := range []string{strings.ToLower(name), strings.ToUpper(name[:1]) + strings.ToLower(name[1:])} {
		if alt == name {
			continue
		}
		if env, found, err = s.Clients.Environments.FindByName(ctx, alt); err != nil || found {
			return env, found, err
		}
	}
	return apiEnvironments.Environment{}, false, nil
}

// targetDeployments is the names of the Deployments this run releases: every
// one the manifest declares, or, with --environment, those bound to it.
func (s *DeployStateV6) targetDeployments() []string {
	var out []string
	for _, d := range s.Manifest.Deployments {
		if s.Opts.Environment == "" || s.deploymentInScope(d.Name) {
			out = append(out, d.Name)
		}
	}
	return out
}

func (s *DeployStateV6) deploymentInScope(name string) bool {
	for _, w := range s.Manifest.Workloads {
		for _, b := range w.Bindings {
			if b.DeploymentName == name && s.inScope(b.EnvironmentName) {
				return true
			}
		}
	}
	return false
}

// boundToTarget reports whether an Environment is bound to a Deployment in
// scope, so its id is needed for the Workload.
func (s *DeployStateV6) boundToTarget(envName string) bool {
	for _, w := range s.Manifest.Workloads {
		for _, b := range w.Bindings {
			if b.EnvironmentName == envName && s.inScope(b.EnvironmentName) {
				return true
			}
		}
	}
	return false
}

// applyDeployments makes every Deployment this run releases exist. One that
// does not exist is created with the manifest's binding policy and, unless the
// manifest sets one, the deployment policy of the Environment it is bound to.
func (s *DeployStateV6) applyDeployments(ctx context.Context) error {
	s.Targets = nil
	for _, name := range s.targetDeployments() {
		want := s.manifestDeployment(name)
		d, err := s.ensureDeployment(ctx, want)
		if err != nil {
			return err
		}
		s.Targets = append(s.Targets, &releaseTarget{Name: name, Deployment: d})
		upsertNamed(&s.Conf.Deployments, contracts.AzionJsonV6Deployment{ID: d.ID, Name: name, Policy: d.DeploymentPolicy, Release: s.deploymentRecord(name).Release},
			func(d contracts.AzionJsonV6Deployment) string { return d.Name })
		s.writeConf()
	}
	return nil
}

// defaultBindingPolicy is the binding policy of a Deployment the manifest does
// not declare. See withDefaults.
const defaultBindingPolicy = apiDeployments.BindingFlexible

func (s *DeployStateV6) manifestDeployment(name string) sdk.DeploymentCreateRequest {
	for _, d := range s.Manifest.Deployments {
		if d.Name == name {
			return d
		}
	}
	return sdk.DeploymentCreateRequest{Name: name, BindingPolicy: defaultBindingPolicy}
}

func (s *DeployStateV6) ensureDeployment(ctx context.Context, want sdk.DeploymentCreateRequest) (apiDeployments.Deployment, error) {
	if rec := s.deploymentRecord(want.Name); rec.ID != "" {
		d, err := s.Clients.Deployments.Get(ctx, rec.ID)
		switch {
		case err == nil:
			return d, nil
		case errors.Is(err, utils.ErrorNotFound404):
			s.say(fmt.Sprintf(msg.V6DeploymentRecreate, rec.ID))
		default:
			return d, err
		}
	}
	d, found, err := s.Clients.Deployments.FindByName(ctx, want.Name)
	if err != nil {
		return d, err
	}
	if !found {
		if want.DeploymentPolicy == nil {
			if policy := s.boundEnvironmentPolicy(want.Name); policy != "" {
				want.SetDeploymentPolicy(policy)
			}
		}
		if d, err = s.Clients.Deployments.Create(ctx, want); err != nil {
			return d, err
		}
		s.say(fmt.Sprintf(msg.V6DeploymentCreated, d.Name, d.ID))
	}

	// A new Deployment may still be settling.
	if d.VersionState != "" && !d.VersionState.Settled() {
		d, err = poll.Until(ctx, s.pollOptions(),
			func(ctx context.Context) (apiDeployments.Deployment, error) {
				return s.Clients.Deployments.Get(ctx, d.ID)
			},
			func(d apiDeployments.Deployment) (bool, error) {
				if d.VersionState.Failed() {
					return true, fmt.Errorf("deployment %q is %s: %w", d.Name, d.VersionState, msg.ErrorV6DeploymentFailed)
				}
				return d.VersionState.Settled(), nil
			})
	}
	return d, err
}

// boundEnvironmentPolicy is the deployment policy of the first Environment a
// Deployment is bound to.
func (s *DeployStateV6) boundEnvironmentPolicy(deployment string) string {
	for _, w := range s.Manifest.Workloads {
		for _, b := range w.Bindings {
			if b.DeploymentName == deployment {
				if env, ok := s.Environments[b.EnvironmentName]; ok {
					return env.DeploymentPolicy
				}
			}
		}
	}
	return ""
}

// releaseCreate opens a release Draft on every target Deployment, naming the
// resources it serves.
func (s *DeployStateV6) releaseCreate(ctx context.Context) error {
	resources := s.releaseResources()
	for _, t := range s.Targets {
		r, err := s.Clients.Deployments.CreateRelease(ctx, t.Deployment.ID, resources, apiDeployments.OriginCLI)
		if err != nil {
			return err
		}
		t.Release = r
		s.say(fmt.Sprintf(msg.V6ReleaseCreated, r.ID, t.Name))
		s.deploymentRecord(t.Name).Release = contracts.AzionJsonV6Release{ID: r.ID, State: string(r.State)}
		s.writeConf()
	}
	return nil
}

// releaseResources lists what a release serves: the Application, Firewall and
// Custom Page of this run and the dependencies they use (Connectors,
// Functions, WAFs, Network Lists), each at the version this run built (the
// version id the API returned, like A32XCE2E). The build steps only finish
// once those versions are ready.
//
// The API does not fill in dependencies: a release without one its
// Application uses fails to build with "missing required subresource
// dependency". A version does not list what it depends on either, so the
// release names every dependency the manifest declares.
func (s *DeployStateV6) releaseResources() []apiDeployments.ReleaseResource {
	var out []apiDeployments.ReleaseResource
	for _, level := range [][]manifestv6.Pending{s.Publishable, s.Dependencies} {
		for _, p := range level {
			out = append(out, apiDeployments.ReleaseResource{ResourceID: p.ID, ResourceType: resourceTypes[p.Kind], VersionID: p.VersionID})
		}
	}
	return out
}

// releaseBuild builds each release and waits. A failed release is deleted, so
// failed runs do not use up the Deployment's 20 release Drafts.
func (s *DeployStateV6) releaseBuild(ctx context.Context) error {
	for _, t := range s.Targets {
		if err := s.buildRelease(ctx, t); err != nil {
			return err
		}
	}
	return nil
}

// discardRelease deletes a release that will never go live and clears it from
// azion.json. A failure to delete is only logged: the deploy already failed.
func (s *DeployStateV6) discardRelease(ctx context.Context, t *releaseTarget) {
	cctx, cancel := cleanupContext(ctx)
	defer cancel()
	if err := s.Clients.Deployments.DeleteRelease(cctx, t.Deployment.ID, t.Release.ID); err != nil {
		logger.Debug("Could not delete a failed release", zap.String("release", t.Release.ID), zap.Error(err))
		return
	}
	s.deploymentRecord(t.Name).Release = contracts.AzionJsonV6Release{}
}

func (s *DeployStateV6) buildRelease(ctx context.Context, t *releaseTarget) error {
	if t.Release.State == "" || t.Release.State == version.Draft {
		r, err := s.Clients.Deployments.BuildRelease(ctx, t.Deployment.ID, t.Release.ID)
		if err != nil {
			// The API refused to build it (a missing dependency, say), so
			// it can never go live.
			s.discardRelease(ctx, t)
			s.writeConf()
			return &version.BuildError{Kind: "release", Name: t.Name, VersionID: t.Release.ID, Err: err}
		}
		if r.State != "" {
			t.Release.State = r.State
		}
	}

	start := time.Now()
	last := t.Release.State
	r, err := poll.Until(ctx, s.pollOptions(), s.getRelease(t),
		func(r apiDeployments.Release) (bool, error) {
			if r.State != last {
				s.say(fmt.Sprintf(msg.V6BuildTransition, "release", r.ID, t.Name, last, r.State, time.Since(start).Round(time.Second)))
				last = r.State
			}
			return version.BuildDone(r.State)
		})
	s.recordWait(time.Since(start))
	if r.ID != "" {
		t.Release = r
	}
	rec := s.deploymentRecord(t.Name)
	rec.Release.State = string(t.Release.State)
	defer s.writeConf()

	fail := func(err error) error {
		return &version.BuildError{Kind: "release", Name: t.Name, VersionID: t.Release.ID, State: t.Release.State, TraceID: t.Release.TraceID, Err: err}
	}
	if ctx.Err() != nil {
		cctx, cancel := cleanupContext(ctx)
		defer cancel()
		if cerr := s.Clients.Deployments.CancelRelease(cctx, t.Deployment.ID, t.Release.ID); cerr != nil {
			logger.Debug("Could not cancel an interrupted release build", zap.Error(cerr))
		}
		return ctx.Err()
	}
	if err != nil {
		if !errors.Is(err, version.ErrBuildFailed) {
			return fail(err)
		}
		s.discardRelease(ctx, t)
		return fail(msg.ErrorV6ReleaseFailed)
	}
	s.say(fmt.Sprintf(msg.V6ReleaseReady, t.Release.ID))
	if s.Opts.NoActivate {
		s.say(fmt.Sprintf(msg.V6ReleaseNotActivated, t.Release.ID))
	}
	return nil
}

// expectedRole is the traffic role an activated release has: ACTIVE on a
// single_version Deployment, VALID_URL on a versioned_urls one.
func expectedRole(d apiDeployments.Deployment) string {
	if d.DeploymentPolicy == apiEnvironments.PolicyVersionedURLs {
		return apiDeployments.RoleValidURL
	}
	return apiDeployments.RoleActive
}

// releaseActivate routes each Deployment's traffic to its new release and
// waits for the traffic role. On a single_version Deployment the switch is
// atomic: if it fails, the previous release stays ACTIVE.
func (s *DeployStateV6) releaseActivate(ctx context.Context) error {
	for _, t := range s.Targets {
		if _, err := s.Clients.Deployments.ActivateRelease(ctx, t.Deployment.ID, t.Release.ID); err != nil {
			return fmt.Errorf("%w: %w", msg.ErrorV6ActivationFailed, err)
		}
		want := expectedRole(t.Deployment)
		start := time.Now()
		r, err := poll.Until(ctx, s.pollOptions(), s.getRelease(t),
			func(r apiDeployments.Release) (bool, error) {
				if r.State.Failed() {
					return true, fmt.Errorf("release %s is %s: %w", r.ID, r.State, msg.ErrorV6ActivationFailed)
				}
				return r.TrafficRole == want, nil
			})
		s.recordWait(time.Since(start))
		if err != nil {
			if errors.Is(err, poll.ErrTimeout) {
				return fmt.Errorf("release %s is %s, waiting for %s: %w: %w", r.ID, r.TrafficRole, want, msg.ErrorV6ActivationFailed, err)
			}
			return err
		}
		t.Release, t.Activated = r, true
		s.deploymentRecord(t.Name).Release.Activated = true
		s.writeConf()
		s.say(fmt.Sprintf(msg.V6ReleaseActivated, r.ID, r.TrafficRole))
	}
	return nil
}

// Activated reports whether any release of this run went live.
func (s *DeployStateV6) Activated() bool {
	for _, t := range s.Targets {
		if t.Activated {
			return true
		}
	}
	return false
}

func (s *DeployStateV6) getRelease(t *releaseTarget) func(context.Context) (apiDeployments.Release, error) {
	return func(ctx context.Context) (apiDeployments.Release, error) {
		return s.Clients.Deployments.GetRelease(ctx, t.Deployment.ID, t.Release.ID)
	}
}

// deploymentRecord is azion.json's entry for a Deployment, by manifest name,
// created when missing.
func (s *DeployStateV6) deploymentRecord(name string) *contracts.AzionJsonV6Deployment {
	for i := range s.Conf.Deployments {
		if s.Conf.Deployments[i].Name == name {
			return &s.Conf.Deployments[i]
		}
	}
	s.Conf.Deployments = append(s.Conf.Deployments, contracts.AzionJsonV6Deployment{Name: name})
	return &s.Conf.Deployments[len(s.Conf.Deployments)-1]
}

// upsertNamed replaces the entry with the same name, or appends.
func upsertNamed[T any](list *[]T, item T, name func(T) string) {
	for i := range *list {
		if name((*list)[i]) == name(item) {
			(*list)[i] = item
			return
		}
	}
	*list = append(*list, item)
}
