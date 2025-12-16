package object

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/storage"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
)

func TestNewObject(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name     string
		request  httpmock.Matcher
		response httpmock.Responder
		args     []string
		output   string
		err      string
	}{
		{
			name: "list success",
			request: func(req *http.Request) bool {
				return strings.EqualFold(req.Method, http.MethodGet) && strings.Contains(req.URL.Path, "/workspace/storage/buckets/my-bucket/objects")
			},
			response: httpmock.JSONFromFile("fixtures/objects_response.json"),
			args:     []string{"--bucket-name", "my-bucket"},
			output:   "KEY      LAST MODIFIED                  \nobject1  2024-08-05 12:00:00 +0000 UTC  \nobject2  2024-08-05 13:00:00 +0000 UTC  \n",
		},
		{
			name: "list 2 items successfully",
			request: func(req *http.Request) bool {
				return strings.EqualFold(req.Method, http.MethodGet) && strings.Contains(req.URL.Path, "/workspace/storage/buckets/my-bucket/objects")
			},
			response: httpmock.JSONFromFile("fixtures/objects_response_2_items.json"),
			args:     []string{"--bucket-name", "my-bucket", "--page-size", "2"},
			output:   "KEY      LAST MODIFIED                  \nobject1  2024-08-05 12:00:00 +0000 UTC  \nobject2  2024-08-05 13:00:00 +0000 UTC  \n",
		},
		{
			name: "failed internal error status 500",
			request: func(req *http.Request) bool {
				return strings.EqualFold(req.Method, http.MethodGet) && strings.Contains(req.URL.Path, "/workspace/storage/buckets/my-bucket/objects")
			},
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			args: []string{"--bucket-name", "my-bucket"},
			err:  fmt.Sprintf(msg.ERROR_LIST_BUCKET, "The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)
			f, out, _ := testutils.NewFactory(mock)
			cmd := NewObject(f)
			cmd.SetArgs(tt.args)
			if err := cmd.Execute(); err != nil {
				if !strings.EqualFold(tt.err, err.Error()) {
					t.Errorf("Error expected: %s got: %s", tt.err, err.Error())
				}
			} else {
				assert.Equal(t, tt.output, out.String())
			}
		})
	}
}
