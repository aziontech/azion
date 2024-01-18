package origin

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("list page 1", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/origins.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "123423424"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,"ORIGIN KEY                            NAME            \n0cee30cd-1743-4202-b0dd-da9b636a6035  Default Origin  \ne4f0761b-d2ac-4168-aa4b-f525d08396fd  Create Origin   \n", stdout.String())
	})

	t.Run("no itens", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/noorigins.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "123423424"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		assert.Equal(t,"ORIGIN KEY                            NAME            \n", stdout.String())
	})
}
