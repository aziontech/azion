package firewallinstance

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions, firewallID int64) (*sdk.PaginatedFirewallFunctionInstanceList, error) {
	logger.Debug("List Firewall Function Instances")

	req := c.apiClient.FirewallsFunctionAPI.
		ListFirewallFunction(ctx, firewallID).
		Page(opts.Page).
		PageSize(opts.PageSize)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Firewall Function Instances", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return nil, err
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, firewallID, functionInstanceID int64) (sdk.FirewallFunctionInstance, error) {
	logger.Debug("Get Firewall Function Instance")
	req := c.apiClient.FirewallsFunctionAPI.RetrieveFirewallFunction(ctx, firewallID, functionInstanceID)
	resp, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while getting Firewall Function Instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.FirewallFunctionInstance{}, err
			}
		}
		return sdk.FirewallFunctionInstance{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}

func (c *Client) Delete(ctx context.Context, firewallID, functionInstanceID int64) error {
	logger.Debug("Delete Firewall Function Instance")
	req := c.apiClient.FirewallsFunctionAPI.DeleteFirewallFunction(ctx, firewallID, functionInstanceID)
	_, httpResp, err := req.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while deleting Firewall Function Instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return err
			}
		}
		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}

func (c *Client) Create(ctx context.Context, firewallID int64, req *CreateRequest) (sdk.FirewallFunctionInstance, error) {
	logger.Debug("Create Firewall Function Instance")
	request := c.apiClient.FirewallsFunctionAPI.CreateFirewallFunction(ctx, firewallID).FirewallFunctionInstanceRequest(req.FirewallFunctionInstanceRequest)
	resp, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating Firewall Function Instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.FirewallFunctionInstance{}, err
			}
		}
		return sdk.FirewallFunctionInstance{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}

func (c *Client) Update(ctx context.Context, firewallID, instanceID int64, req *UpdateRequest) (sdk.FirewallFunctionInstance, error) {
	logger.Debug("Update Firewall Function Instance")
	request := c.apiClient.FirewallsFunctionAPI.PartialUpdateFirewallFunction(ctx, firewallID, instanceID).PatchedFirewallFunctionInstanceRequest(req.PatchedFirewallFunctionInstanceRequest)
	resp, httpResp, err := request.Execute()

	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating Firewall Function Instance", zap.Error(err))
			errBody, err = utils.LogAndRewindBodyV4(httpResp)
			if err != nil {
				return sdk.FirewallFunctionInstance{}, err
			}
		}
		return sdk.FirewallFunctionInstance{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.Data, nil
}
