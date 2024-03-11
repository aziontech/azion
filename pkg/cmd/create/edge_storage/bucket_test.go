package edge_storage

import (
	"github.com/stretchr/testify/require"
	"net/http"
	// "reflect"
	"testing"

	// "github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
	// "github.com/spf13/pflag"
)

func TestNewBucket(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("commando bucket --help ", func(t *testing.T) {
		mock := &httpmock.Registry{}
		request := httpmock.REST(http.MethodPost, "")
		response := httpmock.JSONResponse("")
		mock.Register(request, response)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewBucket(f)
		cmd.SetArgs([]string{"--help"})
		err := cmd.Execute()
		require.NoError(t, err)
	})

	//Error testing
	//t.Run("commando bucket new item", func(t *testing.T) {
	//	mock := &httpmock.Registry{}
	//	request := httpmock.REST("POST", "v4/storage/buckets")
	//	response := httpmock.JSONFromFile("")
	//	mock.Register(request, response)
	//
	//	f, _, _ := testutils.NewFactory(mock)
	//	cmd := NewBucket(f)
	//	cmd.SetArgs([]string{
	//		"--name", "arthur-morgan",
	//		"--edge-access", "read_only",
	//	})
	//
	//	err := cmd.Execute()
	//	require.NoError(t, err)
	//})
}
