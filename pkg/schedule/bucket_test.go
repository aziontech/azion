package schedule

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
)

func TestTriggerDeleteBucket(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	type args struct {
		name string
	}

	tests := []struct {
		name      string
		args      args
		requests  []httpmock.Matcher
		responses []httpmock.Responder
		wantErr   bool
	}{
		{
			name: "success",
			args: args{"arthur-morgan"},
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodDelete, "workspace/storage/buckets/arthur-morgan"),
			},
			responses: []httpmock.Responder{
				httpmock.StatusStringResponse(http.StatusNoContent, ""),
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{"arthur-morgan"},
			requests: []httpmock.Matcher{
				httpmock.REST(http.MethodDelete, "workspace/storage/buckets/arthur-morgan"),
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
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			for i, req := range tt.requests {
				mock.Register(req, tt.responses[i])
			}

			f, _, _ := testutils.NewFactory(mock)
			if err := TriggerDeleteBucket(f, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("TriggerDeleteBucket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
