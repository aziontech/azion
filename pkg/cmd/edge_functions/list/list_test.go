package list

import (
    // "fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
    tblWithFunc string = "\x1b[34;4mID    NAME             LANGUAGE    ACTIVE  \n\x1b[0m\x1b[32m2995  \x1b[0m20220124-batata  javascript  false   \n\x1b[32m3032  \x1b[0mTestandoCLI4     javascript  false   \n"
    tblNoFunc string = "\x1b[34;4mID    NAME             LANGUAGE    ACTIVE  \n\x1b[0m"
)

func TestList(t *testing.T) {
	t.Run("more than one function", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_functions"),
			httpmock.JSONFromFile("./fixtures/functions.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{})

 		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t,tblWithFunc, stdout.String())
})

	t.Run("no functions", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_functions"),
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
