package dryrun

import (
	"bytes"
	"errors"
	"io/fs"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func fileReaderSuccess(path string) ([]byte, error) {
	// Simulated manifest content (as bytes)
	manifest := `{
		"origin": [
			{
				"name": "origin-storage-default",
				"origin_type": "object_storage",
				"prefix": ""
			}
		],
		"cache": [],
		"domain": {},
		"purge": [],
		"rules": [
			{
				"name": "Set Storage Origin for All Requests",
				"phase": "request",
				"description": "",
				"is_active": true,
				"order": 2,
				"criteria": [
					[
						{
							"variable": "${uri}",
							"operator": "matches",
							"conditional": "if",
							"input_value": "^\\\\/"
						}
					]
				],
				"behaviors": [
					{
						"name": "set_origin",
						"target": "origin-storage-default"
					}
				]
			},
			{
				"name": "Deliver Static Assets",
				"phase": "request",
				"description": "",
				"is_active": true,
				"order": 3,
				"criteria": [
					[
						{
							"variable": "${uri}",
							"operator": "matches",
							"conditional": "if",
							"input_value": ".(css|js|ttf|woff|woff2|pdf|svg|jpg|jpeg|gif|bmp|png|ico|mp4|json|xml|html)$"
						}
					]
				],
				"behaviors": [
					{
						"name": "set_origin",
						"target": "origin-storage-default"
					},
					{
						"name": "deliver",
						"target": null
					}
				]
			}
		]
	}`

	// Convert string manifest to byte slice and return it
	return []byte(manifest), nil
}

func TestSimulateDeploy(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name            string
		conf            *contracts.AzionApplicationOptions
		getAzionJsonErr error
		manifestPath    string
		manifestReadErr error
		wantError       bool
		wantMessages    []string
		workingDir      string
		fileReaderFunc  func(path string) ([]byte, error)
	}{
		{
			name: "successful simulation",
			conf: &contracts.AzionApplicationOptions{
				Name: "test-app",
				Application: contracts.AzionJsonDataApplication{
					ID: 12345,
				},
				Domain: contracts.AzionJsonDataDomain{
					Id: 12345,
				},
				Function: contracts.AzionJsonDataFunction{
					ID: 123321,
				},
				NotFirstRun: false,
			},
			manifestPath: "path/to/manifest.json",
			wantMessages: []string{
				"Updating Edge Application ",
				"Creating single Origin named 'test-app_single'",
			},
			fileReaderFunc: fileReaderSuccess,
		},
		{
			name:            "missing Azion JSON content",
			getAzionJsonErr: errors.New("file not found"),
			wantError:       true,
		},
		{
			name: "manifest path not found",
			conf: &contracts.AzionApplicationOptions{
				Name: "test-app",
			},
			manifestPath: "",
			wantMessages: []string{
				"This project has not been built yet",
			},
			fileReaderFunc: fileReaderSuccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := &httpmock.Registry{}

			f, _, _ := testutils.NewFactory(mock)

			dryrun := &DryrunStruct{
				Io: f.IOStreams,
				GetAzionJsonContent: func(path string) (*contracts.AzionApplicationOptions, error) {
					return tt.conf, tt.getAzionJsonErr
				},
				Interpreter: func() *manifestInt.ManifestInterpreter {
					return &manifestInt.ManifestInterpreter{
						GetWorkDir: utils.GetWorkingDir,
						FileReader: tt.fileReaderFunc,
					}
				},
				VersionID: func() string {
					return "v1.0.0"
				},
				F: f,
				Stat: func(name string) (fs.FileInfo, error) {
					return nil, nil
				},
			}

			err := dryrun.SimulateDeploy(tt.workingDir, tt.manifestPath)

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			output := f.IOStreams.Out.(*bytes.Buffer).String()
			for _, msg := range tt.wantMessages {
				assert.Contains(t, output, msg)
			}
		})
	}
}
