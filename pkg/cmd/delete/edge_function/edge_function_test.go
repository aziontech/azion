package edgefunction

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("delete function by id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_functions/1234"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--function-id", "1234"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, "Edge Function 1234 was successfully deleted\n", stdout.String())
	})

	t.Run("delete function that is not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_functions/1234"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--function-id", "1234"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
