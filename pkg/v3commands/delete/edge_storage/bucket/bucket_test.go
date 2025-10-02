package bucket

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
)

func mockBucket(msg string) (string, error) {
	return "arthur-morgan", nil
}

func TestNewBucket(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name              string
		requests          []httpmock.Matcher
		responses         []httpmock.Responder
		args              []string
		output            string
		wantErr           bool
		Err               string
		mockInputs        func(string) (string, error)
		mockDeleteAll     func(*api.Client, context.Context, string, string) error
		mockConfirmDelete func(bool, string, bool) bool
	}{
		{
			name: "delete bucket command bucket of the edge-storage",
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodDelete, "v4/storage/buckets/arthur-morgan"),
			},
			responses: []httpmock.Responder{
				func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNoContent,
					}, nil
				},
			},
			args:   []string{"--name", "arthur-morgan"},
			output: fmt.Sprintf(msg.OUTPUT_DELETE_BUCKET, "arthur-morgan"),
		},
		{
			name: "delete bucket ask for name input",
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodDelete, "v4/storage/buckets/arthur-morgan"),
			},
			responses: []httpmock.Responder{
				func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusNoContent,
					}, nil
				},
			},
			output:     fmt.Sprintf(msg.OUTPUT_DELETE_BUCKET, "arthur-morgan"),
			mockInputs: mockBucket,
		},
		{
			name: "failed delete bucket internal error status 500",
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodPost, "v4/storage/buckets/arthur-morgan"),
			},
			responses: []httpmock.Responder{
				func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       io.NopCloser(bytes.NewBufferString("{}")),
						Header:     http.Header{"Content-Type": []string{"application/json"}},
					}, nil
				},
			},
			args: []string{"--name", "arthur-morgan"},
			Err:  fmt.Sprintf(msg.ERROR_DELETE_BUCKET, "The server could not process the request because an internal and unexpected problem occurred. Wait a few seconds and try again. For more information run the command again using the '--debug' flag. If the problem persists, contact Azion’s support"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			for i, req := range tt.requests {
				mock.Register(req, tt.responses[i])
			}

			f, out, _ := testutils.NewFactory(mock)

			deleteCmd := NewDeleteBucketCmd(f)
			deleteCmd.AskInput = tt.mockInputs
			cobraCmd := NewBucketCmd(deleteCmd, f)

			if len(tt.args) > 0 {
				cobraCmd.SetArgs(tt.args)
			}
			if err := cobraCmd.Execute(); err != nil {
				if !strings.EqualFold(tt.Err, err.Error()) {
					t.Errorf("Error expected: %s got: %s", tt.Err, err.Error())
				}
			} else {
				assert.Equal(t, tt.output, out.String())
			}
		})
	}
}
