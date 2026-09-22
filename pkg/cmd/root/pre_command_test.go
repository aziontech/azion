package root

import (
	"errors"
	"net/http"
	"os"

	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
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
			err := checkAuthorizeMetricsCollection(cmd, tt.globalFlagAll, &settings, "default")
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
			tt.args.fact.factory = f

			token := token.New(&token.Config{
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

// newCommandTree builds a command tree shaped like the CLI's, so the pre-command
// checks are exercised against the same nesting they walk through in production
func newCommandTree() *cobra.Command {
	root := &cobra.Command{Use: "azion"}
	root.PersistentFlags().String("token", "", "")

	create := &cobra.Command{Use: "create"}
	dataStream := &cobra.Command{Use: "data-stream"}
	dataStream.AddCommand(&cobra.Command{Use: "streams"})
	create.AddCommand(dataStream)

	root.AddCommand(
		&cobra.Command{Use: "build"},
		&cobra.Command{Use: "version"},
		&cobra.Command{Use: "__complete"},
		create,
	)

	return root
}

// command returns the command reached by the given path, e.g. "create data-stream streams"
func command(t *testing.T, path ...string) *cobra.Command {
	t.Helper()

	current := newCommandTree()
	for _, name := range path {
		found := false
		for _, sub := range current.Commands() {
			if sub.Name() == name {
				current, found = sub, true
				break
			}
		}
		require.True(t, found, "command %q not found in the test command tree", name)
	}

	return current
}

func TestTopLevelCommandName(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		expected string
	}{
		{
			name:     "root command itself",
			path:     nil,
			expected: "",
		},
		{
			name:     "top-level command",
			path:     []string{"build"},
			expected: "build",
		},
		{
			name:     "nested command returns its top-level parent",
			path:     []string{"create", "data-stream"},
			expected: "create",
		},
		{
			name:     "deeply nested command returns its top-level parent",
			path:     []string{"create", "data-stream", "streams"},
			expected: "create",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, topLevelCommandName(command(t, tt.path...)))
		})
	}
}

// unreachableAPITransport simulates the authentication API being unreachable:
// no network, DNS failure or timeout
type unreachableAPITransport struct{}

func (unreachableAPITransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("dial tcp: lookup api.azion.net: no such host")
}

// forbiddenTransport fails the test if the token check reaches the network at all
type forbiddenTransport struct{ t *testing.T }

func (f forbiddenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	f.t.Errorf("checkTokenNotExpired should not have requested %s", req.URL)
	return nil, errors.New("unexpected request")
}

func respondWith(status int) func(*testing.T) http.RoundTripper {
	return func(t *testing.T) http.RoundTripper {
		mock := &httpmock.Registry{}
		mock.Register(httpmock.REST("GET", "user/me"), httpmock.StatusStringResponse(status, "{}"))
		return mock
	}
}

func noRequestExpected(t *testing.T) http.RoundTripper {
	return forbiddenTransport{t: t}
}

func TestCheckTokenNotExpired(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	configuredToken := &token.Settings{Token: "azion4d277d8dd2ef7597894615d97f17e358959"}

	tests := []struct {
		name        string
		cmd         func(*testing.T) *cobra.Command
		settings    *token.Settings
		transport   func(*testing.T) http.RoundTripper
		expectedErr error
	}{
		{
			name: "token sent through the flag was already validated",
			cmd: func(t *testing.T) *cobra.Command {
				cmd := command(t, "build")
				cmd.Flags().String("token", "", "")
				require.NoError(t, cmd.Flags().Set("token", "azion4d277d8dd2ef7597894615d97f17e358959"))
				return cmd
			},
			settings:  configuredToken,
			transport: noRequestExpected,
		},
		{
			name:      "command that does not require a token",
			cmd:       func(t *testing.T) *cobra.Command { return command(t, "version") },
			settings:  configuredToken,
			transport: noRequestExpected,
		},
		{
			name:      "root command prints the help message",
			cmd:       func(t *testing.T) *cobra.Command { return command(t) },
			settings:  configuredToken,
			transport: noRequestExpected,
		},
		{
			name:      "shell completion command",
			cmd:       func(t *testing.T) *cobra.Command { return command(t, "__complete") },
			settings:  configuredToken,
			transport: noRequestExpected,
		},
		{
			name:      "no settings configured yet",
			cmd:       func(t *testing.T) *cobra.Command { return command(t, "build") },
			settings:  nil,
			transport: noRequestExpected,
		},
		{
			name:      "no token configured yet",
			cmd:       func(t *testing.T) *cobra.Command { return command(t, "build") },
			settings:  &token.Settings{},
			transport: noRequestExpected,
		},
		{
			name:      "valid token",
			cmd:       func(t *testing.T) *cobra.Command { return command(t, "create", "data-stream", "streams") },
			settings:  configuredToken,
			transport: respondWith(http.StatusOK),
		},
		{
			name:        "token rejected by the API",
			cmd:         func(t *testing.T) *cobra.Command { return command(t, "create", "data-stream", "streams") },
			settings:    configuredToken,
			transport:   respondWith(http.StatusUnauthorized),
			expectedErr: utils.ErrorToken401,
		},
		{
			// a network failure says nothing about the token, so the command must run
			// and report the failure itself instead of asking for a new token
			name:      "API unreachable does not report an expired token",
			cmd:       func(t *testing.T) *cobra.Command { return command(t, "build") },
			settings:  configuredToken,
			transport: func(t *testing.T) http.RoundTripper { return unreachableAPITransport{} },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &http.Client{Transport: tt.transport(t)}

			f, _, _ := testutils.NewFactory(&httpmock.Registry{})
			f.HttpClient = client

			fact := &factoryRoot{
				factory: f,
				globals: globals{globalSettings: tt.settings},
			}

			tk := token.New(&token.Config{Client: client, Out: f.IOStreams.Out})

			err := checkTokenNotExpired(tt.cmd(t), fact, tk)
			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				return
			}
			assert.NoError(t, err)
		})
	}
}

// Deleting the profile that held the only working credential used to leave the
// CLI unusable: the commands that create or remove a profile were themselves
// gated behind a valid token, so there was no way back in.
func TestProfileManagementIsNotGatedByTheToken(t *testing.T) {
	root := &cobra.Command{Use: "azion"}
	for _, verb := range []string{"create", "delete"} {
		parent := &cobra.Command{Use: verb}
		parent.AddCommand(&cobra.Command{Use: "profile"})
		parent.AddCommand(&cobra.Command{Use: "variables"})
		root.AddCommand(parent)
	}

	find := func(path ...string) *cobra.Command {
		cur := root
		for _, name := range path {
			for _, c := range cur.Commands() {
				if c.Name() == name {
					cur = c
					break
				}
			}
		}
		return cur
	}

	assert.Equal(t, "create profile", commandPathWithoutBinary(find("create", "profile")))
	assert.Equal(t, "delete profile", commandPathWithoutBinary(find("delete", "profile")))

	assert.True(t, pathsWithoutToken["create profile"], "creating a profile is how a user recovers")
	assert.True(t, pathsWithoutToken["delete profile"], "deleting a broken profile must not need a token")
	assert.False(t, pathsWithoutToken["create variables"], "real resources still require a valid credential")
}
