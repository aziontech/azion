package csr

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/create/csr"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

var generatedCSR string = "-----BEGIN CERTIFICATE REQUEST-----MIIBdummycsr==-----END CERTIFICATE REQUEST-----"

var successResponse string = `
{
  "state": "executed",
  "data": {
    "id": 1337,
    "name": "My CSR",
    "certificate": "string",
    "issuer": "string",
    "subject_name": ["example.com"],
    "validity": "string",
    "type": "certificate",
    "managed": false,
    "status": "pending",
    "status_detail": "",
    "csr": "-----BEGIN CERTIFICATE REQUEST-----MIIBdummycsr==-----END CERTIFICATE REQUEST-----",
    "challenge": "dns",
    "authority": "lets_encrypt",
    "key_algorithm": "rsa_2048",
    "active": true,
    "product_version": "1.0",
    "last_editor": "user@example.com",
    "created_at": "2019-08-24T14:15:22Z",
    "last_modified": "2019-08-24T14:15:22Z",
    "renewed_at": "2019-08-24T14:15:22Z"
  }
}
`

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("create csr with flags", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/tls/csr"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--name", "My CSR",
			"--common-name", "example.com",
			"--country", "US",
			"--state", "California",
			"--locality", "San Francisco",
			"--organization", "Example Corp",
			"--organization-unity", "IT",
			"--email", "admin@example.com",
			"--alternative-names", "www.example.com,api.example.com",
			"--key-algorithm", "rsa_2048",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		// The success message and the generated CSR content are both printed.
		assert.Contains(t, stdout.String(), fmt.Sprintf(msg.CreateOutputSuccess, 1337))
		assert.Contains(t, stdout.String(), generatedCSR)
	})

	t.Run("create csr from file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/tls/csr"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/create.json"})

		err := cmd.Execute()
		require.NoError(t, err)
		assert.Contains(t, stdout.String(), fmt.Sprintf(msg.CreateOutputSuccess, 1337))
	})

	t.Run("invalid json file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--file", "./fixtures/does-not-exist.json"})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorInvalidJSON)
	})

	t.Run("api returns bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("POST", "workspace/tls/csr"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Bad Request"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--name", "My CSR",
			"--common-name", "example.com",
			"--country", "US",
			"--state", "California",
			"--locality", "San Francisco",
			"--organization", "Example Corp",
			"--organization-unity", "IT",
			"--email", "admin@example.com",
		})

		err := cmd.Execute()
		require.Error(t, err)
	})
}
