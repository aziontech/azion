package root

import (
	"net/http"
	"os"

	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckAuthorizeMetricsCollection(t *testing.T) {
	tests := []struct {
		name              string
		authorizeMetrics  int
		globalFlagAll     bool
		mockConfirm       bool
		expectedErr       error
		expectedSettings  token.Settings
		mockWriteSettings func(settings token.Settings) error
	}{
		{
			name:             "Metrics Collection Authorized",
			authorizeMetrics: 1,
			globalFlagAll:    false,
			mockConfirm:      true,
			expectedErr:      nil,
			expectedSettings: token.Settings{
				AuthorizeMetricsCollection: 1,
			},
			mockWriteSettings: func(settings token.Settings) error {
				assert.Equal(t, token.Settings{
					AuthorizeMetricsCollection: 1,
				}, settings)
				return nil
			},
		},
		{
			name:             "Metrics Collection Not Authorized",
			authorizeMetrics: 0,
			globalFlagAll:    false,
			mockConfirm:      false,
			expectedErr:      nil,
			expectedSettings: token.Settings{
				AuthorizeMetricsCollection: 2,
			},
			mockWriteSettings: func(settings token.Settings) error {
				assert.Equal(t, token.Settings{
					AuthorizeMetricsCollection: 2,
				}, settings)
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			settings := token.Settings{AuthorizeMetricsCollection: tt.authorizeMetrics}
			confirmFn = func(globalFlagAll bool, msg string, defaultValue bool) bool {
				return tt.mockConfirm
			}
			err := checkAuthorizeMetricsCollection(cmd, tt.globalFlagAll, &settings)
			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedSettings.AuthorizeMetricsCollection, settings.AuthorizeMetricsCollection)
		})
	}
}

func TestVerifyUserInfo(t *testing.T) {
	tests := []struct {
		name             string
		settings         token.Settings
		expectedResponse bool
	}{
		{
			name: "Complete User Info",
			settings: token.Settings{
				ClientId: "clientID",
				Email:    "email@example.com",
			},
			expectedResponse: true,
		},
		{
			name: "Incomplete User Info",
			settings: token.Settings{
				ClientId: "",
				Email:    "email@example.com",
			},
			expectedResponse: false,
		},
		{
			name: "Empty User Info",
			settings: token.Settings{
				ClientId: "",
				Email:    "",
			},
			expectedResponse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := verifyUserInfo(&tt.settings)
			assert.Equal(t, tt.expectedResponse, result)
		})
	}
}

func TestCheckTokenSent(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type args struct {
		fact     *factoryRoot
		settings *token.Settings
		tokenStr token.Token
	}

	tests := []struct {
		name      string
		request   httpmock.Matcher
		response  httpmock.Responder
		args      args
		expectErr bool
	}{
		{
			name:     "invalid token",
			request:  httpmock.REST("GET", "token"),
			response: httpmock.StatusStringResponse(http.StatusUnauthorized, "{}"),
			args: args{
				fact: &factoryRoot{
					flags: flags{
						tokenFlag: "thisIsNotTheValidToken",
					},
				},
				settings: &token.Settings{},
				tokenStr: token.Token{
					Endpoint: "http://api.azion.net/token",
				},
			},
			expectErr: true,
		},
		{
			name:     "valid token",
			request:  httpmock.REST("GET", "user/me"),
			response: httpmock.StatusStringResponse(http.StatusOK, "{}"),
			args: args{
				fact: &factoryRoot{
					flags: flags{
						tokenFlag: "azion4d277d8dd2ef7597894615d97f17e358959",
					},
				},
				settings: &token.Settings{},
				tokenStr: token.Token{
					Endpoint: "http://api.azion.net/token",
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)

			f, _, _ := testutils.NewFactory(mock)
			tt.args.fact.Factory = f

			token, _ := token.New(&token.Config{
				Client: &http.Client{Transport: mock},
				Out:    os.Stdout,
			})

			err := checkTokenSent(tt.args.fact, tt.args.settings, token)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
