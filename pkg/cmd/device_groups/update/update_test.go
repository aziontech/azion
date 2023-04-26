package update

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/device_groups"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Run("update device group", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/device_groups/1234"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--group-id", "1234",
			"--name", "shokugeki",
			"--user-agent", "Mobile|iP(hone|od)|BlackBerry|IEMobile",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DeviceGroupsUpdateOutputSuccess, 1234), stdout.String())
	})

	t.Run("missing flags", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/device_groups/1234"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--name", "shokugeki",
			"--user-agent", "Mobile|iP(hone|od)|BlackBerry|IEMobile",
		})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorMandatoryFlagsUpdate)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/device_groups/1234"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--group-id", "1234",
			"--in", "./fixtures/update.json",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.DeviceGroupsUpdateOutputSuccess, 1234), stdout.String())
	})
}
