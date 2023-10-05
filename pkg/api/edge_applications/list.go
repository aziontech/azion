package edge_applications

import (
	"context"

	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.GetApplicationsResponse, error) {
	logger.Debug("List Edge Application")

	resp, httpResp, err := c.apiClient.EdgeApplicationsMainSettingsAPI.
		EdgeApplicationsGet(ctx).Page(opts.Page).PageSize(opts.PageSize).Execute()

	if err != nil {
		logger.Debug("Error while listing edge applications", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}
