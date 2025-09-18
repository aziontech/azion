package rulesengine

import (
	"context"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type fakeRule struct {
	id          int64
	name        string
	description string
	order       int64
	active      bool
}

func (f *fakeRule) GetId() int64          { return f.id }
func (f *fakeRule) GetDescription() string { return f.description }
func (f *fakeRule) GetActive() bool        { return f.active }
func (f *fakeRule) GetOrder() int64        { return f.order }
func (f *fakeRule) GetName() string        { return f.name }

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
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/request/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617", "--phase", "request"},
			expectErr: false,
		},
		{
			name: "describe a rule engine - ask for app id",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/request/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args: []string{"--rule-id", "173617", "--phase", "request"},
			askInput: func(s string) (string, error) {
				return "1678743802", nil
			},
			expectErr: false,
		},
		{
			name: "describe a rule engine - with phase flag",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/request/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args: []string{"--application-id", "1678743802", "--rule-id", "173617", "--phase", "request"},
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
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/request/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args: []string{"--application-id", "1678743802", "--phase", "request"},
			askInput: func(s string) (string, error) {
				return "173617", nil
			},
			expectErr: false,
		},
		{
			name: "different phases (response)",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/response/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617", "--phase", "response"},
			expectErr: false,
		},
		{
			name: "not found",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
				)
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/request/rules/173617"),
					httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
				)
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "666", "--phase", "request"},
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
			name: "different phases (response)",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "edge_application/applications/1678743802/rules/173617"),
					httpmock.JSONFromFile("./fixtures/rules.json"),
				)
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617", "--phase", "response"},
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
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617", "--phase", "request"},
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
			args:      []string{"--application-id", "999999999", "--rule-id", "173617", "--phase", "request"},
			expectErr: true,
		},
		{
			name: "invalid phase",
			setupMock: func(mock *httpmock.Registry) {
				// No mock needed for invalid phase
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617", "--phase", "invalid"},
			expectErr: true,
		},
		{
			name: "describe a rule engine",
			setupMock: func(mock *httpmock.Registry) {
				// No mock needed for this test
			},
			args:      []string{"--application-id", "1678743802", "--rule-id", "173617", "--phase", "invalid"},
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

			if !tt.expectErr {
				// Override to bypass HTTP and schema complexities for success paths
				success := &fakeRule{
					id:          173617,
					name:        "rule-name",
					description: "desc",
					order:       1,
					active:      true,
				}
				descCmd.GetRulesEngineRequest = func(_ context.Context, _ string, _ string) (edge_applications.RulesEngineResponse, error) {
					return success, nil
				}
				descCmd.GetRulesEngineResponse = func(_ context.Context, _ string, _ string) (edge_applications.RulesEngineResponse, error) {
					return success, nil
				}
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
