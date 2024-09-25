package deploy

import (
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MockFactory struct {
	cmdutil.Factory
}

type MockDeployCmd struct {
	DeployCmd
}

type MockIOStreams struct {
	cmdutil.Factory
}

// MockFileReader mocks os.ReadFile behavior
func MockCheckToken(f *cmdutil.Factory) error {
	return nil
}

// MockFileReader mocks os.ReadFile behavior
func MockFileReader(path string) ([]byte, error) {
	if path == "validPath" {
		return []byte(`{"valid": "json"}`), nil
	}
	return nil, errors.New("file not found")
}

func MockReadSettings() (token.Settings, error) {
	return token.Settings{Token: "123321", S3AccessKey: "122221", S3SecreKey: "3333322222", S3Bucket: "bucketname"}, nil
}

// MockWriteFile mocks os.WriteFile behavior
func MockWriteFile(filename string, data []byte, perm fs.FileMode) error {
	if filename == "validWritePath" {
		return nil
	}
	return errors.New("write failed")
}

// MockGetWorkDir mocks utils.GetWorkingDir
func MockGetWorkDir() (string, error) {
	return "/mocked/path", nil
}

// MockOpen mocks os.Open
func MockOpen(name string) (*os.File, error) {
	return nil, errors.New("open failed")
}

// MockFilepathWalk mocks filepath.Walk
func MockFilepathWalk(root string, fn filepath.WalkFunc) error {
	return nil
}

// MockGetAzionJsonContent mocks GetAzionJsonContent
func MockGetAzionJsonContent(pathConfig string) (*contracts.AzionApplicationOptions, error) {
	return &contracts.AzionApplicationOptions{
		Name:   "MockApp",
		Prefix: "MockPrefix",
	}, nil
}

// MockWriteAzionJsonContent mocks WriteAzionJsonContent
func MockWriteAzionJsonContent(conf *contracts.AzionApplicationOptions, confConf string) error {
	return nil
}

// MockCallScript mocks the callScript behavior
func MockCallScript(token string, id string, secret string, prefix string, name string, cmd *DeployCmd) (string, error) {
	return "mockExecId", nil
}

// MockOpenBrowser mocks openBrowser behavior
func MockOpenBrowser(f *cmdutil.Factory, urlConsoleDeploy string, cmd *DeployCmd) error {
	return nil
}

// MockCaptureLogs mocks captureLogs behavior
func MockCaptureLogs(execId string, token string, cmd *DeployCmd) error {
	return nil
}

func TestDeploy_Run(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name        string
		setupMocks  func(cmd *DeployCmd)
		expectedErr error
	}{
		{
			name: "successful deploy",
			setupMocks: func(cmd *DeployCmd) {
				cmd.GetWorkDir = MockGetWorkDir
				cmd.FileReader = MockFileReader
				cmd.WriteFile = MockWriteFile
				cmd.GetAzionJsonContent = MockGetAzionJsonContent
				cmd.WriteAzionJsonContent = MockWriteAzionJsonContent
				cmd.CallScript = MockCallScript
				cmd.OpenBrowser = MockOpenBrowser
				cmd.CaptureLogs = MockCaptureLogs
				cmd.CheckToken = MockCheckToken
				cmd.ReadSettings = func() (token.Settings, error) {
					return token.Settings{}, nil
				}
				cmd.UploadFiles = func(f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string, pathStatic, bucket string, cmd *DeployCmd) error {
					return nil
				}
			},
			expectedErr: nil,
		},
		{
			name: "fail on GetWorkDir",
			setupMocks: func(cmd *DeployCmd) {
				cmd.GetWorkDir = func() (string, error) { return "", errors.New("workdir error") }
				cmd.FileReader = MockFileReader
				cmd.WriteFile = MockWriteFile
				cmd.GetAzionJsonContent = MockGetAzionJsonContent
				cmd.WriteAzionJsonContent = MockWriteAzionJsonContent
				cmd.CallScript = MockCallScript
				cmd.OpenBrowser = MockOpenBrowser
				cmd.CaptureLogs = MockCaptureLogs
				cmd.CheckToken = MockCheckToken
				cmd.ReadSettings = MockReadSettings
				cmd.UploadFiles = func(f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string, pathStatic, bucket string, cmd *DeployCmd) error {
					return nil
				}
			},
			expectedErr: errors.New("workdir error"),
		},
		{
			name: "fail on CallScript",
			setupMocks: func(cmd *DeployCmd) {
				cmd.GetWorkDir = MockGetWorkDir
				cmd.FileReader = MockFileReader
				cmd.WriteFile = MockWriteFile
				cmd.GetAzionJsonContent = MockGetAzionJsonContent
				cmd.WriteAzionJsonContent = MockWriteAzionJsonContent
				cmd.CheckToken = MockCheckToken
				cmd.ReadSettings = MockReadSettings
				cmd.UploadFiles = func(f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string, pathStatic, bucket string, cmd *DeployCmd) error {
					return nil
				}
				cmd.CallScript = func(token string, id string, secret string, prefix string, name string, cmd *DeployCmd) (string, error) {
					return "", errors.New("call script failed")
				}
				cmd.OpenBrowser = MockOpenBrowser
				cmd.CaptureLogs = MockCaptureLogs
			},
			expectedErr: errors.New("call script failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(
				httpmock.REST(http.MethodPost, "v4/storage/buckets"),
				httpmock.JSONFromFile("fixtures/response.json"),
			)
			mock.Register(
				httpmock.REST(http.MethodPost, "v4/storage/s3-credentials"),
				httpmock.JSONFromFile("fixtures/responses3.json"),
			)

			f, _, _ := testutils.NewFactory(mock)

			deployCmd := NewDeployCmd(f)
			tt.setupMocks(deployCmd)

			err := deployCmd.Run(f)
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCaptureLogs(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name            string
		logsResponse    string
		resultsResponse string
		logStruct       contracts.Logs
		resultStruct    contracts.Results
		expectError     bool
	}{
		{
			name:            "Successful deployment",
			logsResponse:    `{"status": "succeeded"}`,
			resultsResponse: `{"result": {"azion": {"key": "value"}, "errors": null}}`,
			expectError:     false,
			logStruct:       contracts.Logs{Status: "succeeded"},
		},
		{
			name:            "Deployment with errors",
			logsResponse:    `{"status": "succeeded"}`,
			resultsResponse: `{"result": {"errors": {"stack": "Something went wrong"}}}`,
			logStruct:       contracts.Logs{Status: "succeeded"},
			resultStruct: contracts.Results{Result: contracts.Result{
				Errors: &contracts.ErrorDetails{
					Stack:   "Something went wrong",
					Message: "Something went wrong",
				},
			}},

			expectError: true,
		},
		{
			name:            "Log status running, succeeds on retry",
			logsResponse:    `{"status": "running"}`,
			resultsResponse: `{"result": {"azion": {"key": "value"}, "errors": null}}`,
			expectError:     false,
			logStruct:       contracts.Logs{Status: "succeeded"},
		},
		{
			name:         "Unknown log status",
			logsResponse: `{"status": "salabim"}`,
			expectError:  true,
			logStruct:    contracts.Logs{Status: "salabim"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(
				httpmock.REST(http.MethodGet, "api/script-runner/executions/e89b32a6-c912-4fba-bae7-1e7ff115256f/logs"),
				httpmock.JSONFromString(tt.logsResponse),
			)
			mock.Register(
				httpmock.REST(http.MethodGet, "api/script-runner/executions/e89b32a6-c912-4fba-bae7-1e7ff115256f/results"),
				httpmock.JSONFromString(tt.resultsResponse),
			)

			f, _, _ := testutils.NewFactory(mock)

			// Run the function
			token := "test-token"
			execID := "e89b32a6-c912-4fba-bae7-1e7ff115256f"
			cmd := NewDeployCmd(f)
			Logs = tt.logStruct
			Result = tt.resultStruct

			cmd.WriteAzionJsonContent = func(conf *contracts.AzionApplicationOptions, confConf string) error {
				return nil
			}
			err := cmd.CaptureLogs(execID, token, cmd)

			// Check for expected error
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestNewCobraCommand(t *testing.T) {
	t.Run("Test NewCobraCommand", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)

		deployCmd := NewDeployCmd(f)
		cmd := NewCobraCmd(deployCmd)

		assert.NotNil(t, cmd)

	})
}

func TestNewCmd(t *testing.T) {
	t.Run("Test NewCobraCommand", func(t *testing.T) {
		f, _, _ := testutils.NewFactory(nil)
		cmd := NewCmd(f)
		assert.NotNil(t, cmd)

	})
}

func TestCallScript(t *testing.T) {
	logger.New(zap.DebugLevel)

	// Define the test cases
	tests := []struct {
		name            string
		token           string
		id              string
		secret          string
		prefix          string
		projectName     string
		resultsResponse string
		expectedUUID    string
		expectedErr     bool
	}{
		{
			name:            "Successful API Call",
			token:           "validToken",
			id:              "validID",
			secret:          "validSecret",
			prefix:          "validPrefix",
			projectName:     "myProject",
			resultsResponse: `{"uuid": "12345"}`,
			expectedUUID:    "12345",
			expectedErr:     false,
		},
		{
			name:            "Error in API Response",
			token:           "invalidToken",
			id:              "validID",
			secret:          "validSecret",
			prefix:          "validPrefix",
			projectName:     "myProject",
			resultsResponse: "",
			expectedUUID:    "",
			expectedErr:     true,
		},
	}

	// Run each test case
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new HTTP mock registry
			mock := &httpmock.Registry{}

			// Register the mock API responses for logs and results
			mock.Register(
				httpmock.REST(http.MethodPost, "/api/template-engine/templates/17ac912d-5ce9-4806-9fa7-480779e43f58/instantiate"),
				httpmock.JSONFromString(tt.resultsResponse),
			)

			f, _, _ := testutils.NewFactory(mock)
			cmd := NewDeployCmd(f)
			if tt.expectedErr {
				cmd.Unmarshal = func(data []byte, v interface{}) error {
					return errors.New("error")
				}
			}

			// Call the function under test
			_, err := cmd.CallScript(tt.token, tt.id, tt.secret, tt.prefix, tt.projectName, cmd)

			// Assert the results
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckToken(t *testing.T) {
	tests := []struct {
		name        string
		tokenValue  string
		validToken  bool
		validateErr error
		expectError error
	}{
		{
			name:        "Invalid token",
			tokenValue:  "invalid-token",
			validToken:  false,
			expectError: utils.ErrorTokenNotProvided,
		},
		{
			name:        "Token validation fails",
			tokenValue:  "test-token",
			validateErr: errors.New("Token was not provided; the CLI uses a previous stored token if it was configured. You must provide a valid token, or create a new one, and configure it to use with the CLI. Manage your personal tokens on RTM using the Account Menu > Personal Tokens and configure the token again with the command 'azion -t <token>'"),
			expectError: errors.New("Token was not provided; the CLI uses a previous stored token if it was configured. You must provide a valid token, or create a new one, and configure it to use with the CLI. Manage your personal tokens on RTM using the Account Menu > Personal Tokens and configure the token again with the command 'azion -t <token>'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, _ := testutils.NewFactory(nil)

			err := checkToken(f)

			if tt.expectError != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.expectError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOpenBrowser(t *testing.T) {
	logger.New(zap.DebugLevel)
	tests := []struct {
		name             string
		urlConsoleDeploy string
		openBrowserErr   error
		expectedLogMsg   string
		expectedErr      error
		openFunc         func(input string) error
	}{
		{
			name:             "successful browser open",
			urlConsoleDeploy: "http://example.com",
			openBrowserErr:   nil,
			expectedLogMsg:   "Please visit the following URL: http://example.com",
			expectedErr:      nil,
			openFunc: func(input string) error {
				return nil
			},
		},
		{
			name:             "failed to open browser",
			urlConsoleDeploy: "http://example.com",
			openBrowserErr:   errors.New("failed to open browser"),
			expectedLogMsg:   "Please visit the following URL: http://example.com",
			expectedErr:      errors.New("failed to open browser"),
			openFunc: func(input string) error {
				return errors.New("failed to open browser")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, _ := testutils.NewFactory(nil)

			cmd := NewDeployCmd(f)
			cmd.OpenBrowserFunc = tt.openFunc

			// Call the openBrowser function with the mock factory and deploy command
			err := openBrowser(f, tt.urlConsoleDeploy, cmd)

			// Verify the error behavior
			if tt.expectedErr != nil {
				assert.EqualError(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
