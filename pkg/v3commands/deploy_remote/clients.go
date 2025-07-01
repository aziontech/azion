package deploy

import (
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	apiDomain "github.com/aziontech/azion-cli/pkg/v3api/domain"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	apiEdgeFunction "github.com/aziontech/azion-cli/pkg/v3api/edge_function"
	apiOrigin "github.com/aziontech/azion-cli/pkg/v3api/origin"
	apiStorage "github.com/aziontech/azion-cli/pkg/v3api/storage"
)

type Clients struct {
	EdgeFunction    *apiEdgeFunction.Client
	EdgeApplication *apiEdgeApplications.Client
	Domain          *apiDomain.Client
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
		Domain:          apiDomain.NewClient(httpClient, apiURL, token),
		Origin:          apiOrigin.NewClient(httpClient, apiURL, token),
		Bucket:          apiStorage.NewClient(httpClient, storageURL, token),
		Storage:         apiStorage.NewClient(httpClient, storageURL, token),
	}
}
