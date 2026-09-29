package crl

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/create/crl"
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
    "name": "My CRL",
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

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("create new CRL with flags", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/tls/crls"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "My CRL", "--issuer", "My CA", "--crl", "./fixtures/crl.pem", "--active", "true"})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("create new CRL from file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/tls/crls"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/create.json"})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 1337), stdout.String())
	})

	t.Run("crl file does not exist", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "My CRL", "--issuer", "My CA", "--crl", "./fixtures/does-not-exist.pem"})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorReadCRLFile)
	})

	t.Run("api returns bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/tls/crls"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Bad Request"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--name", "My CRL", "--issuer", "My CA", "--crl", "./fixtures/crl.pem"})

		err := cmd.Execute()
		require.Error(t, err)
	})
}
