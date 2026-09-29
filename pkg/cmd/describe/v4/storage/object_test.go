package storage

import (
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDescribeObject(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name      string
		request   httpmock.Matcher
		response  httpmock.Responder
		args      []string
		expectErr bool
	}{
		{
			name:      "object not found",
			request:   httpmock.REST("GET", "storage/buckets/unknown/objects/unknown-object"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--bucket-name", "unknown", "--object-key", "unknown-object"},
			expectErr: true,
		},
		{
			name:      "missing bucket name",
			request:   httpmock.REST("GET", "storage/buckets/missing/objects/test-object"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--object-key", "test-object"},
			expectErr: true,
		},
		{
			name:      "missing object key",
			request:   httpmock.REST("GET", "storage/buckets/test-bucket/objects/"),
			response:  httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
			args:      []string{"--bucket-name", "test-bucket"},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)

			factory, _, _ := testutils.NewFactory(mock)
			cmd := NewObject(factory)
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
