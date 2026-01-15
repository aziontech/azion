package kv

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/kv-api"
	"go.uber.org/zap"
)

// func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.NamespaceList, error) {
// 	logger.Debug("List namespaces")
// 	if opts.OrderBy == "" {
// 		opts.OrderBy = "id"
// 	}
// 	resp, httpResp, err := c.apiClient.KVNamespacesAPI.ListNamespaces(ctx).
// 		Page(opts.Page).
// 		PageSize(opts.PageSize).
// 		Execute()

// 	if err != nil {
// 		errBody := ""
// 		if httpResp != nil {
// 			logger.Debug("Error while listing the namespaces", zap.Error(err))
// 			errBody, err = utils.LogAndRewindBodyV4(httpResp)
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
// 	}

// 	return resp, nil
// }

func (c *Client) Create(ctx context.Context, req CreateRequest) (*sdk.Namespace, error) {
	logger.Debug("Create namespace")

	request := c.apiClient.KVNamespacesAPI.CreateNamespace(ctx).NamespaceCreateRequest(req.ConnectorPolymorphicRequest)

	response, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a namespace", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return response, nil
}

func (c *Client) Get(ctx context.Context, namespace string) (*sdk.Namespace, error) {
	logger.Debug("Retrieve namespace")
	request := c.apiClient.KVNamespacesAPI.RetrieveNamespace(ctx, namespace)

	res, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while retrieving a namespace", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return res, nil
}
