package applications

import (
	"context"

	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

type ApplicationResponse interface {
	GetId() int64
	GetName() string
	GetActive() bool
	GetDebug() bool
	GetLastEditor() string
	GetModules() sdk.ApplicationModules
	GetProductVersion() string
	GetVersionState() string
}

type CreateRequest struct {
	sdk.ApplicationRequest
}

func (c *Client) Create(ctx context.Context, req *CreateRequest) (ApplicationResponse, error) {
	logger.Debug("Create Application")
	request := c.apiClient.ApplicationsAPI.CreateApplication(ctx).ApplicationRequest(req.ApplicationRequest)
	resp, httpResp, err := request.Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating an Application", zap.Error(err), zap.Any("Name", req.Name))
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
