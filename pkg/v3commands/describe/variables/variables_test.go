package variables

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
		expectErr bool
	}{
		{
			name: "describe a variable",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "variables/32e8ffca-4021-49a4-971f-330935566af4"),
					httpmock.JSONFromFile(".fixtures/variables.json"),
				)
			},
			args:      []string{"--variable-id", "32e8ffca-4021-49a4-971f-330935566af4"},
			expectErr: false,
		},
		{
			name: "not found",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "variables/32e8ffca-4021-49a4-971f-330935566af4"),
					httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
				)
			},
			args:      []string{"--variable-id", "32e8ffca-4021-49a4-971f-330935566af4"},
			expectErr: true,
		},
		{
			name: "no id sent",
			setupMock: func(mock *httpmock.Registry) {
				mock.Register(
					httpmock.REST("GET", "variables/123423424"),
					httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
				)
			},
			args:      []string{"--variable-id", "123423424"},
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

			cmd := NewCmd(f)
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
