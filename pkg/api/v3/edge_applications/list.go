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

	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.EdgeApplicationsMainSettingsAPI.
		EdgeApplicationsGet(ctx).
		OrderBy(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Sort(opts.Sort).Execute()

	if err != nil {
		logger.Debug("Error while listing Edge Applications", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}
