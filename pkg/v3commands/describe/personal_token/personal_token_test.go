package personal_token

import (
	"net/http"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name     string
		request  httpmock.Matcher
		response httpmock.Responder
		args     []string
		output   string
		err      string
	}{
		{
			name:     "Success Status 200 Ok with struct nil no error",
			request:  httpmock.REST(http.MethodGet, "iam/personal_tokens/7b026645-e3dc-4d0f-91c5-6f4030996f6c"),
			response: httpmock.JSONFromString("{}"),
			args:     []string{"--id", "7b026645-e3dc-4d0f-91c5-6f4030996f6c"},
			output:   "Uuid:          null  \nName:          null  \nCreated:       null  \nExpires At:    null  \nDescription:   null  \n",
		},
		{
			name:     "Success Status 200 Ok with items",
			request:  httpmock.REST(http.MethodGet, "iam/personal_tokens/7b026645-e3dc-4d0f-91c5-6f4030996f6c"),
			response: httpmock.JSONFromFile(".fixtures/response.json"),
			args:     []string{"--id", "7b026645-e3dc-4d0f-91c5-6f4030996f6c"},
			output:   "Uuid:          123e4567-e89b-12d3-a456-426614174000  \nName:          MyToken                               \nCreated:       \"2024-07-10T12:34:56Z\"                \nExpires At:    \"2025-07-10T12:34:56Z\"                \nDescription:   \"Token for accessing API\"             \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)
			f, out, _ := testutils.NewFactory(mock)
			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				if !strings.EqualFold(tt.err, err.Error()) {
					t.Errorf("Error expected: %s got: %s", tt.err, err.Error())
				}
			} else {
				assert.Equal(t, tt.output, out.String())
			}
		})
	}
}
