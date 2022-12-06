package webapp

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	buildcmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/build"
	initcmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/init"
	publishcmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/publish"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/require"
)

var successResponseRule string = `
{
	"schema_version": 3,
	"results": {
		"id": 137056,
		"name": "Default Rule",
		"phase": "request",
		"behaviors": [
		  {
			"name": "run_function",
			"target": "6597"
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
		"order": 1
	  }
  }
`

var successResponseRules string = `
{
	"count": 1,
	"total_pages": 1,
	"schema_version": 3,
	"links": {
	  "previous": null,
	  "next": null
	},
	"results": [
	  {
		"id": 137056,
		"name": "Default Rule",
		"phase": "request",
		"behaviors": [
		  {
			"name": "run_function",
			"target": "6597"
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
		"order": 1
	  }
	]
  }
  `

var successResponseApp string = `
{
	"results":{
	   "id":666,
	   "name":"New Edge Applicahvjgjhgjhhgtion",
	   "delivery_protocol":"http",
	   "http_port":80,
	   "https_port":443,
	   "minimum_tls_version":"",
	   "active":true,
	   "application_acceleration":false,
	   "caching":true,
	   "device_detection":false,
	   "edge_firewall":false,
	   "edge_functions":false,
	   "image_optimization":false,
	   "load_balancer":false,
	   "raw_logs":false,
	   "web_application_firewall":false
	},
	"schema_version":3
}
`

var successResponseFunc string = `
{
	"results": {
	  "id": 6597,
	  "name": "imimimi",
	  "language": "javascript",
	  "code": "async function handleRequest(request) {\n    return new Response(\"Hello World!\",\n      {\n          status:204\n      })\n   }\n   addEventListener(\"fetch\", event => {\n    event.respondWith(handleRequest(event.request))\n   })",
	  "json_args": {},
	  "function_to_run": "",
	  "initiator_type": "edge_application",
	  "active": true,
	  "last_editor": "patrickmenott@gmail.com",
	  "modified": "2022-06-30T19:26:17.003242Z",
	  "reference_count": 0
	},
	"schema_version": 3
  }
  `

var successResponseDom string = `
  {
	"results": {
	  "id": 1666281562,
	  "name": "teste",
	  "cnames": [],
	  "cname_access_only": false,
	  "digital_certificate_id": null,
	  "edge_application_id": 666,
	  "is_active": true,
	  "domain_name": "66642069.map.azionedge.net"
	},
	"schema_version": 3
  }
  `

