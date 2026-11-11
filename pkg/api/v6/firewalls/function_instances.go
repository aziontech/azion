package firewalls

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// ListFunctionInstances lists the firewall function instances of one version, one page at a time.
func (c *Client) ListFunctionInstances(ctx context.Context, id int64, versionID string, opts *contracts.ListOptions) (*sdk.PaginatedFirewallFunctionInstanceList, error) {
	logger.Debug("List Firewall firewall function instances")

	req := c.apiClient.FirewallsFunctionAPI.
		ListFirewallVersionFunctions(ctx, id, versionID).Page(opts.Page).PageSize(int32(opts.PageSize))
	if opts.OrderBy != "" {
		req = req.Ordering(opts.OrderBy)
	}
	res, httpResp, err := req.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing firewall function instances", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return res, nil
}

// CreateFunctionInstance creates a firewall function instance in a Draft.
func (c *Client) CreateFunctionInstance(ctx context.Context, id int64, versionID string, req sdk.FirewallFunctionInstanceRequest) (version.Item, error) {
	logger.Debug("Create Firewall firewall function instance")

	res, httpResp, err := c.apiClient.FirewallsFunctionAPI.
		CreateFirewallVersionFunctions(ctx, id, versionID).FirewallFunctionInstanceRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a firewall function instance", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return &res.Data, nil
}

// UpdateFunctionInstance replaces a firewall function instance in a Draft.
func (c *Client) UpdateFunctionInstance(ctx context.Context, id int64, versionID string, subID int64, req sdk.FirewallFunctionInstanceRequest) (version.Item, error) {
	logger.Debug("Update Firewall firewall function instance")

	res, httpResp, err := c.apiClient.FirewallsFunctionAPI.
		UpdateFirewallVersionFunction(ctx, id, subID, versionID).FirewallFunctionInstanceRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating a firewall function instance", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return &res.Data, nil
}

// DeleteFunctionInstance removes a firewall function instance from a Draft.
func (c *Client) DeleteFunctionInstance(ctx context.Context, id int64, versionID string, subID int64) error {
	logger.Debug("Delete Firewall firewall function instance")

	httpResp, err := c.apiClient.FirewallsFunctionAPI.
		DeleteFirewallVersionFunction(ctx, id, subID, versionID).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting a firewall function instance", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return bodyErr
			}
		}
		return version.ErrorFor(errBody, httpResp, err)
	}

	return nil
}
