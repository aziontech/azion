package applications

import (
	"context"

	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedApplicationList, error) {
	logger.Debug("List Applications")

	req := c.apiClient.ApplicationsAPI.
		ListApplications(ctx).Page(opts.Page).PageSize(int32(opts.PageSize)).Ordering(opts.OrderBy)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Applications", zap.Error(err))
			var rewindErr error
			errBody, rewindErr = utils.LogAndRewindBodyV4(httpResp)
			if rewindErr != nil {
				return nil, rewindErr
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
