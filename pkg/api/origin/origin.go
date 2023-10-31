package origin

import (
	"context"
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

type CreateRequest struct {
	sdk.CreateOriginsRequest
}

type UpdateRequest struct {
	sdk.PatchOriginsRequest
}

type Response interface {
	GetOriginKey() string
	GetOriginId() int64
	GetName() string
}

func (c *Client) Get(ctx context.Context, edgeApplicationID int64, originKey string) (sdk.OriginsResultResponse, error) {
	logger.Debug("Get Origin")

	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.
		EdgeApplicationsEdgeApplicationIdOriginsOriginKeyGet(ctx, edgeApplicationID, originKey).Execute()

	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while describing an origin", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return sdk.OriginsResultResponse{}, err
			}
		}
		return sdk.OriginsResultResponse{}, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp.Results, nil
}

func (c *Client) ListOrigins(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.OriginsResponse, error) {
	logger.Debug("List Origins")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.EdgeApplicationsEdgeApplicationIdOriginsGet(ctx, edgeApplicationID).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while listing your origins", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return resp, nil
}

func (c *Client) Create(ctx context.Context, edgeApplicationID int64, req *CreateRequest) (Response, error) {
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

func (c *Client) Update(ctx context.Context, edgeApplicationID int64, originKey string, req *UpdateRequest) (Response, error) {
	logger.Debug("Update Origins")
	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.
		EdgeApplicationsEdgeApplicationIdOriginsOriginKeyPatch(ctx, edgeApplicationID, originKey).PatchOriginsRequest(req.PatchOriginsRequest).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while updating an origin", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}
	return &resp.Results, nil
}

func (c *Client) DeleteOrigins(ctx context.Context, edgeApplicationID int64, originKey string) error {
	logger.Debug("Delete Origins")
	httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.EdgeApplicationsEdgeApplicationIdOriginsOriginKeyDelete(ctx, edgeApplicationID, originKey).Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while deleting an origin", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCode(httpResp, err)
	}
	return nil
}
