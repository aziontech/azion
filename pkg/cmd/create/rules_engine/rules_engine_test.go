package rules_engine

import (
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/create/rules_engine"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	type apiMock struct {
		method, url, path string
	}

	tests := []struct {
		name      string
		args      []string
		mock      apiMock
		wantOut   string
		wantError error
		err       bool
	}{
		{
			name: "success phase request",
			args: []string{"--application-id", "1679423488", "--phase", "request", "--in", "./fixtures/create.json"},
			mock: apiMock{
				method: "POST",
				url:    "edge_applications/1679423488/rules_engine/request/rules",
				path:   "./fixtures/resp_phase_request.json",
			},
			wantOut:   fmt.Sprintf(rules_engine.OutputSuccess, 210543),
			wantError: nil,
			err:       false,
		},
		{
			name: "success phase response",
			args: []string{"--application-id", "1679423488", "--phase", "response", "--in", "./fixtures/create.json"},
			mock: apiMock{
				method: "POST",
				url:    "edge_applications/1679423488/rules_engine/response/rules",
				path:   "./fixtures/resp_phase_response.json",
			},
			wantOut:   fmt.Sprintf(rules_engine.OutputSuccess, 210544),
			wantError: nil,
			err:       false,
		},
		{
			name: "error name empty",
			args: []string{"--application-id", "1679423488", "--phase", "response", "--in", "./fixtures/create_name_empty.json"},
			mock: apiMock{
				method: "POST",
				url:    "edge_applications/1679423488/rules_engine/response/rules",
				path:   "./fixtures/resp_phase_response.json",
			},
			wantError: rules_engine.ErrorNameEmpty,
			err:       true,
		},
		{
			name: "error conditional empty",
			args: []string{"--application-id", "1679423488", "--phase", "response", "--in", "./fixtures/create_conditional_empty.json"},
			mock: apiMock{
				method: "POST",
				url:    "edge_applications/1679423488/rules_engine/response/rules",
				path:   "./fixtures/resp_phase_response.json",
			},
			wantError: rules_engine.ErrorConditionalEmpty,
			err:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}

			mock.Register(
				httpmock.REST(tt.mock.method, tt.mock.url),
				httpmock.JSONFromFile(tt.mock.path),
			)

			f, outGot, _ := testutils.NewFactory(mock)

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()

			if !tt.err && err == nil {
				require.Equal(t, tt.wantOut, outGot.String())
			} else {
				require.ErrorIs(t, err, tt.wantError)
			}
		})
	}
}
