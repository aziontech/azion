package dnssec

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, zoneID int64) (sdk.DNSSEC, error) {
	logger.Debug("Get DNSSEC")

	resp, httpResp, err := c.apiClient.DNSDNSSECAPI.
		RetrieveDnssec(ctx, zoneID).
		Execute()
	if err != nil {
		logger.Debug("Error while getting the DNSSEC", zap.Any("ID", zoneID), zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.DNSSEC{}, err
		}

		return sdk.DNSSEC{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}

func (c *Client) Update(ctx context.Context, req sdk.PatchedDNSSECRequest, zoneID int64) (sdk.DNSSEC, error) {
	logger.Debug("Update DNSSEC")

	request := c.apiClient.DNSDNSSECAPI.
		PartialUpdateDnssec(ctx, zoneID).
		PatchedDNSSECRequest(req)
	resp, httpResp, err := request.Execute()
	if err != nil {
		logger.Debug("Error while updating the DNSSEC", zap.Any("ID", zoneID), zap.Error(err))
		errBody, err := utils.LogAndRewindBodyV4(httpResp)
		if err != nil {
			return sdk.DNSSEC{}, err
		}

		return sdk.DNSSEC{}, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return resp.GetData(), nil
}
