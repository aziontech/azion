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

var successResponse string = `
{
  "data": {
    "id": 1337,
    "name": "Updated Stream",
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "created": "2019-08-24T14:15:22Z",
    "active": false,
    "product_version": "1.0",
    "inputs": [],
    "transform": [],
    "outputs": []
  }
}
`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("update with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/stream/streams/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--stream-id", "1337", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--stream-id", "1337", "--file", "./fixtures/does-not-exist.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/stream/streams/9999"),
			httpmock.StatusStringResponse(http.StatusNotFound, `{"details": "Stream not found"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--stream-id", "9999", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
