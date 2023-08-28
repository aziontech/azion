package delete

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/personal-token"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name   string
		args   []string
		mock   func() *httpmock.Registry
		output string
	}{
		{
			name: "delete personal token by id",
			args: []string{"--id", "5c9c1854-45dd-11ee-be56-0242ac120002"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("DELETE", "iam/personal_tokens/5c9c1854-45dd-11ee-be56-0242ac120002"),
					httpmock.StatusStringResponse(204, ""),
				)
				return &mock
			},
			output: fmt.Sprintf(msg.DeleteOutputSuccess, "5c9c1854-45dd-11ee-be56-0242ac120002"),
		},
		{
			name: "delete personal tokens that is not found",
			args: []string{"--id", "5c9c1854-45dd-11ee-be56-0242ac120002"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("DELETE", "iam/personal_tokens/5c9c1854-45dd-11ee-be56-0242ac120002"),
					httpmock.StatusStringResponse(404, "Not Found"),
				)
				return &mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, out, _ := testutils.NewFactory(tt.mock())

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			if err != nil {
				assert.Error(t, err)
			}

			assert.Equal(t, tt.output, out.String())
		})
	}
}
