package variables

import (
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/variables"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var successResponse string = `
{
	"uuid": "32e8ffca-4021-49a4-971f-330935566af4",
	"key": "Content-Type",
	"value": "json",
	"secret": false,
	"last_editor": "ei@tcha.com",
	"created_at": "2023-06-13T13:17:13.145625Z",
	"updated_at": "2023-06-13T13:17:13.145666Z"
}
`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("update domain", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PUT", "variables/32e8ffca-4021-49a4-971f-330935566af4"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--variable-id", "32e8ffca-4021-49a4-971f-330935566af4", "--key", "Content-Type", "--value", "int", "--secret", "false"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "🚀 Updated Variable with ID 32e8ffca-4021-49a4-971f-330935566af4\n\n", stdout.String())
	})

	t.Run("missing fields", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PUT", "variables/32e8ffca-4021-49a4-971f-330935566af4"),
			httpmock.JSONFromString(successResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--variable-id", "32e8ffca-4021-49a4-971f-330935566af4", "--key", "Content-Type"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorMissingFieldUpdateVariables)

	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PUT", "variables/32e8ffca-4021-49a4-971f-330935566af4"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--variable-id", "32e8ffca-4021-49a4-971f-330935566af4", "--key", "Content-Type", "--value", "int", "--secret", "nottrue"})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorSecretFlag)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PUT", "variables/32e8ffca-4021-49a4-971f-330935566af4"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--file", "./fixtures/variable.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "🚀 Updated Variable with ID 32e8ffca-4021-49a4-971f-330935566af4\n\n", stdout.String())
	})
}
