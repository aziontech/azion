package dnsrecord

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) Create(ctx context.Context, req sdk.RecordRequest, zoneID int64) (sdk.Record, error) {
	logger.Debug("Create DNS Record")

	request := c.apiClient.DNSRecordsAPI.
		CreateDnsRecord(ctx, zoneID).
		RecordRequest(req)
	resp, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while creating a DNS Record", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.Record{}, err
		}

		return sdk.Record{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) Update(ctx context.Context, req sdk.PatchedRecordRequest, zoneID, recordID int64) (sdk.Record, error) {
	logger.Debug("Update DNS Record")

	request := c.apiClient.DNSRecordsAPI.
		PartialUpdateDnsRecord(ctx, recordID, zoneID).
		PatchedRecordRequest(req)
	resp, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating a DNS Record", zap.Any("ID", recordID), zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.Record{}, err
		}

		return sdk.Record{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) List(ctx context.Context, opts *contracts.ListOptions, zoneID int64) (*sdk.PaginatedRecordList, error) {
	logger.Debug("List DNS Records")
	if opts.OrderBy == "" {
		opts.OrderBy = "id"
	}

	resp, httpResp, err := c.apiClient.DNSRecordsAPI.
		ListDnsRecords(ctx, zoneID).
		Ordering(opts.OrderBy).
		Page(opts.Page).
		PageSize(opts.PageSize).
		Search(opts.Filter).
		Execute()
	if err != nil {
		logger.Debug("Error while listing DNS Records", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return nil, err
		}

		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp, nil
}

func (c *Client) Get(ctx context.Context, zoneID, recordID int64) (sdk.Record, error) {
	logger.Debug("Get DNS Record")

	resp, httpResp, err := c.apiClient.DNSRecordsAPI.
		RetrieveDnsRecord(ctx, recordID, zoneID).
		Execute()
	if err != nil {
		logger.Debug("Error while getting a DNS Record", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.Record{}, err
		}

		return sdk.Record{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) Delete(ctx context.Context, zoneID, recordID int64) error {
	logger.Debug("Delete DNS Record")

	req := c.apiClient.DNSRecordsAPI.
		DeleteDnsRecord(ctx, recordID, zoneID)
	_, httpResp, err := req.Execute()
	if err != nil {
		logger.Debug("Error while deleting a DNS Record", zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return err
		}

		return utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return nil
}
