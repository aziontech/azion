package deploy

import (
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiEdgeFunction "github.com/aziontech/azion-cli/pkg/api/edge_function"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	apiStorage "github.com/aziontech/azion-cli/pkg/api/storage"
	apiWorkload "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
)

type Clients struct {
	EdgeFunction    *apiEdgeFunction.Client
	EdgeApplication *apiEdgeApplications.Client
	Workload        *apiWorkload.Client
	Origin          *apiOrigin.Client
	Bucket          *apiStorage.Client
	Storage         *apiStorage.Client
}

func NewClients(f *cmdutil.Factory) *Clients {
	httpClient := f.HttpClient
	apiURL := f.Config.GetString("api_url")
	storageURL := f.Config.GetString("storage_url")
	token := f.Config.GetString("token")

	return &Clients{
		EdgeFunction:    apiEdgeFunction.NewClient(httpClient, apiURL, token),
		EdgeApplication: apiEdgeApplications.NewClient(httpClient, apiURL, token),
		Workload:        apiWorkload.NewClient(httpClient, apiURL, token),
		Origin:          apiOrigin.NewClient(httpClient, apiURL, token),
		Bucket:          apiStorage.NewClient(httpClient, storageURL, token),
		Storage:         apiStorage.NewClient(httpClient, storageURL, token),
	}
}
