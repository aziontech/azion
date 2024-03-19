package edge_storage

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
)

func TestNewBucket(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name     string
		request  httpmock.Matcher
		response httpmock.Responder
		args     []string
		output   string
		wantErr  bool
		Err      string
	}{
		{
			name:     "delete bucket command bucket of the edge-storage",
			request:  httpmock.REST(http.MethodDelete, "v4/storage/buckets/arthur-morgan"),
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusNoContent,
				}, nil
			},
			args:     []string{"--name", "arthur-morgan"},
			output:   fmt.Sprintf(msg.OUTPUT_DELETE_BUCKET, "arthur-morgan"),
		},
		{
			name:    "failed delete bucket internal error status 500",
			request: httpmock.REST(http.MethodPost, "v4/storage/buckets/arthur-morgan"),
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			args:     []string{"--name", "arthur-morgan"},
			Err:  fmt.Sprintf(msg.ERROR_DELETE_BUCKET, "The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)
			f, out, _ := testutils.NewFactory(mock)
			cmd := NewBucket(f)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				if !strings.EqualFold(tt.Err, err.Error()) {
					t.Errorf("Error expected: %s got: %s", tt.Err, err.Error())
				}
			} else {
				assert.Equal(t, tt.output, out.String())
			}
		})
	}
}
