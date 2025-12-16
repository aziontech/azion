package storage

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/storage"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
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
			name:     "update bucket command bucket of the storage",
			request:  httpmock.REST(http.MethodPatch, "workspace/storage/buckets/John-Marston"),
			response: httpmock.StatusStringResponse(http.StatusNoContent, ""),
			args:     []string{"--name", "John-Marston", "--edge-access", "read_only"},
			output:   msg.OUTPUT_UPDATE_BUCKET,
		},
		{
			name:     "create new bucket command bucket of the storage using flag --file",
			request:  httpmock.REST(http.MethodPatch, "workspace/storage/buckets/John-Marston"),
			response: httpmock.StatusStringResponse(http.StatusNoContent, ""),
			args:     []string{"--file", "fixtures/create.json"},
			output:   msg.OUTPUT_UPDATE_BUCKET,
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
			request: httpmock.REST(http.MethodPatch, "workspace/storage/buckets/John-Marston"),
			response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString("")),
					Header:     http.Header{"Content-Type": []string{"application/json"}},
				}, nil
			},
			args: []string{"--file", "fixtures/create.json"},
			Err:  fmt.Sprintf(msg.ERROR_UPDATE_BUCKET, "The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support"),
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
