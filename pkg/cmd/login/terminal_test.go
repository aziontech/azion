package login

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/pkg/token"
	"go.uber.org/zap/zapcore"
)

func Test_login_terminalLogin(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type fields struct {
		factory     *cmdutil.Factory
		askOne      func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
		run         func(input string) error
		server      Server
		token       token.TokenInterface
		marshalToml func(v interface{}) ([]byte, error)
		askInput    func(msg string) (string, error)
		askPassword func(msg string) (string, error)
	}

	type register struct {
		matcher  httpmock.Matcher
		reponder httpmock.Responder
	}

	tests := []struct {
		name     string
		fields   fields
		flags    []string
		register register
		wantErr  bool
	}{
		{
			name: "success flow",
			fields: fields{
				run: func(input string) error {
					return nil
				},
				askOne: func(
					p survey.Prompt,
					response interface{},
					opts ...survey.AskOpt,
				) error {
					// Usamos type assertion para garantir que 'response' seja um ponteiro para string
					if respPtr, ok := response.(*string); ok {
						*respPtr = "terminal"
					}
					return nil
				},
				askInput: func(msg string) (string, error) {
					return "admin", nil
				},
				askPassword: func(msg string) (string, error) {
					return "admin", nil
				},
				token: &token.TokenMock{},
			},
			flags: []string{"--username", "max", "--password", "1235"},
			register: register{
				httpmock.REST("POST", "iam/personal_tokens"),
				httpmock.JSONFromFile("./fixtures/response.json"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.register.matcher, tt.register.reponder)
			f, _, _ := testutils.NewFactory(mock)
			l := &login{
				factory:     f,
				askOne:      tt.fields.askOne,
				run:         tt.fields.run,
				server:      tt.fields.server,
				token:       tt.fields.token,
				marshalToml: tt.fields.marshalToml,
				askInput:    tt.fields.askInput,
				askPassword: tt.fields.askPassword,
			}
			cmd := cmd(l)
			cmd.SetArgs(tt.flags)
			if err := l.terminalLogin(cmd); (err != nil) != tt.wantErr {
				t.Errorf("login.terminalLogin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
