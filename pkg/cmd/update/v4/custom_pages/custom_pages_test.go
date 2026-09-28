package custompages

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/custom_pages"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponse string = `
{
  "data": {
    "id": 1337,
    "name": "Updated Custom Page",
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "created_at": "2019-08-24T14:15:22Z",
    "active": false,
    "product_version": "1.0",
    "pages": [],
    "is_versioned": true,
    "version": 1,
    "version_state": "string",
    "version_id": "string"
  }
}
`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("update with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/custom_pages/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--custom-page-id", "1337", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("file not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--custom-page-id", "1337", "--file", "./fixtures/does-not-exist.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/custom_pages/9999"),
			httpmock.StatusStringResponse(http.StatusNotFound, `{"details": "Custom Page not found"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--custom-page-id", "9999", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
