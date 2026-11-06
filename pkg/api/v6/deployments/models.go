package deployments

import (
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
)

// Binding policies. The CLI creates its Deployments STRICT: a Workload bound
// to it serves exactly the release activated on it.
const (
	BindingStrict   = "STRICT"
	BindingFlexible = "FLEXIBLE"
)

// Traffic roles a release can have.
const (
	RoleActive    = "ACTIVE"
	RoleCandidate = "CANDIDATE"
	RoleValidURL  = "VALID_URL"
	RoleInactive  = "INACTIVE"
)

// OriginCLI marks a release as created by the CLI.
const OriginCLI = "cli"

// Deployment is what the deploy needs from one.
type Deployment struct {
	ID               string
	Name             string
	BindingPolicy    string
	DeploymentPolicy string
	VersionID        string
	VersionState     version.State
}

// ReleaseResource is one resource version a release carries.
type ReleaseResource struct {
	ResourceID   int64
	ResourceType string
	VersionID    string
}

// Release is what the deploy needs from one.
type Release struct {
	ID           string
	DeploymentID string
	State        version.State
	TrafficRole  string
	TraceID      string
	Resources    []ReleaseResource
}

// HierarchyEntry is one resource version a Deployment serves.
type HierarchyEntry struct {
	ResourceType string
	ResourceID   string
	Version      string
}

// Hierarchy is what a Deployment serves: its active release and the resource
// versions in it.
type Hierarchy struct {
	DeploymentID    string
	ActiveReleaseID string
	Resources       []HierarchyEntry
}

func deploymentFrom(d *sdk.Deployment) Deployment {
	return Deployment{
		ID: d.GetId(), Name: d.GetName(), BindingPolicy: d.GetBindingPolicy(),
		DeploymentPolicy: d.GetDeploymentPolicy(), VersionID: d.GetVersionId(),
		VersionState: version.State(d.GetVersionState()),
	}
}

func releaseFrom(r *sdk.Release) Release {
	out := Release{
		ID: r.GetId(), DeploymentID: r.GetDeploymentId(), State: version.State(r.GetState()),
		TrafficRole: r.GetTrafficRole(), TraceID: r.GetTraceId(),
	}
	for _, res := range r.GetResources() {
		out.Resources = append(out.Resources, ReleaseResource{
			ResourceID: res.GetResourceId(), ResourceType: res.GetResourceType(), VersionID: res.GetVersionId(),
		})
	}
	return out
}

// asyncRelease reads the answer of the calls that start work on a release
// (create, build, activate), which only carry its id and state.
func asyncRelease(deploymentID string, res *sdk.ReleaseAsyncResponse) Release {
	return Release{
		ID: res.Data.GetId(), DeploymentID: deploymentID,
		State: version.State(res.Data.GetState()), TraceID: res.Data.GetTraceId(),
	}
}

func hierarchyEntry(e *sdk.ResourceHierarchyResponseDataResourcesApplication) HierarchyEntry {
	return HierarchyEntry{ResourceType: e.GetResourceType(), ResourceID: e.GetResourceId(), Version: e.GetResourceVersion()}
}

func hierarchyFrom(h *sdk.ResourceHierarchyResponseData) Hierarchy {
	out := Hierarchy{DeploymentID: h.GetDeploymentId(), ActiveReleaseID: h.GetActiveReleaseId()}
	r := h.GetResources()
	for _, single := range []sdk.NullableResourceHierarchyResponseDataResourcesApplication{r.Application, r.Firewall, r.CustomPage} {
		if e := single.Get(); e != nil {
			out.Resources = append(out.Resources, hierarchyEntry(e))
		}
	}
	for _, many := range [][]sdk.ResourceHierarchyResponseDataResourcesWafInner{r.Waf, r.Function, r.Connector, r.NetworkList} {
		for i := range many {
			e := many[i]
			out.Resources = append(out.Resources, HierarchyEntry{
				ResourceType: e.GetResourceType(), ResourceID: e.GetResourceId(), Version: e.GetResourceVersion(),
			})
		}
	}
	return out
}
