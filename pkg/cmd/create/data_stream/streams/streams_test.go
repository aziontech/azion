package streams

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/streams"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var errorResponse string = `
{
  "errors": [
    {
      "status": "500",
      "code": "500",
      "title": "string",
      "detail": "string"
    }
  ]
}
`

var successResponse string = `
{
  "data": {
    "id": 1337,
    "name": "My Stream",
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "created": "2019-08-24T14:15:22Z",
    "active": true,
    "product_version": "1.0",
    "inputs": [],
    "transform": [],
    "outputs": []
  }
}
`

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("create with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/stream/streams"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/create.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/does-not-exist.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("api error", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/stream/streams"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, errorResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/create.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
