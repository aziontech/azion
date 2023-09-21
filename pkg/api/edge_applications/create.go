package edge_applications

import (
	"context"
	"io"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (c *Client) Create(ctx context.Context, req *CreateRequest) (EdgeApplicationsResponse, error) {
	logger.Debug("Create Edge Application")
	request := c.apiClient.EdgeApplicationsMainSettingsApi.
		EdgeApplicationsPost(ctx).CreateApplicationRequest(req.CreateApplicationRequest)

	edgeApplicationsResponse, httpResp, err := request.Execute()
	if err != nil {
		if httpResp != nil {
			logger.Debug("Error while creating an edge application", zap.Error(err))
			logger.Debug("", zap.Any("Status Code", httpResp.StatusCode))
			logger.Debug("", zap.Any("Headers", httpResp.Header))
			body, err := io.ReadAll(httpResp.Body)
			if err != nil {
				logger.Debug("Error while reading body of the http response", zap.Error(err))
				return nil, utils.ErrorPerStatusCode(httpResp, err)
			}
			logger.Debug("", zap.Any("Body", body))
		}
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &edgeApplicationsResponse.Results, nil
}
