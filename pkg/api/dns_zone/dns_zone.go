package dnszone

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) Create(ctx context.Context, req sdk.ZoneRequest) (sdk.Zone, error) {
	logger.Debug("Create DNS Zone")

	request := c.apiClient.DNSZonesAPI.
		CreateDnsZone(ctx).
		ZoneRequest(req)
	resp, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating a DNS Zone", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.Zone{}, err
		}

		return sdk.Zone{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) Update(ctx context.Context, req sdk.PatchedUpdateZoneRequest, zoneID int64) (sdk.Zone, error) {
	logger.Debug("Update DNS Zone")

	request := c.apiClient.DNSZonesAPI.
		PartialUpdateDnsZone(ctx, zoneID).
		PatchedUpdateZoneRequest(req)
	resp, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a DNS Zone", zap.Any("ID", zoneID), zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.Zone{}, err
		}

		return sdk.Zone{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedZoneList, error) {
	logger.Debug("List DNS Zones")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.DNSZonesAPI.
		ListDnsZones(ctx).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Filter).
		Execute()
	if err != nil {
		logger.Debug("Error while listing DNS Zones", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return nil, err
		}

		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, zoneID int64) (sdk.Zone, error) {
	logger.Debug("Get DNS Zone")

	resp, httpResp, err := c.apiClient.DNSZonesAPI.
		RetrieveDnsZone(ctx, zoneID).
		Execute()
	if err != nil {
		logger.Debug("Error while getting a DNS Zone", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.Zone{}, err
		}

		return sdk.Zone{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) Delete(ctx context.Context, zoneID int64) error {
	logger.Debug("Delete DNS Zone")

	req := c.apiClient.DNSZonesAPI.
		DeleteDnsZone(ctx, zoneID)
	_, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting a DNS Zone", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return err
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}
