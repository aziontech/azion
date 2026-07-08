package datastream

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/data-stream-api"
	"go.uber.org/zap"
)

func (c *Client) ListDataSources(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedDataSourceList, error) {
	logger.Debug("List Data Sources")
	if opts.OrderBy == "" {
		opts.OrderBy = "name"
	}
	resp, httpResp, err := c.apiClient.DataStreamDataSourcesAPI.ListDataSources(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Sort).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the data sources", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
