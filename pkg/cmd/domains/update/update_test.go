package update

import (
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

var successResponse string = `
{
	"results": {
	  "id": 1337,
	  "name": "ZA WARUDO",
	  "cnames": [
		"www.test123.com",
		"www.pudim.com"
	   ],
	  "cname_access_only": true,
	  "digital_certificate_id": null,
	  "edge_application_id": 1674767911,
	  "is_active": true,
	  "domain_name": "euxhjonxrr.map.azionedge.net"
	},
	"schema_version": 3
  }
`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("update domain", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "domains/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--domain-id", "1337", "--name", "ATUALIZANDO"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Updated Domain with ID 1337\n", stdout.String())
	})

	t.Run("update all fields", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "domains/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-d", "1337", "--name", "ATUALIZANDO", "--cnames", "www.test.com,www.pudim.com", "--cname-access-only", "false", "--application-id", "1674767911", "--active", "true"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Updated Domain with ID 1337\n", stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "domains/1234"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-d", "1234", "--active", "unactive"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "domains/1337"),
			httpmock.JSONFromString(successResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--in", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, "Updated Domain with ID 1337\n", stdout.String())
	})
}
