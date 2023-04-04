package list

import (
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewCmd(t *testing.T) {
	t.Run("list all edge functions instances", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "/edge_applications/12312333/functions_instances"),
			httpmock.JSONFromFile("./.fixtures/resp.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{"-a", "12312333"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("no itens", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "/edge_applications/12312333/functions_instances"),
			httpmock.JSONFromFile("./.fixtures/no_resp.json"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
