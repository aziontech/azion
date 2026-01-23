package firewallinstance

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/update/firewall_instance"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("update Firewall Function Instance", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/functions/5678"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--instance-id", "5678", "--name", "Updated Instance", "--active", "false"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 0), stdout.String())
	})

	t.Run("update with function-id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/functions/5678"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--firewall-id", "1234",
			"--instance-id", "5678",
			"--function-id", "9999",
			"--name", "Updated Instance",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 0), stdout.String())
	})

	t.Run("update with args", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/functions/5678"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--firewall-id", "1234",
			"--instance-id", "5678",
			"--args", "./fixtures/args.json",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 0), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/functions/5678"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--instance-id", "5678", "--active", "invalid"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/functions/5678"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--instance-id", "5678", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 0), stdout.String())
	})

	t.Run("Error Field active not is boolean", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/firewalls/1234/functions/5678"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--instance-id", "5678", "--active", "not_a_bool"})

		err := cmd.Execute()
		stringErr := fmt.Sprintf("%s: %q", msg.ErrorIsActiveFlag, "not_a_bool")
		if stringErr == err.Error() {
			return
		}
		t.Fatalf("Error: %q", err)
	})
}
