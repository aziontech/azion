package deploy

import (
	apiDomain "github.com/aziontech/azion-cli/pkg/api/domain"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	apiEdgeFunction "github.com/aziontech/azion-cli/pkg/api/edge_function"
	apiOrigin "github.com/aziontech/azion-cli/pkg/api/origin"
	apiStorage "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
)

type Clients struct {
	EdgeFunction    *apiEdgeFunction.Client
	EdgeApplication *apiEdgeApplications.Client
	Domain          *apiDomain.Client
	Origin          *apiOrigin.Client
	Bucket          *apiStorage.ClientStorage
	Storage         *apiStorage.Client
}

func NewClients(f *cmdutil.Factory) *Clients {
	httpClient := f.HttpClient
	apiURL := f.Config.GetString("api_url")
	token := f.Config.GetString("token")

	return &Clients{
		EdgeFunction:    apiEdgeFunction.NewClient(httpClient, apiURL, token),
		EdgeApplication: apiEdgeApplications.NewClient(httpClient, apiURL, token),
		Domain:          apiDomain.NewClient(httpClient, apiURL, token),
		Origin:          apiOrigin.NewClient(httpClient, apiURL, token),
		Bucket:          apiStorage.NewClientStorage(httpClient, "https://stage-api.azion.com", "azion6a25b8d66b05215640b496e73b1281a1811"),
		Storage:         apiStorage.NewClient(httpClient, apiURL, token),
	}
}
