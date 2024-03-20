package edge_storage

import (
	"bytes"
	"fmt"
	"github.com/aziontech/azion-cli/utils"
	"io"
	"net/http"
	"strings"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewObjects(t *testing.T) {
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
			name:     "create new object command object of the edge-storage",
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--bucket-name", "nomedobucket", "--object-key", "nomedoobject", "--source", "fixtures/index.html"},
			output:   msg.OUTPUT_CREATE_OBJECT,
		},
		{
			name:     "create new object command object of the edge-storage using flag --file",
			request:  httpmock.REST(http.MethodPost, "v4/storage/buckets/nomedobucket/objects/nomedoobject"),
			response: httpmock.JSONFromFile("fixtures/response_object.json"),
			args:     []string{"--file", "fixtures/create_object.json"},
			output:   msg.OUTPUT_CREATE_OBJECT,
		},
		{
			name:     "input file json err --file",
			request:  httpmock.REST(http.MethodPost, "/"),
			response: httpmock.JSONFromFile("/"),
			args:     []string{"--file", "fixtures/create_error.json"},
			Err:      utils.ErrorUnmarshalReader.Error(),
		},
		{
			name:    "failed internal error status 500",
			request: httpmock.REST(http.MethodPost, "v4/storage/buckets"),
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, utils.ErrorInternalServerError
			},
			args: []string{"--file", "fixtures/create_object.json"},
			Err:  fmt.Sprintf(msg.ERROR_CREATE_OBJECT, "The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			mock.Register(tt.request, tt.response)
			f, out, _ := testutils.NewFactory(mock)
			cmd := NewObjects(f)
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
