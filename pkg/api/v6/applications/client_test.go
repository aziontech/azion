package applications

import (
	"context"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

// A 2xx the SDK cannot decode must surface as an error. Overwriting the SDK
// error with the result of rewinding the body used to turn it into (nil, nil),
// and the caller then dereferenced the nil response.
func TestUndecodableSuccessIsAnError(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	mock := &httpmock.Registry{}
	mock.Register(
		httpmock.REST("GET", "workspace/applications/1"),
		httpmock.StatusStringResponse(http.StatusOK, "not json"),
	)
	client := NewClient(&http.Client{Transport: mock}, "https://api.test/v6", "token")

	resp, err := client.Get(context.Background(), 1)

	require.Error(t, err)
	require.Nil(t, resp)
}
