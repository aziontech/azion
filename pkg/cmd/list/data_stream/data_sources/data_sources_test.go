package datasources

import (
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tblOneDataSource   string = "SLUG  NAME    ACTIVE  \nhttp  string  true    \n"
	tblManyDataSources string = "SLUG  NAME       ACTIVE  \nhttp  string     true    \nrtld  real-time  false   \n"
	tblNoDataSource    string = "SLUG  NAME    ACTIVE  \n"
)

func TestList(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("single data source", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/data_sources"),
			httpmock.JSONFromFile("./fixtures/onedatasource.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblOneDataSource, stdout.String())
	})

	t.Run("more than one data source", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/data_sources"),
			httpmock.JSONFromFile("./fixtures/data_sources.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblManyDataSources, stdout.String())
	})

	t.Run("list with details renders the same columns", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/data_sources"),
			httpmock.JSONFromFile("./fixtures/onedatasource.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--details"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblOneDataSource, stdout.String())
	})

	t.Run("passes ordering and pagination flags", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/data_sources"),
			httpmock.JSONFromFile("./fixtures/onedatasource.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--order-by", "slug", "--page", "2", "--page-size", "10", "--sort", "asc"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblOneDataSource, stdout.String())
	})

	t.Run("list api error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/data_sources"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{"--page", "0"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("no data sources", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "workspace/stream/data_sources"),
			httpmock.JSONFromFile("./fixtures/nodatasource.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
		assert.Equal(t, tblNoDataSource, stdout.String())
	})
}
