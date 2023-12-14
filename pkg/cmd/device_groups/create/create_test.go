package create

// import (
// 	"fmt"
// 	"github.com/aziontech/azion-cli/pkg/logger"
// 	"go.uber.org/zap/zapcore"
// 	"net/http"
// 	"testing"

// 	msg "github.com/aziontech/azion-cli/messages/device_groups"
// 	"github.com/aziontech/azion-cli/pkg/httpmock"
// 	"github.com/aziontech/azion-cli/pkg/testutils"
// 	"github.com/stretchr/testify/require"
// )

// func TestCreate(t *testing.T) {
// 	logger.New(zapcore.DebugLevel)
// 	t.Run("create new device groups", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("POST", "edge_applications/1673635846/device_groups"),
// 			httpmock.JSONFromFile(".fixtures/response.json"),
// 		)

// 		f, stdout, _ := testutils.NewFactory(mock)
// 		cmd := NewCmd(f)
// 		cmd.SetArgs([]string{
// 			"--application-id", "1673635846",
// 			"--name", "Mob123il22e",
// 			"--user-agent", "Mobile|Android|iPhone",
// 		})

// 		err := cmd.Execute()
// 		require.NoError(t, err)
// 		require.Equal(t, fmt.Sprintf(msg.DeviceGroupsCreateOutputSuccess, 2260), stdout.String())
// 	})

// 	t.Run("create with file", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("POST", "edge_applications/1673635846/device_groups"),
// 			httpmock.JSONFromFile(".fixtures/response.json"),
// 		)

// 		f, stdout, _ := testutils.NewFactory(mock)

// 		cmd := NewCmd(f)
// 		cmd.SetArgs([]string{
// 			"--application-id", "1673635846",
// 			"--in", ".fixtures/request.json",
// 		})

// 		err := cmd.Execute()
// 		require.NoError(t, err)
// 		require.Equal(t, fmt.Sprintf(msg.DeviceGroupsCreateOutputSuccess, 2260), stdout.String())
// 	})

// 	t.Run("bad request status 400", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("POST", "edge_applications/1673635846/device_groups"),
// 			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
// 		)

// 		f, _, _ := testutils.NewFactory(mock)

// 		cmd := NewCmd(f)
// 		err := cmd.Execute()
// 		require.Error(t, err)
// 	})

// 	t.Run("internal server error 500", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("POST", "edge_applications/1673635846/device_groups"),
// 			httpmock.StatusStringResponse(http.StatusInternalServerError, "Invalid"),
// 		)

// 		f, _, _ := testutils.NewFactory(mock)

// 		cmd := NewCmd(f)
// 		err := cmd.Execute()
// 		require.Error(t, err)
// 	})
// }
