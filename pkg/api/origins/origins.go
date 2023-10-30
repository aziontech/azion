package origins

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"
)

type Request struct {
	sdk.CreateOriginsRequest
}

type Response interface {
	GetOriginKey() string
	GetOriginId() int64
	GetName() string
}

func (c *Client) Create(ctx context.Context, edgeApplicationID int64, req *Request) (Response, error) {
	logger.Debug("Create Origins")

	resp, httpResp, err := c.apiClient.EdgeApplicationsOriginsAPI.
		EdgeApplicationsEdgeApplicationIdOriginsPost(ctx, edgeApplicationID).
		CreateOriginsRequest(req.CreateOriginsRequest).Execute()

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
