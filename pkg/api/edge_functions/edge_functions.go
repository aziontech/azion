package edge_funtions

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	sdk "github.com/aziontech/azionapi-go-sdk/edgefunctions"
)

type Client struct {
	apiClient *sdk.APIClient
}

type EdgeFunction interface {
	GetId() int64 // Should be uint64
	GetName() string
	GetLanguage() string
	GetReferenceCount() int64 // Should be uint64
	GetModified() string
	GetInitiatorType() string
	GetLastEditor() string
	GetFunctionToRun() string
	GetJsonArgs() map[string]interface{}
	GetCode() string
}

func NewClient(c *http.Client, url string, token string) *Client {
	conf := sdk.NewConfiguration()
	conf.HTTPClient = c
	conf.AddDefaultHeader("Authorization", "token "+token)
	conf.AddDefaultHeader("Accept", "application/json;version=3")
	conf.Servers = sdk.ServerConfigurations{
		{URL: url},
	}

	return &Client{
		apiClient: sdk.NewAPIClient(conf),
	}
}

func (c *Client) Get(ctx context.Context, id int64) (EdgeFunction, error) {
	req := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdGet(ctx, id)

	res, _, err := req.Execute()

	if err != nil {
		return nil, err
	}

	return res.Results, nil
}

func (c *Client) Delete(ctx context.Context, id int64) error {
	req := c.apiClient.EdgeFunctionsApi.EdgeFunctionsIdDelete(ctx, id)

	_, err := req.Execute()

func (c *Client) Create(ctx context.Context, req CreateRequest) (uint64, error) {
	var body sdk.InlineObject

	body.SetActive(req.Active)
	body.SetCode(req.Code)
	body.SetName(req.Name)
	body.SetLanguage(req.Language)

	if req.Args != "" {
		args := make(map[string]interface{})
		if err := json.Unmarshal([]byte(req.Args), &args); err != nil {
			return 0, fmt.Errorf("failed to encode json args: %w", err)
		}
		body.SetJsonArgs(args)
	}

	request := c.apiClient.EdgeFunctionsApi.EdgeFunctionsPost(ctx).InlineObject(body)

	res, err := request.Execute()
	if err != nil {
		responseBody, _ := ioutil.ReadAll(res.Body)
		return 0, fmt.Errorf("%w: %s", err, responseBody)
	}

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return nil

}
