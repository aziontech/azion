package custompages

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// Create creates a Custom Page. Its first version comes back as a Draft.
func (c *Client) Create(ctx context.Context, req sdk.CustomPageRequest) (version.Object, error) {
	logger.Debug("Create Custom Page")

	res, httpResp, err := c.apiClient.CustomPagesAPI.
		CreateCustomPage(ctx).CustomPageRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Custom Page", zap.Error(err))
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

// FindByName returns the Custom Page with exactly this name, if there is one. The
// list filter may also return near matches, which are ignored.
func (c *Client) FindByName(ctx context.Context, name string) (version.Object, bool, error) {
	logger.Debug("List Custom Pages")

	res, httpResp, err := c.apiClient.CustomPagesAPI.
		ListCustomPages(ctx).Name(name).PageSize(100).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Custom Pages", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, false, bodyErr
			}
		}
		return nil, false, version.ErrorFor(errBody, httpResp, err)
	}

	results := res.GetResults()
	for i := range results {
		if obj := &results[i]; obj.GetName() == name {
			return obj, true, nil
		}
	}
	return nil, false, nil
}
