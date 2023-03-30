package update

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/rules_engine"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	t.Run("update ", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/rules_engine/request/rules/1234"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--phase", "request",
			"--in", "./fixtures/update.json",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.RulesEngineUpdateOutputSuccess, 1234), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST(http.MethodPatch, "edge_applications/1673635839/rules_engine/request/rules/1234"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--application-id", "1673635839",
			"--rule-id", "1234",
			"--phase", "request",
		})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorMandatoryFlagsUpdate)
	})
}
