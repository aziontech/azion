package deploy

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponseApp string = `
{
	"results":{
		"id":1697666970,
		"name":"New Edge Applicahvjgjhgjhhgtion",
		"delivery_protocol":"http",
		"http_port":80,
		"https_port":443,
		"minimum_tls_version":"",
		"active":true,
		"application_acceleration":false,
		"caching":true,
   		"debug_rules": true,
   		"http3": false,
		"supported_ciphers": "asdf",
		"device_detection":false,
		"edge_firewall":false,
		"edge_functions":false,
		"image_optimization":false,
		"load_balancer":false,
		"raw_logs":false,
		"web_application_firewall":false,
		"l2_caching": false
	},
	"schema_version":3
}
`

var sucRespDomain string = `{
    "results": {
        "id": 1702659986,
        "name": "My Domain",
        "cnames": [],
        "cname_access_only": false,
        "digital_certificate_id": null,
        "edge_application_id": 1697666970,
        "is_active": true,
        "domain_name": "ja65r2loc3.map.azionedge.net",
        "is_mtls_enabled": false,
        "mtls_verification": "enforce",
        "mtls_trusted_ca_certificate_id": null
    },
    "schema_version": 3
}`

var sucRespOrigin string = `{
    "results": {
        "origin_id": 116207,
        "origin_key": "35e3a635-2227-4bb6-976c-5e8c8fa58a67",
        "name": "Create Origin22",
        "origin_type": "single_origin",
        "addresses": [
            {
                "address": "httpbin.org",
                "weight": null,
                "server_role": "primary",
                "is_active": true
            }
        ],
        "origin_protocol_policy": "http",
        "is_origin_redirection_enabled": false,
        "host_header": "${host}",
        "method": "",
        "origin_path": "/requests",
        "connection_timeout": 60,
        "timeout_between_bytes": 120,
        "hmac_authentication": false,
        "hmac_region_name": "",
        "hmac_access_key": "",
        "hmac_secret_key": ""
    },
    "schema_version": 3
}`

var sucRespCacheSettings = `{
    "results": {
        "id": 138708,
        "name": "Default Cache Settings2234",
        "browser_cache_settings": "override",
        "browser_cache_settings_maximum_ttl": 0,
        "cdn_cache_settings": "override",
        "cdn_cache_settings_maximum_ttl": 60,
        "cache_by_query_string": "ignore",
        "query_string_fields": null,
        "enable_query_string_sort": false,
        "cache_by_cookies": "ignore",
        "cookie_names": null,
        "adaptive_delivery_action": "ignore",
        "device_group": [],
        "enable_caching_for_post": false,
        "l2_caching_enabled": false,
        "is_slice_configuration_enabled": false,
        "is_slice_edge_caching_enabled": false,
        "is_slice_l2_caching_enabled": false,
        "slice_configuration_range": 1024,
        "enable_caching_for_options": false,
        "enable_stale_cache": true,
        "l2_region": null
    },
    "schema_version": 3
}`

var sucRespRules = `{
    "results": {
        "id": 214790,
        "name": "testorigin2",
        "phase": "request",
        "behaviors": [
            {
                "name": "set_origin",
                "target": "116207"
            }
        ],
        "criteria": [
            [
                {
                    "variable": "${uri}",
                    "operator": "starts_with",
                    "conditional": "if",
                    "input_value": "/"
                }
            ]
        ],
        "is_active": true,
        "order": 1,
        "description": ""
    },
    "schema_version": 3
}`

func TestDeployCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("full flow manifest empty", func(t *testing.T) {
		mock := &httpmock.Registry{}

		options := &contracts.AzionApplicationOptions{
			Name: "LovelyName",
		}

		dat, _ := os.ReadFile("./fixtures/create_app.json")
		_ = json.Unmarshal(dat, options)

		mock.Register(
			httpmock.REST("POST", "edge_applications"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("POST", "domains"),
			httpmock.JSONFromString(sucRespDomain),
		)

		mock.Register(
			httpmock.REST("PATCH", "edge_applications/1697666970"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1697666970/functions_instances"),
			httpmock.JSONFromString(successResponseApp),
		)

		f, _, _ := testutils.NewFactory(mock)
		deployCmd := NewDeployCmd(f)
		clients := NewClients(f)

		manifest := Manifest{}

		deployCmd.FilepathWalk = func(root string, fn filepath.WalkFunc) error {
			return nil
		}
		deployCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions) error {
			return nil
		}

		err := manifest.Interpreted(f, deployCmd, options, clients)
		require.NoError(t, err)
	})

	//t.Run("full flow manifest completed", func(t *testing.T) {
	//	mock := &httpmock.Registry{}
	//
	//	options := &contracts.AzionApplicationOptions{
	//		Name: "LovelyName",
	//	}
	//
	//	dat, _ := os.ReadFile("./fixtures/create_app.json")
	//	_ = json.Unmarshal(dat, options)
	//
	//	mock.Register(
	//		httpmock.REST("POST", "edge_applications"),
	//		httpmock.JSONFromString(successResponseApp),
	//	)
	//
	//	mock.Register(
	//		httpmock.REST("POST", "domains"),
	//		httpmock.JSONFromString(sucRespDomain),
	//	)
	//
	//	mock.Register(
	//		httpmock.REST("POST", "edge_applications/1697666970/origins"),
	//		httpmock.JSONFromString(sucRespOrigin),
	//	)
	//
	//	mock.Register(
	//		httpmock.REST("POST", "edge_applications/1697666970/cache_settings"),
	//		httpmock.JSONFromString(sucRespCacheSettings),
	//	)
	//
	//	mock.Register(
	//		httpmock.REST("POST", "edge_applications/1697666970/rules_engine/request/rules"),
	//		httpmock.JSONFromString(sucRespRules),
	//	)
	//
	//	mock.Register(
	//		httpmock.REST("POST", "edge_applications/1697666971/rules_engine/request/rules"),
	//		httpmock.JSONFromString(sucRespRules),
	//	)
	//
	//	mock.Register(
	//		httpmock.REST("PATCH", "edge_applications/1697666970"),
	//		httpmock.JSONFromString(successResponseApp),
	//	)
	//
	//	mock.Register(
	//		httpmock.REST("POST", "edge_applications/1697666970/functions_instances"),
	//		httpmock.JSONFromString(successResponseApp),
	//	)
	//
	//	f, _, _ := testutils.NewFactory(mock)
	//	deployCmd := NewDeployCmd(f)
	//	clients := NewClients(f)
	//
	//	manifest := Manifest{
	//		Routes: []Routes{
	//			{
	//				From:     "/_next/static/",
	//				To:       ".edge/storage",
	//				Priority: 1,
	//				Type:     "deliver",
	//			},
	//
	//			{
	//				From:     "/_next/data/",
	//				To:       ".edge/storage",
	//				Priority: 2,
	//				Type:     "deliver",
	//			},
	//			{
	//				From:     "\\.(css|js|ttf|woff|woff2|pdf|svg|jpg|jpeg|gif|bmp|png|ico|mp4)$",
	//				To:       ".edge/storage",
	//				Priority: 3,
	//				Type:     "deliver",
	//			},
	//			{
	//				From:     "/",
	//				To:       ".edge/worker.js",
	//				Priority: 4,
	//				Type:     "compute",
	//			},
	//		},
	//	}
	//
	//	deployCmd.FilepathWalk = func(root string, fn filepath.WalkFunc) error {
	//		return nil
	//	}
	//	deployCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions) error {
	//		return nil
	//	}
	//
	//	err := manifest.Interpreted(f, deployCmd, options, clients)
	//	require.NoError(t, err)
	//})

	t.Run("without azion.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		deployCmd := NewDeployCmd(f)

		deployCmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := NewCobraCmd(deployCmd)

		cmd.SetArgs([]string{""})

		err := cmd.Execute()

		require.EqualError(t, err, "Failed to build your resource. Azion configuration not found. Make sure you are in the root directory of your local repository and have already initialized or linked your resource with the commands 'azion init' or 'azion link'")
	})

	t.Run("failed to create application", func(t *testing.T) {

		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{
			Name: "NotAVeryGoodName",
		}

		dat, _ := os.ReadFile("./fixtures/create_app.json")
		_ = json.Unmarshal(dat, options)

		mock.Register(
			httpmock.REST("POST", "edge_applications"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)
		ctx := context.Background()

		cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

		cmd := NewDeployCmd(f)

		_, err := cmd.createApplication(cliapp, ctx, options)
		require.ErrorContains(t, err, "Failed to create the Edge Application")
	})

	t.Run("create application success", func(t *testing.T) {

		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{
			Name: "LovelyName",
		}

		dat, _ := os.ReadFile("./fixtures/create_app.json")
		_ = json.Unmarshal(dat, options)

		mock.Register(
			httpmock.REST("POST", "edge_applications"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("PATCH", "edge_applications/666"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/666/functions_instances"),
			httpmock.JSONFromString(successResponseApp),
		)

		f, _, _ := testutils.NewFactory(mock)
		ctx := context.Background()

		cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

		cmd := NewDeployCmd(f)

		_, err := cmd.createApplication(cliapp, ctx, options)
		require.NoError(t, err)
	})
}
