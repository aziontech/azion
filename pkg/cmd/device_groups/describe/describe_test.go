package describe

// import (
// 	"github.com/aziontech/azion-cli/pkg/logger"
// 	"go.uber.org/zap/zapcore"
// 	"log"
// 	"net/http"
// 	"os"
// 	"testing"

// 	msg "github.com/aziontech/azion-cli/messages/device_groups"
// 	"github.com/aziontech/azion-cli/pkg/httpmock"
// 	"github.com/aziontech/azion-cli/pkg/testutils"
// 	"github.com/stretchr/testify/require"
// )

// func TestDescribe(t *testing.T) {
// 	logger.New(zapcore.DebugLevel)
// 	t.Run("describe a device group", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("GET", "edge_applications/1676400693/device_groups/2259"),
// 			httpmock.JSONFromFile("./fixtures/groups.json"),
// 		)

// 		f, _, _ := testutils.NewFactory(mock)

// 		cmd := NewCmd(f)
// 		cmd.SetArgs([]string{"-a", "1676400693", "-g", "2259"})

// 		err := cmd.Execute()
// 		require.NoError(t, err)
// 	})
// 	t.Run("not found", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("GET", "edge_applications/1676400693/device_groups/666"),
// 			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
// 		)

// 		f, _, _ := testutils.NewFactory(mock)

// 		cmd := NewCmd(f)

// 		err := cmd.Execute()
// 		require.Error(t, err)
// 	})

// 	t.Run("missing mandatory flag", func(t *testing.T) {
// 		mock := &httpmock.Registry{}
// 		mock.Register(
// 			httpmock.REST("GET", "edge_applications/1676400693/device_groups/2259"),
// 			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
// 		)

// 		f, _, _ := testutils.NewFactory(mock)
// 		cmd := NewCmd(f)
// 		cmd.SetArgs([]string{})

// 		err := cmd.Execute()
// 		require.ErrorIs(t, err, msg.ErrorMandatoryFlags)
// 	})

// 	t.Run("export to a file", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("GET", "edge_applications/1676400693/device_groups/2259"),
// 			httpmock.JSONFromFile("./fixtures/groups.json"),
// 		)

// 		f, stdout, _ := testutils.NewFactory(mock)

// 		cmd := NewCmd(f)
// 		path := "./out.json"
// 		cmd.SetArgs([]string{"-a", "1676400693", "-g", "2259", "--out", path})

// 		err := cmd.Execute()
// 		if err != nil {
// 			log.Println("error executing cmd err: ", err.Error())
// 		}

// 		_, err = os.ReadFile(path)
// 		if err != nil {
// 			t.Fatalf("error reading `out.json`: %v", err)
// 		}
// 		defer func() {
// 			_ = os.Remove(path)
// 		}()

// 		require.NoError(t, err)

// 		require.Equal(t, `File successfully written to: out.json
// `, stdout.String())
// 	})
// }
