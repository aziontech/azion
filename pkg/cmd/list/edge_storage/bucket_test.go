package edge_storage

import (
	"net/http"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

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
			request:  httpmock.REST(http.MethodGet, "v4/storage/buckets"),
			response: httpmock.JSONFromFile("fixtures/response.json"),
			args:     []string{""},
			output:   "NAME                EDGE ACCESS  \narthur-morgan02     read_only    \narthur-morgan03     read_only    \narthur-morgan05     read_only    \narthur-morgan06     read_only    \nblue-bilbo          read_write   \ncourageous-thunder  read_write   \n",
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
				if !strings.EqualFold(tt.err, err.Error()) {
					t.Errorf("Error expected: %s got: %s", tt.err, err.Error())
				}
			} else {
				assert.Equal(t, tt.output, out.String())
			}
		})
	}
}
