package templates

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
	tblWithTemplate    string = "ID  NAME    ACTIVE  \n0   string  true    \n"
	tblNoTemplate      string = "ID  NAME    ACTIVE  \n"
	tblDetailsTemplate string = "ID  NAME    ACTIVE  CUSTOM  LAST EDITOR       LAST MODIFIED                  \n0   string  true    true    user@example.com  2019-08-24 14:15:22 +0000 UTC  \n"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("more than one template", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/templates"),
			httpmock.JSONFromFile("./fixtures/templates.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblWithTemplate, stdout.String())
	})

	t.Run("list with details", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/templates"),
			httpmock.JSONFromFile("./fixtures/templates.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--details"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblDetailsTemplate, stdout.String())
	})

	t.Run("list api error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/templates"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--page", "0"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("no templates", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/templates"),
			httpmock.JSONFromFile("./fixtures/notemplate.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblNoTemplate, stdout.String())
	})
}
