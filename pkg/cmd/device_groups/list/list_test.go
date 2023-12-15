package list

// import (
// 	"github.com/aziontech/azion-cli/pkg/httpmock"
// 	"github.com/aziontech/azion-cli/pkg/logger"
// 	"github.com/aziontech/azion-cli/pkg/testutils"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"go.uber.org/zap/zapcore"
// 	"testing"
// )

// func TestList(t *testing.T) {
// 	logger.New(zapcore.DebugLevel)
// 	t.Run("command list with successes", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("GET", "edge_applications/1673635846/device_groups"),
// 			httpmock.JSONFromFile(".fixtures/resp.json"),
// 		)

// 		f, stdout, _ := testutils.NewFactory(mock)
// 		cmd := NewCmd(f)

// 		cmd.SetArgs([]string{"-a", "1673635846"})

// 		_, err := cmd.ExecuteC()
// 		require.NoError(t, err)
// 		assert.Equal(t, "ID    NAME    \n2257  Mob1le  \n", stdout.String())
// 	})

// 	t.Run("command list response without items", func(t *testing.T) {
// 		mock := &httpmock.Registry{}

// 		mock.Register(
// 			httpmock.REST("GET", "edge_applications/1673635847/device_groups"),
// 			httpmock.JSONFromFile(".fixtures/resp_without_items.json"),
// 		)

// 		f, stdout, _ := testutils.NewFactory(mock)
// 		cmd := NewCmd(f)

// 		cmd.SetArgs([]string{"-a", "1673635847"})

// 		_, err := cmd.ExecuteC()
// 		require.NoError(t, err)
// 		assert.Equal(t, "ID    NAME    \n", stdout.String())
// 	})
// }
