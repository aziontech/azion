package networklist

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/network_list"
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
      "detail": "string",
      "source": {
        "pointer": "string",
        "parameter": "string",
        "header": "string"
      },
      "meta": {
        "property1": null,
        "property2": null
      }
    }
  ]
}
`

var successResponse string = `
{
  "data": {
    "id": 1337,
    "name": "My IP List",
    "type": "ip_cidr",
    "items": [
      "192.168.0.1/32",
      "10.0.0.0/8"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("create new Network List", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/network_lists"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "My IP List", "--type", "ip_cidr", "--items", "192.168.0.1/32,10.0.0.0/8", "--active", "true"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("create ASN Network List", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/network_lists"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "ASN List", "--type", "asn", "--items", "1234,5678"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("create countries Network List", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/network_lists"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "Countries List", "--type", "countries", "--items", "US,BR,JP"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/network_lists"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "Bad List", "--type", "invalid_type", "--items", "test"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("create with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/network_lists"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--file", "./fixtures/create.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("Error Field active not is boolean", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/network_lists"),
			httpmock.JSONFromString(successResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "Test List", "--type", "ip_cidr", "--items", "192.168.1.0/24", "--active", "invalid_bool"})

		err := cmd.Execute()
		stringErr := fmt.Sprintf("%s: %s", msg.ErrorActiveFlag, "invalid_bool")
		if stringErr == err.Error() {
			return
		}
		t.Fatalf("Error: %q", err)
	})

	t.Run("Error create network list request api", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/network_lists"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, errorResponse),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--name", "Error List", "--type", "ip_cidr", "--items", "192.168.1.0/24", "--active", "true"})

		err := cmd.Execute()
		if err != nil {
			return
		}
		t.Fatalf("Error: %q", err)
	})
}
