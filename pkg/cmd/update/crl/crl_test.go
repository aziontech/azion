package crl

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/update/crl"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var successResponse string = `
{
  "state": "executed",
  "data": {
    "id": 1337,
    "name": "My Updated CRL",
    "active": true,
    "last_editor": "user@example.com",
    "created_at": "2019-08-24T14:15:22Z",
    "last_modified": "2019-08-24T14:15:22Z",
    "product_version": "1.0",
    "issuer": "My CA",
    "last_update": "2019-08-24T14:15:22Z",
    "next_update": "2019-08-24T14:15:22Z",
    "crl": "-----BEGIN X509 CRL-----\nMIIBdummy==\n-----END X509 CRL-----\n"
  }
}
`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("update crl with flags", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/tls/crls/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--crl-id", "1337", "--name", "My Updated CRL", "--issuer", "My CA", "--crl", "./fixtures/crl.pem", "--active", "true"})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("update crl from file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/tls/crls/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--crl-id", "1337", "--file", "./fixtures/update.json"})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, 1337), stdout.String())
	})

	t.Run("invalid active flag", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--crl-id", "1337", "--active", "notabool"})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorActiveFlag)
	})

	t.Run("crl file does not exist", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--crl-id", "1337", "--crl", "./fixtures/does-not-exist.pem"})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorReadCRLFile)
	})

	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/tls/crls/9999"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--crl-id", "9999", "--name", "My Updated CRL"})

		err := cmd.Execute()
		require.Error(t, err)
	})
}
