package list

import (
	"net/http"
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
			httpmock.REST(http.MethodGet, "api/variables"),
			httpmock.JSONFromFile(".fixtures/variables.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)

		wantOutput := "ID                                    KEY            VALUE  \n32e8ffca-4021-49a4-971f-330935566af4  Content-Type   json   \ne314a185-d775-40f9-9b68-714bbbfbd442  Content-Type2  json   \n"
		assert.Equal(t, wantOutput, stdout.String())
	})

	t.Run("no itens", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "api/variables"),
			httpmock.JSONFromFile(".fixtures/nocontent.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, "ID                                    KEY            VALUE  \n", stdout.String())
	})
}
