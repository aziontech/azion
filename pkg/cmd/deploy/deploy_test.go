package deploy

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

// var successRespRule string = `
// {
// 	"results": {
// 	  "id": 234567,
// 	  "name": "enable gzip",
// 	  "phase": "response",
// 	  "behaviors": [
// 		{
// 		  "name": "enable_gzip",
// 		}
// 	  ],
// 	  "criteria": [
// 		[
// 		  {
// 			"variable": "${uri}",
// 			"operator": "exists",
// 			"conditional": "if",
// 			"input_value": ""
// 		  }
// 		]
// 	  ],
// 	  "is_active": true,
// 	  "order": 1,
// 	},
// 	"schema_version": 3
//   }
// `

// var successRespOrigin string = `
//
//	{
//		"results": {
//		  "origin_id": 0,
//		  "origin_key": "000000-000000-00000-00000-000000",
//		  "name": "name",
//		  "origin_type": "single_origin",
//		  "addresses": [
//			{
//			  "address": "httpbin.org",
//			  "weight": null,
//			  "server_role": "primary",
//			  "is_active": true
//			}
//		  ],
//		  "origin_protocol_policy": "http",
//		  "is_origin_redirection_enabled": false,
//		  "host_header": "${host}",
//		  "method": "",
//		  "origin_path": "/requests",
//		  "connection_timeout": 60,
//		  "timeout_between_bytes": 120,
//		  "hmac_authentication": false,
//		  "hmac_region_name": "",
//		  "hmac_access_key": "",
//		  "hmac_secret_key": ""
//		},
//		"schema_version": 3
//	  }
//
// `
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

// var sucRespInst string = `{
// 	"results": {
// 	  "edge_function_id": 1111,
// 	  "name": "Edge Function",
// 	  "args": {},
// 	  "id": 101001
// 	},
// 	"schema_version": 3
//   }
// `

// var sucRespFunc string = `{
// 	"results": {
// 	  "id": 1111,
// 	  "name": "Function Test API",
// 	  "language": "javascript",
// 	  "code": "{\r\n    async function handleRequest(request) {\r\n        return new Response(\"Hello world in a new response\");\r\n    }\r\n\r\n    addEventListener(\"fetch\", (event) => {\r\n        event.respondWith(handleRequest(event.request));\r\n    });\r\n}",
// 	  "json_args": {
// 		"key": "value"
// 	  },
// 	  "function_to_run": "",
// 	  "initiator_type": "edge_application",
// 	  "active": true,
// 	  "last_editor": "mail@mail.com",
// 	  "modified": "2023-04-27T17:37:12.389389Z",
// 	  "reference_count": 1
// 	},
// 	"schema_version": 3
//   }
// `

// var sucRespDomain string = `{
//     "results": {
//         "id": 1702659986,
//         "name": "My Domain",
//         "cnames": [],
//         "cname_access_only": false,
//         "digital_certificate_id": null,
//         "edge_application_id": 1697666970,
//         "is_active": true,
//         "domain_name": "ja65r2loc3.map.azionedge.net",
//         "is_mtls_enabled": false,
//         "mtls_verification": "enforce",
//         "mtls_trusted_ca_certificate_id": null
//     },
//     "schema_version": 3
// }`

// var sucRespOrigin string = `{
//     "results": {
//         "origin_id": 116207,
//         "origin_key": "35e3a635-2227-4bb6-976c-5e8c8fa58a67",
//         "name": "Create Origin22",
//         "origin_type": "single_origin",
//         "addresses": [
//             {
//                 "address": "httpbin.org",
//                 "weight": null,
//                 "server_role": "primary",
//                 "is_active": true
//             }
//         ],
//         "origin_protocol_policy": "http",
//         "is_origin_redirection_enabled": false,
//         "host_header": "${host}",
//         "method": "",
//         "origin_path": "/requests",
//         "connection_timeout": 60,
//         "timeout_between_bytes": 120,
//         "hmac_authentication": false,
//         "hmac_region_name": "",
//         "hmac_access_key": "",
//         "hmac_secret_key": ""
//     },
//     "schema_version": 3
// }`

// var sucRespCacheSettings = `{
//     "results": {
//         "id": 138708,
//         "name": "Default Cache Settings2234",
//         "browser_cache_settings": "override",
//         "browser_cache_settings_maximum_ttl": 0,
//         "cdn_cache_settings": "override",
//         "cdn_cache_settings_maximum_ttl": 60,
//         "cache_by_query_string": "ignore",
//         "query_string_fields": null,
//         "enable_query_string_sort": false,
//         "cache_by_cookies": "ignore",
//         "cookie_names": null,
//         "adaptive_delivery_action": "ignore",
//         "device_group": [],
//         "enable_caching_for_post": false,
//         "l2_caching_enabled": false,
//         "is_slice_configuration_enabled": false,
//         "is_slice_edge_caching_enabled": false,
//         "is_slice_l2_caching_enabled": false,
//         "slice_configuration_range": 1024,
//         "enable_caching_for_options": false,
//         "enable_stale_cache": true,
//         "l2_region": null
//     },
//     "schema_version": 3
// }`