func TestWebappCmd(t *testing.T) {
	t.Run("nextjs testing", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions"),
			httpmock.JSONFromString(successResponseFunc),
		)

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

		mock.Register(
			httpmock.REST("GET", "edge_applications/666/rules_engine/request/rules"),
			httpmock.JSONFromString(successResponseRules),
		)

		mock.Register(
			httpmock.REST("PATCH", "edge_applications/666/rules_engine/request/rules/137056"),
			httpmock.JSONFromString(successResponseRule),
		)

		mock.Register(
			httpmock.REST("POST", "domains"),
			httpmock.JSONFromString(successResponseDom),
		)

		options := &contracts.AzionApplicationOptions{}

		f, _, _ := testutils.NewFactory(mock)
		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		buildEnvs := []string{"AWS_ACCESS_KEY_ID=123456789", "AWS_SECRET_ACCESS_KEY=987654321"}

		// INIT

		initCmd := initcmd.NewInitCmd(f)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}

		initCmd.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
		}

		initCmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}

		initCmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		initCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}

		cmdInit := initcmd.NewCobraCmd(initCmd)

		cmdInit.SetArgs([]string{"--name", "functional_testing_nextjs", "--type", "nextjs"})
		err := cmdInit.Execute()

		//BUILD

		buildCommand := buildcmd.NewBuildCmd(f)

		buildCommand.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"build": {"cmd": "./azion/webdev.sh build", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
		}

		buildCommand.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "Build completed", 0, nil
		}

		buildCommand.EnvLoader = func(path string) ([]string, error) {
			return buildEnvs, nil
		}

		buildCommand.GetWorkDir = func() (string, error) {
			return "", nil
		}

		cmdBuild := buildcmd.NewCobraCmd(buildCommand)

		errBuild := cmdBuild.Execute()

		//PUBLISH

		publishCmd := publishcmd.NewPublishCmd(f)

		publishCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"pre_cmd": "./azion/webdev.sh publish", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "nextjs"}`), nil
		}

		publishCmd.EnvLoader = func(path string) ([]string, error) {
			return buildEnvs, nil
		}

		publishCmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "", 0, nil
		}

		publishCmd.BuildCmd = func(f *cmdutil.Factory) *buildcmd.BuildCmd {
			return buildCommand
		}

		publishCmd.GetAzionJsonContent = func() (*contracts.AzionApplicationOptions, error) {
			return options, nil
		}

		publishCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions) error {
			return nil
		}

		cmdPublish := publishcmd.NewCobraCmd(publishCmd)

		errPublish := cmdPublish.Execute()
		fmt.Println(errPublish)

		require.NoError(t, err)

		require.NoError(t, errBuild)

		require.NoError(t, errPublish)

	})

	t.Run("flareact testing", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_functions"),
			httpmock.JSONFromString(successResponseFunc),
		)

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

		mock.Register(
			httpmock.REST("GET", "edge_applications/666/rules_engine/request/rules"),
			httpmock.JSONFromString(successResponseRules),
		)

		mock.Register(
			httpmock.REST("PATCH", "edge_applications/666/rules_engine/request/rules/137056"),
			httpmock.JSONFromString(successResponseRule),
		)

		mock.Register(
			httpmock.REST("POST", "domains"),
			httpmock.JSONFromString(successResponseDom),
		)

		options := &contracts.AzionApplicationOptions{}

		f, _, _ := testutils.NewFactory(mock)
		envs := []string{"UEBA=OBA", "FAZER=UM_PENSO"}
		buildEnvs := []string{"AWS_ACCESS_KEY_ID=123456789", "AWS_SECRET_ACCESS_KEY=987654321"}

		// INIT

		initCmd := initcmd.NewInitCmd(f)

		initCmd.LookPath = func(bin string) (string, error) {
			return "", nil
		}

		initCmd.WriteFile = func(filename string, data []byte, perm fs.FileMode) error {
			return nil
		}

		initCmd.Rename = func(oldpath string, newpath string) error {
			return nil
		}

		initCmd.Stat = func(path string) (fs.FileInfo, error) {
			return nil, nil
		}

		initCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"init": {"cmd": "ls", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "flareact"}`), nil
		}

		initCmd.EnvLoader = func(path string) ([]string, error) {
			return envs, nil
		}

		initCmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "my command output", 0, nil
		}

		initCmd.Mkdir = func(path string, perm os.FileMode) error {
			return nil
		}

		initCmd.GitPlainClone = func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error) {
			return &git.Repository{}, nil
		}

		cmdInit := initcmd.NewCobraCmd(initCmd)

		cmdInit.SetArgs([]string{"--name", "functional_testing_flareact", "--type", "flareact"})
		err := cmdInit.Execute()

		//BUILD

		buildCommand := buildcmd.NewBuildCmd(f)

		buildCommand.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"build": {"cmd": "./azion/webdev.sh build", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "flareact"}`), nil
		}

		buildCommand.CommandRunner = func(cmd string, envs []string) (string, int, error) {
			return "Build completed", 0, nil
		}

		buildCommand.EnvLoader = func(path string) ([]string, error) {
			return buildEnvs, nil
		}

		buildCommand.GetWorkDir = func() (string, error) {
			return "", nil
		}

		cmdBuild := buildcmd.NewCobraCmd(buildCommand)

		errBuild := cmdBuild.Execute()

		//PUBLISH

		publishCmd := publishcmd.NewPublishCmd(f)

		publishCmd.FileReader = func(path string) ([]byte, error) {
			return []byte(`{"publish": {"pre_cmd": "./azion/webdev.sh publish", "env": "./azion/init.env", "output-ctrl": "on-error"}, "type": "flareact"}`), nil
		}

		publishCmd.EnvLoader = func(path string) ([]string, error) {
			return buildEnvs, nil
		}

		publishCmd.CommandRunner = func(cmd string, env []string) (string, int, error) {
			return "", 0, nil
		}

		publishCmd.BuildCmd = func(f *cmdutil.Factory) *buildcmd.BuildCmd {
			return buildCommand
		}

		publishCmd.GetAzionJsonContent = func() (*contracts.AzionApplicationOptions, error) {
			return options, nil
		}

		publishCmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions) error {
			return nil
		}

		cmdPublish := publishcmd.NewCobraCmd(publishCmd)

		errPublish := cmdPublish.Execute()
		fmt.Println(errPublish)

		require.NoError(t, err)

		require.NoError(t, errBuild)

		require.NoError(t, errPublish)

	})
}
