package storage

import (
    // "github.com/aziontech/azion-cli/pkg/cmd/version"
    "context"
    "net/http"
    "time"
    "os"

    "github.com/aziontech/azion-cli/utils"
    sdk "github.com/aziontech/azionapi-go-sdk/storageapi"
    // "context"
)

type Client struct {
    apiClient *sdk.APIClient
}

type UploadRequest struct {
    sdk.ApiStorageVersionIdPostRequest
}

// func NewClient(c *http.Client, url string, token string) *Client {
//     conf := sdk.NewConfiguration()
//     conf.HTTPClient = c
//     conf.AddDefaultHeader("Authorization", "token "+token)
//     conf.AddDefaultHeader("Accept", "application/json;version=3")
//     conf.UserAgent = "Azion_CLI/" + version.BinVersion
//     conf.Servers = sdk.ServerConfigurations{
//         {URL: url},
//     }
//     conf.HTTPClient.Timeout = 30 * time.Second
//
//     return &Client{
//         apiClient: sdk.NewAPIClient(conf),
//     }
// }

func NewClient(c *http.Client, url string, token string) *Client {
    conf := sdk.NewConfiguration()
    conf.AddDefaultHeader("Authorization", "Token <token-RTM>")
    conf.AddDefaultHeader("Accept", "application/json; version=3")
    conf.AddDefaultHeader("Content-Type", "application/json")
    // conf.AddDefaultHeader("Content-Disposition", "attachment; filename=\"hello.js\"")
    conf.Servers = sdk.ServerConfigurations{
        {URL: url},
    }
    conf.HTTPClient.Timeout = 30 * time.Second
    return &Client{
        sdk.NewAPIClient(conf),
    }
}


// api := api.NewAPIClient(conf).DefaultApi
// versionID := "b5392936c6c4af3c123c1d0a94d8c576"
//
// fileContent, err := os.Open("file_js_local.js")
// if err != nil {
// 	log.Fatalf("An error occured while reading fileContent: %v", err)
// 	return
// }

func (c *Client) Upload(ctx context.Context, versionID, path string, file *os.File) (error) {
    req := c.apiClient.DefaultApi.StorageVersionIdPost(ctx, versionID).XAzionStaticPath(path).Body(file)
    _, httpResp, err := req.Execute()
    if err != nil {
        return utils.ErrorPerStatusCode(httpResp, err)
    }
    return nil
}
