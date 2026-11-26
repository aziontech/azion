package connectors

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
	"go.uber.org/zap"
)

// Create creates a Connector. Its first version comes back as a Draft.
func (c *Client) Create(ctx context.Context, req sdk.ConnectorRequest) (version.Object, error) {
	logger.Debug("Create Connector")

	res, httpResp, err := c.apiClient.ConnectorsAPI.
		CreateConnector(ctx).ConnectorRequest(req).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while creating a Connector", zap.Error(err))
			var bodyErr error
			errBody, bodyErr = utils.LogAndRewindBodyV4(httpResp)
			if bodyErr != nil {
				return nil, bodyErr
			}
		}
		return nil, version.ErrorFor(errBody, httpResp, err)
	}

	return object(&res.Data), nil
}

// FindByName returns the Connector with exactly this name, if there is one. The
// list filter may also return near matches, which are ignored.
func (c *Client) FindByName(ctx context.Context, name string) (version.Object, bool, error) {
	logger.Debug("List Connectors")

	res, httpResp, err := c.apiClient.ConnectorsAPI.
		ListConnectors(ctx).Name(name).PageSize(100).Execute()
	if err != nil {
		errBody := ""
		if httpResp != nil {
			logger.Debug("Error while listing Connectors", zap.Error(err))
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
		if obj := object(&results[i]); obj.GetName() == name {
			return obj, true, nil
		}
	}
	return nil, false, nil
}

// object returns the variant a Connector holds (HTTP, storage or live
// ingest); each has the getters version.Object asks for.
// versionMeta reads the variant a ConnectorVersion holds. An empty one (no
// variant matched) has no id and no state, which no wait treats as settled.
func versionMeta(v *sdk.ConnectorVersion) version.Meta {
	if v != nil {
		if m, ok := v.GetActualInstance().(version.Meta); ok {
			return m
		}
	}
	return &sdk.ConnectorHTTPVersion{}
}

func object(c *sdk.Connector) version.Object {
	if obj, ok := c.GetActualInstance().(version.Object); ok {
		return obj
	}
	return &sdk.ConnectorHTTP{}
}
