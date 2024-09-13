package login

import (
	"fmt"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/spf13/cobra"
	"go.uber.org/zap/zapcore"
)

func Test_cmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	type args struct {
		l *login
	}

	f, _, _ := testutils.NewFactory(&httpmock.Registry{})
	tests := []struct {
		name string
		args args
		cmd  func(l *login) *cobra.Command
	}{
		{
			name: "success flow",
			args: args{
				l: &login{
					factory: f,
					askOne: func(
						p survey.Prompt,
						response interface{},
						opts ...survey.AskOpt,
					) error {
						// Usamos type assertion para garantir que 'response' seja um ponteiro para string
						if respPtr, ok := response.(*string); ok {
							*respPtr = "browser"
						}
						return nil
					},
					run:    func(input string) error { return nil },
					server: &ServerMock{},
				},
			},
			cmd: cmd,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmd(tt.args.l)
			err := cmd.Execute()
			fmt.Println(err)
		})
	}
}
