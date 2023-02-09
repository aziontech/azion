package create

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/origins"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	t.Run("create new domains", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/origins"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "one piece is the best",
			"--addresses", "asdfsd.cvdf",
			"--host-header", "asdfsdfsd.cvdf",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OriginsCreateOutputSuccess,  92779), stdout.String())
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
      "--in", "./fixtures/create.json",
    })

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OriginsCreateOutputSuccess, 92779), stdout.String())
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
		err := cmd.Execute()
		require.Error(t, err)
	})
}
