package applications

import (
	"context"

	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

type UpdateRequest struct {
	sdk.PatchedApplicationRequest
	Id int64
}

func (c *Client) Update(ctx context.Context, req *UpdateRequest) (ApplicationResponse, error) {
	logger.Debug("Update Application")
	request := c.apiClient.ApplicationsAPI.PartialUpdateApplication(ctx, req.Id).PatchedApplicationRequest(req.PatchedApplicationRequest)

	resp, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while updating an Application", zap.Error(err), zap.Any("ID", req.Id), zap.Any("Name", req.Name))
			var rewindErr error
			errBody, rewindErr = utils.LogAndRewindBodyV4(httpResp)
			if rewindErr != nil {
				return nil, rewindErr
			}
		}
		return nil, utils.ErrorPerStatusCodeV4(errBody, httpResp, err)
	}

	return &resp.Data, nil
}
