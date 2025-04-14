package login

import (
	"errors"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/pkg/token"
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
		name    string
		args    args
		cmd     func(l *login) *cobra.Command
		wantErr bool
	}{
		{
			name: "error select login mode",
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
						return errors.New("error select login mode")
					},
					run:    func(input string) error { return nil },
					server: &ServerMock{Cancel: true},
					token:  &token.TokenMock{},
					marshalToml: func(v interface{}) ([]byte, error) {
						return []byte(""), nil
					},
				},
			},
			cmd:     cmd,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.cmd(tt.args.l)
			if err := cmd.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_New(t *testing.T) {
	mock := &httpmock.Registry{}
	f, _, _ := testutils.NewFactory(mock)
	New(f)
}
