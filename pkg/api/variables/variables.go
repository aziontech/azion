package variables

import (
	"context"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/utils"

	sdk "github.com/aziontech/azionapi-go-sdk/variables"
)

type Client struct {
	apiClient *sdk.APIClient
}

type VariableResponse interface {
	GetUuid() string
	GetKey() string
	GetValue() string
	GetSecret() bool
	GetLastEditor() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	conf.HTTPClient.Timeout = 30 * time.Second

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) List(ctx context.Context) ([]VariableResponse, error) {
	resp, httpResp, err := c.apiClient.VariablesApi.ApiVariablesList(ctx).Execute()

	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	var result []VariableResponse

	for i := range resp {
		result = append(result, &resp[i])
	}

	return result, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	req := c.apiClient.VariablesApi.ApiVariablesDestroy(ctx, id)

	httpResp, err := req.Execute()

	if err != nil {
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, id string) (VariableResponse, error) {
	req := c.apiClient.VariablesApi.ApiVariablesRetrieve(ctx, id)
	res, httpResp, err := req.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return res, nil
}
