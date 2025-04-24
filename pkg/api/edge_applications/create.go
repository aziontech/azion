package edge_applications

import (
	"context"

	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

type ResponseCreate interface {
	GetId() int64
	GetName() string
	GetActive() bool
	GetApplicationAcceleration() bool
	GetCaching() bool
	GetDeliveryProtocol() string
	GetDeviceDetection() bool
	GetEdgeFirewall() bool
	GetEdgeFunctions() bool
	GetHttpPort() interface{}
	GetHttpsPort() interface{}
	GetImageOptimization() bool
	GetL2Caching() bool
	GetLoadBalancer() bool
	GetMinimumTlsVersion() string
	GetRawLogs() bool
	GetWebApplicationFirewall() bool
}

type EdgeApplicationsResponse interface {
	GetId() int64
	GetName() string
}

type CreateRequest struct {
	sdk.EdgeApplicationRequest
}

func (c *Client) Create(ctx context.Context, req *CreateRequest,
) (EdgeApplicationsResponse, error) {
	logger.Debug("Create Edge Application")
	request := c.apiClient.EdgeApplicationsAPI.
		CreateEdgeApplication(ctx).EdgeApplicationRequest(req.EdgeApplicationRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating an Edge Application", zap.Error(err))
			err := utils.LogAndRewindBody(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Data, nil
}
