package applications

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name     string
		args     []string
		response httpmock.Responder
		contains []string
		err      error
	}{
		{
			name:     "listing with success",
			response: httpmock.JSONFromFile("./fixtures/response.json"),
			contains: []string{"ID", "NAME", "ACTIVE", "1337", "my-app", "true"},
		},
		{
			name:     "listing with details",
			args:     []string{"--details"},
			response: httpmock.JSONFromFile("./fixtures/response.json"),
			contains: []string{"LAST EDITOR", "LAST MODIFIED", "DEBUG", "user@azion.com", "2026-10-01"},
		},
		{
			name:     "no items",
			response: httpmock.JSONFromFile("./fixtures/no_items.json"),
			contains: []string{"ID", "NAME", "ACTIVE"},
		},
		{
			name:     "json invalid",
			response: httpmock.JSONFromString("{'name': 'some name',}"),
			err:      errors.New("Failed to list your Applications: invalid character '\\'' looking for beginning of object key string. Check your settings and try again. If the error persists, contact Azion support."),
		},
		{
			name:     "unauthorized",
			response: httpmock.StatusStringResponse(http.StatusUnauthorized, "{}"),
			err:      utils.ErrorToken401,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(httpmock.REST("GET", "workspace/applications"), tt.response)

			f, out, _ := testutils.NewFactory(mock)

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			if tt.err != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.err.Error())
				return
			}
			require.NoError(t, err)
			for _, s := range tt.contains {
				assert.Contains(t, out.String(), s)
			}
		})
	}
}
