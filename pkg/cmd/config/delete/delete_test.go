package delete

import (
	"errors"
	"io/fs"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/config/delete"
)

func init() {
	logger.New(zapcore.DebugLevel)
}

func mockConfirmYes(string) (string, error) {
	return "y", nil
}

func mockConfirmNo(string) (string, error) {
	return "n", nil
}

func mockGetAzionJson(confPath string) (*contracts.AzionApplicationOptions, error) {
	return &contracts.AzionApplicationOptions{
		Application: contracts.AzionJsonDataApplication{
			ID:   1234,
			Name: "test-application",
		},
		Function: []contracts.AzionJsonDataFunction{
			{
				ID:           5678,
				Name:         "test-function",
				InstanceID:   9012,
				InstanceName: "test-instance",
			},
		},
		RulesEngine: contracts.AzionJsonDataRulesEngine{
			Created: true,
			Rules: []contracts.AzionJsonDataRules{
				{
					Id:    3456,
					Name:  "test-rule",
					Phase: "request",
				},
			},
		},
		CacheSettings: []contracts.AzionJsonDataCacheSettings{
			{
				Id:   7890,
				Name: "test-cache-setting",
			},
		},
		Workloads: contracts.AzionJsonDataWorkload{
			Id:   1111,
			Name: "test-workload",
		},
		Bucket: "test-bucket",
		Firewalls: []contracts.AzionJsonDataFirewall{
			{
				Id:   2222,
				Name: "test-firewall",
				Rules: []contracts.AzionJsonDataFirewallRule{
					{
						Id:   3333,
						Name: "test-fw-rule",
					},
				},
			},
		},
		Connectors: []contracts.AzionJsonDataConnectors{
			{
				Id:   4444,
				Name: "test-connector",
			},
		},
	}, nil
}

// Mock without bucket for tests that don't need storage
func mockGetAzionJsonWithoutBucket(confPath string) (*contracts.AzionApplicationOptions, error) {
	return &contracts.AzionApplicationOptions{
		Application: contracts.AzionJsonDataApplication{
			ID:   1234,
			Name: "test-application",
		},
		Function: []contracts.AzionJsonDataFunction{
			{
				ID:           5678,
				Name:         "test-function",
				InstanceID:   9012,
				InstanceName: "test-instance",
			},
		},
		RulesEngine: contracts.AzionJsonDataRulesEngine{
			Created: true,
			Rules: []contracts.AzionJsonDataRules{
				{
					Id:    3456,
					Name:  "test-rule",
					Phase: "request",
				},
			},
		},
		CacheSettings: []contracts.AzionJsonDataCacheSettings{
			{
				Id:   7890,
				Name: "test-cache-setting",
			},
		},
		Workloads: contracts.AzionJsonDataWorkload{
			Id:   1111,
			Name: "test-workload",
		},
		Bucket: "", // No bucket
		Firewalls: []contracts.AzionJsonDataFirewall{
			{
				Id:   2222,
				Name: "test-firewall",
				Rules: []contracts.AzionJsonDataFirewallRule{
					{
						Id:   3333,
						Name: "test-fw-rule",
					},
				},
			},
		},
		Connectors: []contracts.AzionJsonDataConnectors{
			{
				Id:   4444,
				Name: "test-connector",
			},
		},
	}, nil
}

func mockGetEmptyAzionJson(confPath string) (*contracts.AzionApplicationOptions, error) {
	return &contracts.AzionApplicationOptions{}, nil
}

func mockGetAzionJsonError(confPath string) (*contracts.AzionApplicationOptions, error) {
	return nil, errors.New("error reading azion.json")
}

func mockGetAzionJsonNotFound(confPath string) (*contracts.AzionApplicationOptions, error) {
	return nil, errors.New("open azion/azion.json: no such file or directory")
}

func mockWriteFileSuccess(filename string, data []byte, perm fs.FileMode) error {
	return nil
}

func mockWriteFileError(filename string, data []byte, perm fs.FileMode) error {
	return errors.New("error writing file")
}

func mockGetWorkDirSuccess() (string, error) {
	return "/tmp", nil
}

func mockGetWorkDirError() (string, error) {
	return "", errors.New("error getting working directory")
}

