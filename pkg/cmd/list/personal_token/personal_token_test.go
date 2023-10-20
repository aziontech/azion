package personaltoken

import (
	"errors"
	"testing"

	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name string
		args []string
		mock func() *httpmock.Registry
		err  error
	}{
		{
			name: "successful listing",
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("GET", "iam/personal_tokens"),
					httpmock.JSONFromFile("./fixtures/items.json"),
				)
				return &mock
			},
			err: nil,
		},
		{
			name: "no items",
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("GET", "iam/personal_tokens"),
					httpmock.JSONFromFile("./fixtures/no_items.json"),
				)
				return &mock
			},
			err: nil,
		},
		{
			name: "json invalid",
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("GET", "iam/personal_tokens"),
					httpmock.JSONFromString("{'name': 'some name',}"),
				)
				return &mock
			},
			err: errors.New("Failed to list your personal tokens: invalid character '\\'' looking for beginning of object key string. Check your settings and try again. If the error persists, contact Azion support."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, _ := testutils.NewFactory(tt.mock())

			cmd := NewCmd(f)
			cmd.SetArgs(tt.args)
			_, err := cmd.ExecuteC()

			if err != nil && !(err.Error() == tt.err.Error()) {
				t.Errorf("Executec() err = %v, \nexpected %v", err, tt.args)
			}
		})
	}
}
