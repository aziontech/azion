package create

import (
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/personal-token"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
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
			name: "Create personal token with success",
			args: []string{"--name", "sakura", "--expiration", "9m", "--description", "example"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "iam/personal_tokens"),
					httpmock.JSONFromFile("./fixtures/successfully_created.json"),
				)
				return &mock
			},
			output: fmt.Sprintf("🚀 "+msg.CreateOutputSuccess, "super-secret-token\n"),
			err:    nil,
		},
		{
			name: "Create new personal token with json file using --in flag",
			args: []string{"--in", "./fixtures/complete_structure.json"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				mock.Register(
					httpmock.REST("POST", "iam/personal_tokens"),
					httpmock.JSONFromFile("./fixtures/successfully_created.json"),
				)
				return &mock
			},
			output: fmt.Sprintf("🚀 "+msg.CreateOutputSuccess, "super-secret-token\n"),
			err:    nil,
		},
		{
			name: "Failed to create missing required flag",
			args: []string{"--name", "sakura", "--description", "example"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				return &mock
			},
			err: msg.ErrorMandatoryCreateFlags,
		},
		{
			name: "Failed to create argument empty",
			args: []string{"--name", "", "--expiration", "9m", "--description", "example"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				return &mock
			},
			err: msg.ErrorMandatoryCreateFlags,
		},
		{
			name: "Missing expiration field in json structure passed by --in flag",
			args: []string{"--in", "./fixtures/structure_without_expiration_field.json"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				return &mock
			},
			err: msg.ErrorMandatoryCreateFlags,
		},
		{
			name: "Failure to create expiration date with invalid format",
			args: []string{"--name", "luffy", "--expiration", "2323", "--description", "example"},
			mock: func() *httpmock.Registry {
				mock := httpmock.Registry{}
				return &mock
			},
			err: fmt.Errorf("invalid date format, what do we expect: \"1d\", \"2w\", \"2m\", \"1y\", \"18/08/2023\", \"2023-02-12\""),
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
