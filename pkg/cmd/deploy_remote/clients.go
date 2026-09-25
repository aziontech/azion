package deploy

import (
	apiApplications "github.com/aziontech/azion-cli/pkg/api/applications"
	apiFunction "github.com/aziontech/azion-cli/pkg/api/function"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	apiStorage "github.com/aziontech/azion-cli/pkg/api/storage"
	apiWorkload "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/registry"
)

type Clients struct {
	Function    *apiFunction.Client
	Application *apiApplications.Client
	Workload    *apiWorkload.Client
	Origin      *apiOrigin.Client
	Bucket      *apiStorage.Client
	Storage     *apiStorage.Client
}

func NewClients(f *cmdutil.Factory) *Clients {
	return &Clients{
		Function:    registry.Function(f),
		Application: registry.Applications(f),
		Workload:    registry.Workloads(f),
		Origin:      registry.Origin(f),
		Bucket:      registry.Storage(f),
		Storage:     registry.Storage(f),
	}
}
