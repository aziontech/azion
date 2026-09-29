package firewallinstance

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/create/firewall_instance"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("create new Firewall Function Instance", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/firewalls/1234/functions"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--function-id", "5678", "--name", "My Instance", "--active", "true"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 0), stdout.String())
	})

	t.Run("create with args", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/firewalls/1234/functions"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--firewall-id", "1234",
			"--function-id", "5678",
			"--name", "Instance with Args",
			"--active", "true",
			"--args", "./fixtures/args.json",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 0), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/firewalls/1234/functions"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--function-id", "5678", "--name", "Test Instance", "--active", "true"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("create with file", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/firewalls/1234/functions"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--file", "./fixtures/create.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 0), stdout.String())
	})

	t.Run("Error Field active not is boolean", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/firewalls/1234/functions"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--function-id", "5678", "--name", "Test Instance", "--active", "invalid_value"})

		err := cmd.Execute()
		stringErr := fmt.Sprintf("%s: %q", msg.ErrorIsActiveFlag, "invalid_value")
		if stringErr == err.Error() {
			return
		}
		t.Fatalf("Error: %q", err)
	})

	t.Run("Error create firewall instance request api", func(t *testing.T) {
		t.Parallel()

		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "workspace/firewalls/1234/functions"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, "Internal Server Error"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--function-id", "5678", "--name", "Test Instance", "--active", "true"})

		err := cmd.Execute()
		if err != nil {
			return
		}
		t.Fatalf("Error: %q", err)
	})
}
