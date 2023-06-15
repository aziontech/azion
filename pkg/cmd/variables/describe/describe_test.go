package describe

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestDescribe(t *testing.T) {
	t.Run("describe an variables", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "api/variables/32e8ffca-4021-49a4-971f-330935566af4"),
			httpmock.JSONFromFile(".fixtures/variables.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"-v", "32e8ffca-4021-49a4-971f-330935566af4"})

		err := cmd.Execute()
		require.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "api/variables/32e8ffca-4021-49a4-971f-330935566af4"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("no id sent", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "api/variables/123423424"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"-v", "123423424"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("export to a file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "api/variables/32e8ffca-4021-49a4-971f-330935566af4"),
			httpmock.JSONFromFile(".fixtures/variables.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		path := "./out.json"
		cmd.SetArgs([]string{"-v", "32e8ffca-4021-49a4-971f-330935566af4", "--out", path})

		err := cmd.Execute()
		if err != nil {
			log.Println("error executing cmd err: ", err.Error())
		}

		_, err = os.ReadFile(path)
		if err != nil {
			t.Fatalf("error reading `out.json`: %v", err)
		}
		defer func() {
			_ = os.Remove(path)
		}()

		require.NoError(t, err)
		fmt.Println(">>> ", stdout.String())
		require.Equal(t, `File successfully written to: out.json`, stdout.String())
	})
}
