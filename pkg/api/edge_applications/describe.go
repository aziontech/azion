package edge_applications

import (
	"context"
	"fmt"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (c *Client) Get(ctx context.Context, id string) (EdgeApplicationResponse, error) {
	logger.Debug("Get Edge Application")

	res, httpResp, err := c.apiClient.EdgeApplicationsMainSettingsApi.
		EdgeApplicationsIdGet(ctx, id).Execute()

	fmt.Println("res: ", res)
	fmt.Println("err: ", err)
	if err != nil {
		logger.Debug("Error while getting an edge application", zap.Error(err))
		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return &res.Results, nil
}
