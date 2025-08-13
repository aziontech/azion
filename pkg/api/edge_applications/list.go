package edge_applications

import (
	"context"

	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge-api"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedEdgeApplicationList, error) {
	logger.Debug("List Edge Applications")

	req := c.apiClient.EdgeApplicationsAPI.
		ListEdgeApplications(ctx).Page(opts.Page).PageSize(opts.PageSize)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Edge Applications", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
