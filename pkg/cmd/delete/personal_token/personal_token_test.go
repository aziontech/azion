package personaltoken

import (
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/delete/personal_token"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name   string
		args   []string
		mock   func() *httpmock.Registry
		output string
		err    error
	}{
		{
			name: "Delete personal token by id",
			args: []string{"--id", "5c9c1854-45dd-11ee-be56-0242ac120002"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("DELETE", "iam/personal_tokens/5c9c1854-45dd-11ee-be56-0242ac120002"),
					httpmock.StatusStringResponse(204, ""),
				)
				return &mock
			},
			output: fmt.Sprintf(personaltoken.OutputSuccess, "5c9c1854-45dd-11ee-be56-0242ac120002"),
			err:    nil,
		},
		{
			name: "Delete personal tokens that is not found",
			args: []string{"--id", "5c9c1854-45dd-11ee-be56-0242ac120002"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("DELETE", "iam/personal_tokens/5c9c1854-45dd-11ee-be56-0242ac120002"),
					httpmock.StatusStringResponse(404, "Not Found"),
				)
				return &mock
			},
			err: fmt.Errorf(personaltoken.ErrorFailToDelete.Error(), utils.ErrorNotFound404),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, out, _ := testutils.NewFactory(tt.mock())

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			if err != nil && !(err.Error() == tt.err.Error()) {
				t.Errorf("Executec() err = %v, \nexpected %v", err, tt.args)
			}

			assert.Equal(t, tt.output, out.String())
		})
	}
}
