package edgeapplication

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("delete application by id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1234"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "1234"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, "Edge Application 1234 was successfully deleted\n", stdout.String())
	})

	t.Run("delete application - not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1234"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "1234"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("delete cascade with no azion.json file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--cascade"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
	t.Run("cascade delete application", func(t *testing.T) {
		mock := &httpmock.Registry{}
		options := &contracts.AzionApplicationOptions{}

		dat, _ := os.ReadFile("./fixtures/azion.json")
		_ = json.Unmarshal(dat, options)

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/666"),
			httpmock.StatusStringResponse(204, ""),
		)
		mock.Register(
			httpmock.REST("DELETE", "edge_functions/123"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		del := NewDeleteCmd(f)
		del.GetAzion = func() (*contracts.AzionApplicationOptions, error) {
			return options, nil
		}
		del.UpdateJson = func(cmd *DeleteCmd) error {
			return nil
		}
		del.f = f
		del.Io = f.IOStreams

		cmd := NewCobraCmd(del)

		cmd.SetArgs([]string{"--cascade"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, "Cascade delete carried out successfully\n", stdout.String())
	})
}
