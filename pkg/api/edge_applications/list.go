package edge_applications

import (
	"context"

	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedResponseListEdgeApplicationList, error) {
	logger.Debug("List Edge Applications")

	resp, httpResp, err := c.apiClient.EdgeApplicationsAPI.
		ListEdgeApplications(ctx).Page(opts.Page).PageSize(opts.PageSize).Execute()

	if err != nil {
		logger.Debug("Error while listing Edge Applications", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return resp, nil
}
