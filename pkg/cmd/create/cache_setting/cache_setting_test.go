package cachesetting

import (
	"fmt"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/aziontech/azion-cli/utils"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("create new cache_setting", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_application/applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_application/applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "BetterLesson",
			"--adaptive-delivery-action", "ignore",

			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",

			"--slice-configuration-enabled", "false",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 0), stdout.String())
	})

	t.Run("create with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_application/applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_application/applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--file", "./fixtures/create.json",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.CreateOutputSuccess, 0), stdout.String())
	})

	t.Run("wrong caching for post boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_application/applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_application/applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "BetterLesson",
			"--adaptive-delivery-action", "ignore",

			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "incorrect",
			"--enable-caching-string-sort", "true",

			"--slice-configuration-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorCachingForPostFlag)
	})

	t.Run("wrong caching string sort boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_application/applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_application/applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "BetterLesson",
			"--adaptive-delivery-action", "ignore",

			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "false",
			"--enable-caching-for-post", "false",
			"--enable-caching-string-sort", "precise",

			"--slice-configuration-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorCachingStringSortFlag)
	})

	t.Run("wrong slice configuration enabled boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_application/applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_application/applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--name", "BetterLesson",
			"--adaptive-delivery-action", "ignore",

			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",

			"--slice-configuration-enabled", "faithful",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, msg.ErrorSliceConfigurationFlag)
	})

	t.Run("error unmarshall", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_application/applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_application/applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--file", "./fixtures/error",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, utils.ErrorUnmarshalReader)
	})

}
