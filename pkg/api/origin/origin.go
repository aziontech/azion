package origin

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
)

type Client struct {
	apiClient *sdk.APIClient
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

type CreateOriginsRequest struct {
	sdk.CreateOriginsRequest
}

type UpdateOriginsRequest struct {
	sdk.PatchOriginsRequest
}

type OriginsResponse interface {
	GetOriginKey() string
	GetOriginId() int64
	GetName() string
}

func (c *Client) Get(ctx context.Context, edgeApplicationID, originID int64) (sdk.OriginsResultResponse, error) {
	logger.Debug("Get Origin")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.EdgeApplicationsEdgeApplicationIdOriginsGet(ctx, edgeApplicationID).Execute()
	if err != nil {
		logger.Debug("Error while getting an origin", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return sdk.OriginsResultResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}
	if len(resp.Results) > 0 {
		for _, result := range resp.Results {
			if result.OriginId == originID {
				return result, nil
			}
		}
	}
	return sdk.OriginsResultResponse{}, utils.ErrorPerStatusCode(&http.Response{Status: "404 Not Found", StatusCode: http.StatusNotFound}, errors.New("404 Not Found"))
}

func (c *Client) ListOrigins(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.OriginsResponse, error) {
	logger.Debug("List Origins")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.EdgeApplicationsEdgeApplicationIdOriginsGet(ctx, edgeApplicationID).Execute()
	if err != nil {
		logger.Debug("Error while listing origins", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return &sdk.OriginsResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) CreateOrigins(ctx context.Context, edgeApplicationID int64, req *CreateOriginsRequest) (OriginsResponse, error) {
	logger.Debug("Create Origins")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.EdgeApplicationsEdgeApplicationIdOriginsPost(ctx, edgeApplicationID).CreateOriginsRequest(req.CreateOriginsRequest).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating an origin", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) UpdateOrigins(ctx context.Context, edgeApplicationID int64, originKey string, req *UpdateOriginsRequest) (OriginsResponse, error) {
	logger.Debug("Update Origins")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.
		EdgeApplicationsEdgeApplicationIdOriginsOriginKeyPatch(ctx, edgeApplicationID, originKey).PatchOriginsRequest(req.PatchOriginsRequest).Execute()
	if err != nil {
		logger.Debug("Error while updating an origin", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) DeleteOrigins(ctx context.Context, edgeApplicationID int64, originKey string) error {
	logger.Debug("Delete Origins")
	httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.EdgeApplicationsEdgeApplicationIdOriginsOriginKeyDelete(ctx, edgeApplicationID, originKey).Execute()
	if err != nil {
		logger.Debug("Error while deleting an origin", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
