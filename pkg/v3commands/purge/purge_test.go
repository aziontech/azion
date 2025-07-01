package purge

import (
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/purge"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("purge urls", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "purge/url"),
			httpmock.StatusStringResponse(201, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		err := purgeUrls([]string{"www.example.com", "www.httpin.com"}, f)

		require.NoError(t, err)
	})

	t.Run("purge wildcard", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "purge/wildcard"),
			httpmock.StatusStringResponse(201, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		err := purgeWildcard([]string{"www.example.com/*"}, f)

		require.NoError(t, err)
	})

	t.Run("purge wildcard - more than one item", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "purge/wildcard"),
			httpmock.StatusStringResponse(201, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		err := purgeWildcard([]string{"www.example.com/*", "www.pudim.com/*"}, f)

		require.ErrorContains(t, err, msg.ErrorTooManyUrls.Error())
	})

	t.Run("purge cache-keys", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "purge/cachekey"),
			httpmock.StatusStringResponse(201, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		err := purgeCacheKeys([]string{"www.example.com/@@cookie_name=cookie_value"}, f)

		require.NoError(t, err)
	})

}

func TestPurge(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name         string
		request      httpmock.Matcher
		response     httpmock.Responder
		args         []string
		expectErr    bool
		mockInput    func() ([]string, error)
		mockGetPurge func() (string, error)
	}{
		{
			name:      "purge urls",
			request:   httpmock.REST("POST", "purge/url"),
			response:  httpmock.StatusStringResponse(201, ""),
			args:      []string{"--urls", "http://www.example.com/,http://www.pudim.com/"},
			expectErr: false,
		},
		{
			name:      "purge urls - ask input",
			request:   httpmock.REST("POST", "purge/url"),
			response:  httpmock.StatusStringResponse(201, ""),
			expectErr: false,
			mockInput: func() ([]string, error) {
				return []string{"www.example.com/", "www.pudim.com/"}, nil
			},
			mockGetPurge: func() (string, error) {
				return "URLs", nil
			},
		},
		{
			name:      "purge wildcard - ask input",
			request:   httpmock.REST("POST", "purge/wildcard"),
			response:  httpmock.StatusStringResponse(201, ""),
			expectErr: false,
			mockInput: func() ([]string, error) {
				return []string{"www.example.com/*"}, nil
			},
			mockGetPurge: func() (string, error) {
				return "Wildcard", nil
			},
		},
		{
			name:      "purge cachekey - ask input",
			request:   httpmock.REST("POST", "purge/cachekey"),
			response:  httpmock.StatusStringResponse(201, ""),
			expectErr: false,
			mockInput: func() ([]string, error) {
				return []string{"www.domain.com/@@cookie_name=cookie_value"}, nil
			},
			mockGetPurge: func() (string, error) {
				return "Cache-Key", nil
			},
		},
		{
			name:      "purge wildcard",
			request:   httpmock.REST("POST", "purge/wildcard"),
			response:  httpmock.StatusStringResponse(201, ""),
			args:      []string{"--wildcard", "www.example.com/*"},
			expectErr: false,
		},
		{
			name:      "purge cache keys",
			request:   httpmock.REST("POST", "purge/cachekey"),
			response:  httpmock.JSONFromFile("./fixtures/response.json"),
			args:      []string{"--cache-key", "www.domain.com/@@cookie_name=cookie_value,www.domain.com/test.js"},
			expectErr: false,
		},
		{
			name:      "invalid urls",
			request:   httpmock.REST("POST", "purge/url"),
			response:  httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid URL"),
			args:      []string{"--urls", "invalid-url"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)

			f, _, _ := testutils.NewFactory(mock)
			purgeCmd := NewPurgeCmd(f)
			purgeCmd.AskForInput = tt.mockInput
			purgeCmd.GetPurgeType = tt.mockGetPurge
			cmd := NewCobraCmd(purgeCmd, f)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetPurgeType(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantAnswer string
		wantErr    bool
	}{
		{
			name:       "select URLs",
			input:      "URLs\n",
			wantAnswer: "URLs",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getPurgeType()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAskForInput(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantAnswer []string
		wantErr    bool
	}{
		{
			name:       "single URL",
			input:      "www.example.com\n",
			wantAnswer: []string{"www.example.com"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := askForInput()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