// var sucRespRules = `{
//     "results": {
//         "id": 214790,
//         "name": "testorigin2",
//         "phase": "request",
//         "behaviors": [
//             {
//                 "name": "set_origin",
//                 "target": "116207"
//             }
//         ],
//         "criteria": [
//             [
//                 {
//                     "variable": "${uri}",
//                     "operator": "starts_with",
//                     "conditional": "if",
//                     "input_value": "/"
//                 }
//             ]
//         ],
//         "is_active": true,
//         "order": 1,
//         "description": ""
//     },
//     "schema_version": 3
// }`

func TestDeployCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	// t.Run("full flow manifest empty", func(t *testing.T) {
	// 	mock := &httpmock.Registry{}

	// 	options := &contracts.AzionApplicationOptions{
	// 		Name:   "LovelyName",
	// 		Bucket: "LovelyName",
	// 	}

	// 	dat, _ := os.ReadFile("./fixtures/create_app.json")
	// 	_ = json.Unmarshal(dat, options)

	// 	mock.Register(
	// 		httpmock.REST("POST", "edge_applications"),
	// 		httpmock.JSONFromString(successResponseApp),
	// 	)

	// 	mock.Register(
	// 		httpmock.REST("POST", "v4/storage/buckets"),
	// 		httpmock.JSONFromString(""),
	// 	)

	// 	mock.Register(
	// 		httpmock.REST("POST", "edge_functions"),
	// 		httpmock.JSONFromString(sucRespFunc),
	// 	)

	// 	mock.Register(
	// 		httpmock.REST("POST", "edge_applications/1697666970/origins"),
	// 		httpmock.JSONFromString(successRespOrigin),
	// 	)

	// 	mock.Register(
	// 		httpmock.REST("POST", "edge_applications/1697666970/functions_instances"),
	// 		httpmock.JSONFromString(sucRespInst),
	// 	)

	// 	mock.Register(
	// 		httpmock.REST("POST", "domains"),
	// 		httpmock.JSONFromString(sucRespDomain),
	// 	)

	// 	mock.Register(
	// 		httpmock.REST("PATCH", "edge_applications/1697666970"),
	// 		httpmock.JSONFromString(successResponseApp),
	// 	)

	// 	mock.Register(
	// 		httpmock.REST("POST", "edge_applications/1697666970/rules_engine/response/rules"),
	// 		httpmock.JSONFromString(successRespRule),
	// 	)

	// 	f, _, _ := testutils.NewFactory(mock)
	// 	deployCmd := NewDeployCmd(f)

	// 	deployCmd.FilepathWalk = func(root string, fn filepath.WalkFunc) error {
	// 		return nil
	// 	}
	// 	deployCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions) error {
	// 		return nil
	// 	}

	// 	deployCmd.FileReader = func(path string) ([]byte, error) {
	// 		return []byte{}, nil
	// 	}

	// 	deployCmd.Unmarshal = func(data []byte, v interface{}) error {
	// 		return nil
	// 	}

	// 	deployCmd.GetAzionJsonContent = func() (*contracts.AzionApplicationOptions, error) {
	// 		return &contracts.AzionApplicationOptions{}, nil
	// 	}

	// 	deployCmd.Interpreter = func() *manifestInt.ManifestInterpreter {
	// 		return &manifestInt.ManifestInterpreter{
	// 			FileReader: func(path string) ([]byte, error) {
	// 				return []byte{'{', '}'}, nil
	// 			},
	// 			WriteAzionJsonContent: func(conf *contracts.AzionApplicationOptions) error {
	// 				return nil
	// 			},
	// 			GetWorkDir: func() (string, error) {
	// 				return "", nil
	// 			},
	// 		}
	// 	}

	// 	err := deployCmd.Run(f)

	// 	// err := manifest.Interpreted(f, deployCmd, options, clients)
	// 	require.NoError(t, err)
	// })

	t.Run("without azion.json", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		deployCmd := NewDeployCmd(f)

		deployCmd.FileReader = func(path string) ([]byte, error) {
			return nil, os.ErrNotExist
		}

		cmd := NewCobraCmd(deployCmd)

		cmd.SetArgs([]string{""})

		err := cmd.Execute()

		require.Error(t, err)
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
			httpmock.REST("PATCH", "edge_applications/1697666970"),
			httpmock.JSONFromString(successResponseApp),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1697666970/functions_instances"),
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
