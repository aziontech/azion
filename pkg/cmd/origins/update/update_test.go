package update

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/origins"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Run("update edge application", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/origins/03a6e7bf-8e26-49c7-a66e-ab8eaa425086"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--origin-key", "03a6e7bf-8e26-49c7-a66e-ab8eaa425086",
			"--name", "onepiece",
			"--addresses", "asdfsd.cvdf",
			"--host-header", "asdfsdfsd.cvdf",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OriginsUpdateOutputSuccess, "03a6e7bf-8e26-49c7-a66e-ab8eaa425086"), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/origins/03a6e7bf-8e26-49c7-a66e-ab8eaa425086"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--origin-key", "03a6e7bf-8e26-49c7-a66e-ab8eaa425086",
			"--name", "onepiece",
			"--addresses", "asdfsd.cvdf",
			"--host-header", "asdfsdfsd.cvdf",
		})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/origins/03a6e7bf-8e26-49c7-a66e-ab8eaa425086"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--origin-key", "03a6e7bf-8e26-49c7-a66e-ab8eaa425086",
			"--in", "./fixtures/update.json",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OriginsUpdateOutputSuccess, "03a6e7bf-8e26-49c7-a66e-ab8eaa425086"), stdout.String())
	})
}
