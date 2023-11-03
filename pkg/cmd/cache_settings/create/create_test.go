package create

import (
	"fmt"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/messages/cache_settings"
	"go.uber.org/zap/zapcore"
	"testing"

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
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)
		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "false",
			"--slice-l2-caching-enabled", "false",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(cache_settings.CacheSettingsCreateOutputSuccess, 115255), stdout.String())
	})

	t.Run("create with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--in", "./fixtures/create.json",
		})

		err := cmd.Execute()
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(cache_settings.CacheSettingsCreateOutputSuccess, 115255), stdout.String())
	})

	t.Run("no acceleration error --in flag", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result_no_acceleration.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--in", "./fixtures/create.json",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorApplicationAccelerationNotEnabled)
	})

	t.Run("no acceleration error with args", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result_no_acceleration.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "false",
			"--slice-l2-caching-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorApplicationAccelerationNotEnabled)
	})

	t.Run("override but no ttl", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "false",
			"--slice-l2-caching-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorBrowserMaximumTtlNotSent)
	})

	t.Run("no acceleration error with args", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result_no_acceleration.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "false",
			"--slice-l2-caching-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorApplicationAccelerationNotEnabled)
	})

	t.Run("wrong l2 boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "false",
			"--slice-l2-caching-enabled", "troo",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorSliceL2CachingFlag)
	})

	t.Run("wrong caching for options boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "untrue",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "false",
			"--slice-l2-caching-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorCachingForOptionsFlag)
	})

	t.Run("wrong caching for post boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "incorrect",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "false",
			"--slice-l2-caching-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorCachingForPostFlag)
	})

	t.Run("wrong caching string sort boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "false",
			"--enable-caching-for-post", "false",
			"--enable-caching-string-sort", "precise",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "false",
			"--slice-l2-caching-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorCachingStringSortFlag)
	})

	t.Run("wrong slice configuration enabled boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "faithful",
			"--slice-l2-caching-enabled", "false",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorSliceConfigurationFlag)
	})

	t.Run("wrong slice l2 caching enabled boolean var", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"-a", "1673635841",
			"--name", "fmaiswaybetter",
			"--adaptive-delivery-action", "ignore",
			"--browser-cache-settings", "override",
			"--cdn-cache-settings", "honor",
			"--cnd-cache-settings-maximum-ttl", "60",
			"--browser-cache-settings-maximum-ttl", "60",
			"--cache-by-cookies", "whitelist",
			"--cookie-names", "aa,123,987",
			"--cache-by-query-string", "whitelist",
			"--query-string-fields", "heyy,yoo",
			"--enable-caching-for-options", "true",
			"--enable-caching-for-post", "true",
			"--enable-caching-string-sort", "true",
			"--l2-caching-enabled", "true",
			"--slice-configuration-enabled", "true",
			"--slice-l2-caching-enabled", "erroneous",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, cache_settings.ErrorSliceL2CachingFlag)
	})

	t.Run("error unmarshall", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_applications/1673635841"),
			httpmock.JSONFromFile("./fixtures/app_result.json"),
		)

		mock.Register(
			httpmock.REST("POST", "edge_applications/1673635841/cache_settings"),
			httpmock.JSONFromFile("./fixtures/result.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1673635841",
			"--in", "./fixtures/error",
		})

		err := cmd.Execute()
		require.ErrorIs(t, err, utils.ErrorUnmarshalReader)
	})

}
