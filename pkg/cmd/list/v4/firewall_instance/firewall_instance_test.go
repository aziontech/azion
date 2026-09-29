package firewallinstance

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tblWithInstance string = "ID  NAME    ACTIVE  FUNCTION  \n0   string  true    1         \n"
	tblNoInstance   string = "ID  NAME    ACTIVE  FUNCTION  \n"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("more than one firewall function instance", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/firewalls/1234/functions"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		listCmd.AskInput = func(s string) (string, error) {
			return "1234", nil
		}
		cmd := NewCobraCmd(listCmd, f)

		cmd.SetArgs([]string{"--firewall-id", "1234"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblWithInstance, stdout.String())
	})

	t.Run("list - page 0 should generate an error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/firewalls/1234/functions"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		listCmd.AskInput = func(s string) (string, error) {
			return "1234", nil
		}
		cmd := NewCobraCmd(listCmd, f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--page", "0"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("no firewall function instances", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/firewalls/1234/functions"),
			httpmock.JSONFromFile("./fixtures/noinstance.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		listCmd.AskInput = func(s string) (string, error) {
			return "1234", nil
		}
		cmd := NewCobraCmd(listCmd, f)

		cmd.SetArgs([]string{"--firewall-id", "1234"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, tblNoInstance, stdout.String())
	})

	t.Run("list with details", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/firewalls/1234/functions"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		listCmd.AskInput = func(s string) (string, error) {
			return "1234", nil
		}
		cmd := NewCobraCmd(listCmd, f)

		cmd.SetArgs([]string{"--firewall-id", "1234", "--details"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("ask for firewall id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/firewalls/1234/functions"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		listCmd.AskInput = func(s string) (string, error) {
			return "1234", nil
		}
		cmd := NewCobraCmd(listCmd, f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblWithInstance, stdout.String())
	})
}
