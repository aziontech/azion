package rulesengine

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name      string
		setupMock func(mock *httpmock.Registry)
		args      []string
		askInput  func(s string) (string, error)
		expectErr bool
	}{
		{
			name: "describe a rule engine",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617"},
			expectErr: false,
		},
		{
			name: "describe a rule engine - ask for app id",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args: []string{"--rule-id", "173617"},
			askInput: func(s string) (string, error) {
				return "1678743802", nil
			},
			expectErr: false,
		},
		{
			name: "describe a rule engine - ask for phase",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args: []string{"--application-id", "1678743802", "--rule-id", "173617"},
			askInput: func(s string) (string, error) {
				return "request", nil
			},
			expectErr: false,
		},
		{
			name: "describe a rule engine - ask for rule id",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args: []string{"--application-id", "1678743802"},
			askInput: func(s string) (string, error) {
				return "173617", nil
			},
			expectErr: false,
		},
		{
			name: "not found",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
				)
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "666"},
			expectErr: true,
		},
		{
			name: "missing mandatory flag",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/1"),
					httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
				)
			},
			args:      []string{},
			expectErr: true,
		},
		{
			name: "different phases",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617"},
			expectErr: false,
		},
		{
			name: "invalid JSON response",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.StringResponse("{invalid json"),
				)
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617"},
			expectErr: true,
		},
		{
			name: "non-existent application ID",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/999999999/rules/173617"),
					httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
				)
			},
			args:      []string{"--application-id", "999999999", "--rule-id", "173617"},
			expectErr: true,
		},
		{
			name: "invalid phase",
			setupMock: func(mock *httpmock.Registry) {
				// No mock needed for invalid phase
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			if tt.setupMock != nil {
				tt.setupMock(mock)
			}

			f, _, _ := testutils.NewFactory(mock)

			descCmd := NewDescribeCmd(f)
			if tt.askInput != nil {
				descCmd.AskInput = tt.askInput
			}

			cmd := NewCobraCmd(descCmd, f)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
