package kv

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/kv-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.NamespaceList, error) {
	logger.Debug("List Functions")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}
	resp, httpResp, err := c.apiClient.KVNamespacesAPI.ListNamespaces(ctx).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing the Functions", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}
