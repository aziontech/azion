package bucket

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

func TestNewBucket(t *testing.T) {
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
			name:     "list success",
			request:  httpmock.REST(http.MethodGet, "workspace/storage/buckets"),
			response: httpmock.JSONFromFile("fixtures/response.json"),
			args:     []string{""},
			output:   "NAME                EDGE ACCESS  \narthur-morgan02     read_only    \narthur-morgan03     read_only    \narthur-morgan05     read_only    \narthur-morgan06     read_only    \nblue-bilbo          read_write   \ncourageous-thunder  read_write   \n",
		},
		{
			name:     "list 2 items successfully",
			request:  httpmock.REST(http.MethodGet, "workspace/storage/buckets"),
			response: httpmock.JSONFromFile("fixtures/response_2_items.json"),
			args:     []string{"--page", "1", "--page-size", "2"},
			output:   "NAME                EDGE ACCESS  \narthur-morgan02     read_only    \narthur-morgan03     read_only    \n",
		},
		{
			name:    "failed internal error status 500",
			request: httpmock.REST(http.MethodGet, "workspace/storage/buckets"),
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			err: fmt.Sprintf(msg.ERROR_LIST_BUCKET, "The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support"),
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)
			f, out, _ := testutils.NewFactory(mock)
			cmd := NewBucket(f)
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
