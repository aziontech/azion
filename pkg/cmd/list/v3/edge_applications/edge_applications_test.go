package edge_applications

import (
	"fmt"
	"testing"

	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/list/applications"
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
			args: []string{},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("GET", "edge_applications"),
					httpmock.JSONFromFile("./fixtures/response.json"),
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
					httpmock.REST("GET", "edge_applications"),
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
					httpmock.REST("GET", "edge_applications"),
					httpmock.JSONFromString("{'name': 'some name',}"),
				)
				return &mock
			},
			// Built from the constant the command uses, so a reworded message
			// cannot leave this expectation silently stale.
			err: fmt.Errorf(msg.ErrorGetAll.Error(), "invalid character '\\'' looking for beginning of object key string"),
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
