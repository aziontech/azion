package http

import (
	"errors"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/pkg/v3api/graphql/http"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func mockHTTPEvents(f *cmdutil.Factory, logTime time.Time, limit string) (http.HTTPEventsResponse, error) {
	return http.HTTPEventsResponse{
		HTTPEvents: []http.HTTPEvent{
			{
				Host:              "nf6sxm2b1k.map.azionedge.net",
				GeolocRegion:      "Rio Grande do Sul",
				HTTPUserAgent:     "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
				RequestURI:        "/css/images/vector-bg.svg",
				Status:            200,
				Ts:                logTime,
				UpstreamBytesSent: 1703,
				RequestTime:       "1.543",
				RequestMethod:     "GET",
			},
		},
	}, nil
}

func mockHTTPEventsError(f *cmdutil.Factory, logTime time.Time, limit string) (http.HTTPEventsResponse, error) {
	return http.HTTPEventsResponse{}, errors.New("error")
}

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}

	f, _, _ := testutils.NewFactory(mock)

	cmd := NewLogsCmd(f)
	cobracmd := NewCobraCmd(cmd, f)
	assert.NotNil(t, cmd)
	assert.IsType(t, &cobra.Command{}, cobracmd)
}

func TestPrintLogs(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	tests := []struct {
		name           string
		mockHTTP       func(*cmdutil.Factory, time.Time, string) (http.HTTPEventsResponse, error)
		expectedOutput string
		expectError    bool
		tail           bool
		pretty         bool
	}{
		{
			name:           "successful log retrieval",
			mockHTTP:       mockHTTPEvents,
			expectedOutput: "Timestamp: ",
			expectError:    false,
			tail:           false,
			pretty:         true,
		},
		{
			name:           "error in log retrieval",
			mockHTTP:       mockHTTPEventsError,
			expectedOutput: "",
			expectError:    true,
			tail:           false,
			pretty:         true,
		},
		{
			name:           "pretty print logs",
			mockHTTP:       mockHTTPEvents,
			expectedOutput: "Timestamp: ",
			expectError:    false,
			tail:           false,
			pretty:         true,
		},
		{
			name:           "tail logs",
			mockHTTP:       mockHTTPEvents,
			expectedOutput: "Timestamp: ",
			expectError:    false,
			tail:           true,
			pretty:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var arguments []string

			mock := &httpmock.Registry{}

			f, stdout, _ := testutils.NewFactory(mock)

			cmd := NewLogsCmd(f)
			cmd.GetEvents = tt.mockHTTP
			cmd.Tail = tt.tail
			cobracmd := NewCobraCmd(cmd, f)
			if tt.pretty {
				arguments = append(arguments, "--pretty")
			}

			cobracmd.SetArgs(arguments)

			err := cobracmd.Execute()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				output := stdout.String()
				assert.Contains(t, output, tt.expectedOutput)
			}
		})
	}
}
