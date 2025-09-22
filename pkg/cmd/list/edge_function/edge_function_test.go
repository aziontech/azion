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

var (
	tblWithFunc string = "ID    NAME    ACTIVE  \n1337  string  true    \n"
	tblNoFunc   string = "ID    NAME    ACTIVE  \n"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("more than one function", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_functions/functions"),
			httpmock.JSONFromFile("./fixtures/functions.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblWithFunc, stdout.String())
	})

	t.Run("list - page 0 should generate an error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_functions/functions"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--page", "0"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("no functions", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_functions/functions"),
			httpmock.JSONFromFile("./fixtures/nofunction.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t, tblNoFunc, stdout.String())
	})
}
