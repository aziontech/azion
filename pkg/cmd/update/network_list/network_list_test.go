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

var successResponse string = `
{
  "data": {
    "id": 1337,
    "name": "Updated Network List",
    "type": "ip_cidr",
    "items": [
      "192.168.1.0/24",
      "10.0.0.0/8"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("update Network List name and active", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--name", "Updated Network List", "--active", "true"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("update type and items", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--type", "ip_cidr", "--items", "192.168.1.0/24,10.0.0.0/8"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("update only name", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--name", "Updated Network List"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("update only active status", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--active", "false"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("bad request - invalid active value", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1234"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1234", "--active", "invalid"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/9999"),
			httpmock.StatusStringResponse(http.StatusNotFound, `{"details": "Network List not found"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "9999", "--name", "Test"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("add single item", func(t *testing.T) {
		mock := &httpmock.Registry{}

		// Mock GET to retrieve current items
		getCurrentResponse := `
{
  "data": {
    "id": 1337,
    "name": "Test Network List",
    "type": "ip_cidr",
    "items": [
      "192.168.1.0/24",
      "10.0.0.0/8"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`
		mock.Register(
			httpmock.REST("GET", "workspace/network_lists/1337"),
			httpmock.JSONFromString(getCurrentResponse),
		)

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--add-item", "203.0.113.0/24"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("add multiple items", func(t *testing.T) {
		mock := &httpmock.Registry{}

		getCurrentResponse := `
{
  "data": {
    "id": 1337,
    "name": "Test Network List",
    "type": "ip_cidr",
    "items": [
      "192.168.1.0/24"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`
		mock.Register(
			httpmock.REST("GET", "workspace/network_lists/1337"),
			httpmock.JSONFromString(getCurrentResponse),
		)

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--add-item", "203.0.113.0/24,172.16.0.0/12"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("remove single item", func(t *testing.T) {
		mock := &httpmock.Registry{}

		getCurrentResponse := `
{
  "data": {
    "id": 1337,
    "name": "Test Network List",
    "type": "ip_cidr",
    "items": [
      "192.168.1.0/24",
      "10.0.0.0/8",
      "203.0.113.0/24"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`
		mock.Register(
			httpmock.REST("GET", "workspace/network_lists/1337"),
			httpmock.JSONFromString(getCurrentResponse),
		)

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--remove-item", "10.0.0.0/8"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("remove multiple items", func(t *testing.T) {
		mock := &httpmock.Registry{}

		getCurrentResponse := `
{
  "data": {
    "id": 1337,
    "name": "Test Network List",
    "type": "ip_cidr",
    "items": [
      "192.168.1.0/24",
      "10.0.0.0/8",
      "203.0.113.0/24",
      "172.16.0.0/12"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`
		mock.Register(
			httpmock.REST("GET", "workspace/network_lists/1337"),
			httpmock.JSONFromString(getCurrentResponse),
		)

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--remove-item", "10.0.0.0/8,172.16.0.0/12"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("add and remove items together", func(t *testing.T) {
		mock := &httpmock.Registry{}

		getCurrentResponse := `
{
  "data": {
    "id": 1337,
    "name": "Test Network List",
    "type": "ip_cidr",
    "items": [
      "192.168.1.0/24",
      "10.0.0.0/8"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`
		mock.Register(
			httpmock.REST("GET", "workspace/network_lists/1337"),
			httpmock.JSONFromString(getCurrentResponse),
		)

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--add-item", "203.0.113.0/24", "--remove-item", "10.0.0.0/8"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("add item that already exists", func(t *testing.T) {
		mock := &httpmock.Registry{}

		getCurrentResponse := `
{
  "data": {
    "id": 1337,
    "name": "Test Network List",
    "type": "ip_cidr",
    "items": [
      "192.168.1.0/24",
      "10.0.0.0/8"
    ],
    "last_editor": "user@example.com",
    "last_modified": "2019-08-24T14:15:22Z",
    "active": true
  }
}
`
		mock.Register(
			httpmock.REST("GET", "workspace/network_lists/1337"),
			httpmock.JSONFromString(getCurrentResponse),
		)

		mock.Register(
			httpmock.REST("PATCH", "workspace/network_lists/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--network-list-id", "1337", "--add-item", "192.168.1.0/24"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})
}
