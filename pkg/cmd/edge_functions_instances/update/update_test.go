package update

import (
	"fmt"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/edge_functions_instances"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("update instance", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1678743802/functions_instances/9810"),
			httpmock.JSONFromFile("./fixtures/instance.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1678743802",
			"--instance-id", "9810",
			"--function-id", "8065",
			"--name", "updated",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.EdgeFuncInstanceUpdateOutputSuccess, 9810), stdout.String())
	})

	t.Run("missing flags", func(t *testing.T) {

		f, _, _ := testutils.NewFactory(nil)

		cmd := NewCmd(f)
		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorMandatoryUpdateFlags)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/functions_instances/9810"),
			httpmock.JSONFromFile("./fixtures/instance.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--in", "./fixtures/update.json",
			"--instance-id", "9810",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.EdgeFuncInstanceUpdateOutputSuccess, 9810), stdout.String())
	})
}
