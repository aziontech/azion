package login

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

func Test_login_browserLogin(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		askOne      func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
		run         func(input string) error
		server      Server
		token       token.TokenInterface
		marshalToml func(v interface{}) ([]byte, error)
		askInput    func(msg string) (string, error)
	}
	type args struct {
		srv Server
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error Open browser",
			fields: fields{
				run: func(input string) error {
					return errors.New("error open browser")
				},
			},
			args: args{
				srv: &ServerMock{
					Cancel:               true,
					ErrorListenAndServer: nil,
					ErrorShutdown:        nil,
				},
			},
			wantErr: true,
		},
		{
			name: "error shotdown",
			fields: fields{
				run: func(input string) error {
					return nil
				},
			},
			args: args{
				srv: &ServerMock{
					Cancel:               true,
					ErrorListenAndServer: nil,
					ErrorShutdown:        errors.Errorf("iiiisssssh"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _, _ := testutils.NewFactory(&httpmock.Registry{})
			l := &login{
				factory:     f,
				askOne:      tt.fields.askOne,
				run:         tt.fields.run,
				server:      tt.fields.server,
				token:       tt.fields.token,
				marshalToml: tt.fields.marshalToml,
				askInput:    tt.fields.askInput,
			}
			enableHandlerRouter = false
			if err := l.browserLogin(); (err != nil) != tt.wantErr {
				t.Errorf("login.browserLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
