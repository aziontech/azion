package delete

import (
	"fmt"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/origins"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {

	t.Run("delete domain by id", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1673635839/origins/58755fef-e830-4ea4-b9e0-6481f1ef496d"),
			httpmock.StatusStringResponse(204, ""),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--application-id", "1673635839", "--origin-key", "58755fef-e830-4ea4-b9e0-6481f1ef496d"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf(msg.OriginsDeleteOutputSuccess, "58755fef-e830-4ea4-b9e0-6481f1ef496d"), stdout.String())
	})

	t.Run("delete domain that is not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("DELETE", "edge_applications/1673635839/origins/58755fef-e830-4ea4-b9e0-6481f1ef497a"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{"-d", "1234"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})
}
