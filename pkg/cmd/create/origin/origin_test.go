package origin

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/origin"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("create new origin single origin", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/origins"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "onepieceisthebest",
			"--addresses", "asdfsd.cvdf",
			"--host-header", "asdfsdfsd.cvdf",
			"--origin-type", "single_origin",
			"--origin-protocol-policy", "asdf",
			"--origin-path", "asdf",
			"--hmac-region-name", "asdf",
			"--hmac-access-key", "asdf",
			"--hmac-secret-key", "asdf",
			"--hmac-authentication", "true",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, "176fa5e2-a895-4862-9657-e2e37d9125a7"), stdout.String())
	})

	t.Run("create new origin object storage", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/origins"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "onepieceisthebest",
			"--addresses", "asdfsd.cvdf",
			"--host-header", "asdfsdfsd.cvdf",
			"--origin-type", "object_storage",
			"--origin-protocol-policy", "asdf",
			"--origin-path", "asdf",
			"--hmac-region-name", "asdf",
			"--hmac-access-key", "asdf",
			"--hmac-secret-key", "asdf",
			"--hmac-authentication", "true",
			"--bucket", "truebucket",
			"--prefix", "123123211123213",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, "176fa5e2-a895-4862-9657-e2e37d9125a7"), stdout.String())
	})

	t.Run("create with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/origins"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--file", "./fixtures/create.json",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, "176fa5e2-a895-4862-9657-e2e37d9125a7"), stdout.String())
	})

	t.Run("bad request status 400", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/origin"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("internal server error 500", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/origin"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "onepieceisthebest",
			"--addresses", "asdfsd.cvdf",
			"--host-header", "asdfsdfsd.cvdf",
			"--origin-type", "single_origin",
		})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("error file flag not exist", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/origins"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--file", "./fixtures/no_exist.json",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, utils.ErrorUnmarshalReader)
	})

	t.Run("Error field expected bool, not is bool", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/origins"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "onepieceisthebest",
			"--addresses", "asdfsd.cvdf",
			"--host-header", "asdfsdfsd.cvdf",
			"--origin-type", "single_origin",
			"--origin-protocol-policy", "asdf",
			"--origin-path", "asdf",
			"--hmac-region-name", "asdf",
			"--hmac-access-key", "asdf",
			"--hmac-secret-key", "asdf",
			"--hmac-authentication", "1232132",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorHmacAuthenticationFlag)
	})
}
