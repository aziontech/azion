package domains

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/domains"
)

type Client struct {
	apiClient *sdk.APIClient
}

type CreateRequest struct {
	sdk.CreateDomainRequest
}

type UpdateRequest struct {
	sdk.UpdateDomainRequest
	DomainId string
}

type DomainResponse interface {
	GetId() int64
	GetDomainName() string
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

func (c *Client) Create(ctx context.Context, req *CreateRequest) (DomainResponse, error) {

	request := c.apiClient.DomainsApi.CreateDomain(ctx).CreateDomainRequest(req.CreateDomainRequest)

	domainsResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp == nil || httpResp.StatusCode >= 500 {
			err := utils.CheckStatusCode500Error(err)
			return nil, err
		}
		responseBody, _ := ioutil.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("%w: %s", err, responseBody)
	}

	return &domainsResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (DomainResponse, error) {

	request := c.apiClient.DomainsApi.UpdateDomain(ctx, req.DomainId).UpdateDomainRequest(req.UpdateDomainRequest)

	domainsResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp == nil || httpResp.StatusCode >= 500 {
			err := utils.CheckStatusCode500Error(err)
			return nil, err
		}
		responseBody, _ := ioutil.ReadAll(httpResp.Body)
		return nil, fmt.Errorf("%w: %s", err, responseBody)
	}

	return &domainsResponse.Results, nil
}
