package domains

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
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
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &domainsResponse.Results, nil
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (DomainResponse, error) {

	request := c.apiClient.DomainsApi.UpdateDomain(ctx, req.DomainId).UpdateDomainRequest(req.UpdateDomainRequest)

	domainsResponse, httpResp, err := request.Execute()
	if err != nil {
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &domainsResponse.Results, nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (sdk.DomainResponseWithResults, error) {
	resp, httpResp, err := c.apiClient.DomainsApi.GetDomains(ctx).Execute()

	if err != nil {
		return sdk.DomainResponseWithResults{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	str := strconv.FormatInt(id, 10)
	req := c.apiClient.DomainsApi.DelDomain(ctx, str)

	httpResp, err := req.Execute()

	if err != nil {
		return utils.ErrorPerStatusCode(httpResp, err)
	}

	return nil
}
