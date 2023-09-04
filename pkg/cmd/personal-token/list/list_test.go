package list

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
			name: "listing with success",
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
			name: "not found",
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("GET", "iam/personal_tokens"),
					httpmock.StatusStringResponse(404, "Not Found"),
				)
				return &mock
			},
			err: errors.New("Failed to describe the personal token: The given web page URL or API's endpoint doesn't exist or isn't available. Check that the identifying information is correct. If the error persists, contact Azion's support. Check your settings and try again. If the error persists, contact Azion support."),
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
			err: errors.New("Failed to describe the personal token: invalid character '\\'' looking for beginning of object key string. Check your settings and try again. If the error persists, contact Azion support."),
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
