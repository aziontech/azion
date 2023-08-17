package personal_token

import (
	"context"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"

	sdk "github.com/aziontech/azionapi-go-sdk/personal_tokens"
)

type Client struct {
	apiClient *sdk.APIClient
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.UserAgent = "Azion_CLI/" + version.BinVersion
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}
	conf.HTTPClient.Timeout = 30 * time.Second

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

type Response interface {
	GetUuid() string
	GetName() string
	GetKey() string
	GetUserId() float32
	GetCreated() time.Time
	GetExpiresAt() time.Time
	GetDescription() string
}

type Request struct {
	sdk.CreatePersonalTokenRequest
}

func (c *Client) Create(ctx context.Context, req *Request) (Response, error) {
	logger.Debug("Create Personal Token")

	response, httpResp, err := c.apiClient.PersonalTokenApi.CreatePersonalToken(ctx).
		CreatePersonalTokenRequest(req.CreatePersonalTokenRequest).Execute()

	if err != nil {
		logger.Error("Error while creating a personal token", zap.Error(err))
		logger.Debug("Status Code", zap.Any("http", httpResp.StatusCode))
		logger.Debug("Headers", zap.Any("http", httpResp.Header))
		logger.Debug("Response body", zap.Any("http", httpResp.Body))

		return nil, utils.ErrorPerStatusCode(httpResp, err)
	}

	return response, nil
}
