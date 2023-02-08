package describe

import (
	"log"
	"os"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)


func TestDescribe(t *testing.T) {
	t.Run("describe an domains", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/origins.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"-a", "123423424", "-o", "88144"})

		err := cmd.Execute()
		require.NoError(t, err)
  })
	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origin"),
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
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"-a", "123423424", "-o", "88149"})

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("export to a file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/123423424/origins"),
			httpmock.JSONFromFile("./fixtures/origins.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		path := "./out.json"
		cmd.SetArgs([]string{"-a", "123423424", "-o", "88144", "--out", path})

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

		require.Equal(t, `File successfully written to: out.json
`, stdout.String())
	})
}
