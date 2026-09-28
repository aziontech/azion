package custompages

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
	tblWithCustomPage string = "ID  NAME    ACTIVE  \n0   string  true    \n"
	tblNoCustomPage   string = "ID  NAME    ACTIVE  \n"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("more than one custom page", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/custom_pages"),
			httpmock.JSONFromFile("./fixtures/custompages.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblWithCustomPage, stdout.String())
	})

	t.Run("list api error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/custom_pages"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--page", "0"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("no custom pages", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/custom_pages"),
			httpmock.JSONFromFile("./fixtures/nocustompage.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblNoCustomPage, stdout.String())
	})
}
