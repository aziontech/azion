package dnszone

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/dns_zone"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockZoneID(msg string) (string, error) {
	return "1234", nil
}

func mockInvalidZoneID(msg string) (string, error) {
	return "invalid", nil
}

func mockParseErrorZoneID(msg string) (string, error) {
	return "invalid", utils.ErrorParseResponse
}

func TestDelete(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name           string
		zoneID         string
		method         string
		endpoint       string
		statusCode     int
		responseBody   string
		expectedOutput string
		expectError    bool
		mockInputs     func(string) (string, error)
	}{
		{
			name:           "delete DNS zone by id",
			zoneID:         "1234",
			method:         "DELETE",
			endpoint:       "workspace/dns/zones/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DNSZoneDeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockZoneID,
		},
		{
			name:         "delete DNS zone - not found",
			zoneID:       "1234",
			method:       "DELETE",
			endpoint:     "workspace/dns/zones/1234",
			statusCode:   404,
			responseBody: "Not Found",
			expectError:  true,
			mockInputs:   mockZoneID,
		},
		{
			name:           "ask for zone id success",
			zoneID:         "",
			method:         "DELETE",
			endpoint:       "workspace/dns/zones/1234",
			statusCode:     204,
			responseBody:   "",
			expectedOutput: fmt.Sprintf(msg.DNSZoneDeleteOutputSuccess, 1234),
			expectError:    false,
			mockInputs:     mockZoneID,
		},
		{
			name:        "ask for zone id conversion failure",
			zoneID:      "",
			method:      "",
			endpoint:    "",
			statusCode:  0,
			expectError: true,
			mockInputs:  mockInvalidZoneID,
		},
		{
			name:        "error - parse answer",
			zoneID:      "",
			method:      "",
			endpoint:    "",
			statusCode:  0,
			expectError: true,
			mockInputs:  mockParseErrorZoneID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(
				httpmock.REST(tt.method, tt.endpoint),
				httpmock.StatusStringResponse(tt.statusCode, tt.responseBody),
			)

			f, stdout, _ := testutils.NewFactory(mock)

			deleteCmd := NewDeleteCmd(f)
			deleteCmd.AskInput = tt.mockInputs
			cobraCmd := NewCobraCmd(deleteCmd, f)

			if tt.zoneID != "" {
				cobraCmd.SetArgs([]string{"--zone-id", tt.zoneID})
			}

			_, err := cobraCmd.ExecuteC()
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedOutput, stdout.String())
			}
		})
	}
}
