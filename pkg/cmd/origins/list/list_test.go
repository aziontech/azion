package list

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("list page 1", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/origins.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "123423424"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, "ID     NAME            \n88144  Default Origin  \n91799  Create Origin   \n", stdout.String())
	})

	t.Run("no itens", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/noorigins.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "123423424"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, "ID     NAME            \n", stdout.String())
	})
}
