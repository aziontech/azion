package rulesengine

import (
	"log"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDescribe(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("describe a rule engine", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1678743802/rules_engine/request/rules/173617"),
			httpmock.JSONFromFile("./fixtures/rules.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		err := cmd.Execute()
		if err != nil {
			log.Println("error executing cmd err: ", err.Error())
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1678743802/rules_engine/request/rules/666"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		err := cmd.Execute()
		require.Error(t, err)
	})

	t.Run("missing mandatory flag", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("GET", "edge_applications/1678743802/rules_engine/request/rules/1"),
			httpmock.StatusStringResponse(http.StatusNotFound, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{})

		err := cmd.Execute()
		require.Error(t, err)
	})
}
