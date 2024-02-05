package variables

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("delete variable by id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "variables/7a187044-4a00-4a4a-93ed-d230900421f3"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--variable-id", "7a187044-4a00-4a4a-93ed-d230900421f3"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, "🚀 Variable 7a187044-4a00-4a4a-93ed-d230900421f3 was successfully deleted\n\n", stdout.String())
	})

	t.Run("delete variable that is not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "variables/7a187044-4a00-4a4a-93ed-d230900421f3"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--variable-id", "7a187044-4a00-4a4a-93ed-d230900421f3"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