func TestConfigDelete(t *testing.T) {
	tests := []struct {
		name           string
		force          bool
		mockAzion      func(confPath string) (*contracts.AzionApplicationOptions, error)
		mockAskInput   func(string) (string, error)
		mockWriteFile  func(filename string, data []byte, perm fs.FileMode) error
		mockWorkDir    func() (string, error)
		statusCode     int
		expectError    bool
		expectedOutput string
	}{
		{
			name:           "delete all resources with force flag",
			force:          true,
			mockAzion:      mockGetAzionJsonWithoutBucket,
			mockAskInput:   nil, // not needed with force
			mockWriteFile:  mockWriteFileSuccess,
			mockWorkDir:    mockGetWorkDirSuccess,
			statusCode:     204,
			expectError:    false,
			expectedOutput: msg.DeleteSuccess,
		},
		{
			name:           "delete with user confirmation",
			force:          false,
			mockAzion:      mockGetAzionJsonWithoutBucket,
			mockAskInput:   mockConfirmYes,
			mockWriteFile:  mockWriteFileSuccess,
			mockWorkDir:    mockGetWorkDirSuccess,
			statusCode:     204,
			expectError:    false,
			expectedOutput: msg.DeleteSuccess,
		},
		{
			name:           "user aborts deletion",
			force:          false,
			mockAzion:      mockGetAzionJson,
			mockAskInput:   mockConfirmNo,
			mockWriteFile:  mockWriteFileSuccess,
			mockWorkDir:    mockGetWorkDirSuccess,
			statusCode:     0,
			expectError:    false,
			expectedOutput: msg.DeletionAborted,
		},
		{
			name:           "no resources to delete",
			force:          true,
			mockAzion:      mockGetEmptyAzionJson,
			mockAskInput:   nil,
			mockWriteFile:  mockWriteFileSuccess,
			mockWorkDir:    mockGetWorkDirSuccess,
			statusCode:     0,
			expectError:    false,
			expectedOutput: "No resources found in azion.json to delete\n",
		},
		{
			name:           "azion.json not found",
			force:          true,
			mockAzion:      mockGetAzionJsonNotFound,
			mockAskInput:   nil,
			mockWriteFile:  mockWriteFileSuccess,
			mockWorkDir:    mockGetWorkDirSuccess,
			statusCode:     0,
			expectError:    true,
			expectedOutput: "",
		},
		{
			name:           "error reading azion.json",
			force:          true,
			mockAzion:      mockGetAzionJsonError,
			mockAskInput:   nil,
			mockWriteFile:  mockWriteFileSuccess,
			mockWorkDir:    mockGetWorkDirSuccess,
			statusCode:     0,
			expectError:    true,
			expectedOutput: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}

			// Register mock responses for various API calls
			if tt.statusCode == 204 {
				// List request phase rules - return the rule that needs to be deleted
				mock.Register(
					httpmock.REST("GET", "workspace/applications/1234/request_rules"),
					httpmock.StatusStringResponse(200, `{"count":1,"results":[{"id":3456,"name":"test-rule","active":true}]}`),
				)
				// List response phase rules - return empty
				mock.Register(
					httpmock.REST("GET", "workspace/applications/1234/response_rules"),
					httpmock.StatusStringResponse(200, `{"count":0,"results":[]}`),
				)
				// Rules Engine delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/applications/1234/request_rules/3456"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// List firewall rules
				mock.Register(
					httpmock.REST("GET", "workspace/firewalls/2222/request_rules"),
					httpmock.StatusStringResponse(200, `{"count":1,"results":[{"id":3333,"name":"test-fw-rule","active":true}]}`),
				)
				// Firewall rule delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/firewalls/2222/request_rules/3333"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// List function instances
				mock.Register(
					httpmock.REST("GET", "workspace/applications/1234/functions"),
					httpmock.StatusStringResponse(200, `{"count":1,"results":[{"id":9012,"name":"test-instance","function":5678}]}`),
				)
				// Function Instance delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/applications/1234/functions/9012"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// List cache settings
				mock.Register(
					httpmock.REST("GET", "workspace/applications/1234/cache_settings"),
					httpmock.StatusStringResponse(200, `{"count":1,"results":[{"id":7890,"name":"test-cache-setting"}]}`),
				)
				// Cache Setting delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/applications/1234/cache_settings/7890"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// Application delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/applications/1234"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// Function delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/functions/5678"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// Workload delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/workloads/1111"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// Firewall delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/firewalls/2222"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// Connector delete
				mock.Register(
					httpmock.REST("DELETE", "workspace/connectors/4444"),
					httpmock.StatusStringResponse(tt.statusCode, ""),
				)
				// Storage bucket delete - this is more complex, skip for now
			}

			f, stdout, _ := testutils.NewFactory(mock)

			deleteCmd := NewDeleteCmd(f)
			deleteCmd.GetAzion = tt.mockAzion
			if tt.mockAskInput != nil {
				deleteCmd.AskInput = tt.mockAskInput
			}
			deleteCmd.WriteFile = tt.mockWriteFile
			deleteCmd.GetWorkDir = tt.mockWorkDir

			cobraCmd := NewCobraCmd(deleteCmd)

			if tt.force {
				cobraCmd.SetArgs([]string{"--force"})
			}

			_, err := cobraCmd.ExecuteC()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Contains(t, stdout.String(), tt.expectedOutput)
			}
		})
	}
}

func TestResetAzionJson(t *testing.T) {
	tests := []struct {
		name          string
		mockWriteFile func(filename string, data []byte, perm fs.FileMode) error
		mockWorkDir   func() (string, error)
		expectError   bool
	}{
		{
			name:          "reset azion.json success",
			mockWriteFile: mockWriteFileSuccess,
			mockWorkDir:   mockGetWorkDirSuccess,
			expectError:   false,
		},
		{
			name:          "reset azion.json write error",
			mockWriteFile: mockWriteFileError,
			mockWorkDir:   mockGetWorkDirSuccess,
			expectError:   true,
		},
		{
			name:          "reset azion.json workdir error",
			mockWriteFile: mockWriteFileSuccess,
			mockWorkDir:   mockGetWorkDirError,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, _ := testutils.NewFactory(nil)

			deleteCmd := NewDeleteCmd(f)
			deleteCmd.WriteFile = tt.mockWriteFile
			deleteCmd.GetWorkDir = tt.mockWorkDir

			err := deleteCmd.resetAzionJson("azion")
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
