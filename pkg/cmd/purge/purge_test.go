package purge

import (
	"testing"

	msg "github.com/aziontech/azion-cli/messages/purge"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("purge urls", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "purge/url"),
			httpmock.StatusStringResponse(201, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		err := purgeUrls([]string{"www.example.com", "www.httpin.com"}, f)

		require.NoError(t, err)
	})

	t.Run("purge wildcard", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "purge/wildcard"),
			httpmock.StatusStringResponse(201, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		err := purgeWildcard([]string{"www.example.com/*"}, f)

		require.NoError(t, err)
	})

	t.Run("purge wildcard - more than one item", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "purge/wildcard"),
			httpmock.StatusStringResponse(201, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		err := purgeWildcard([]string{"www.example.com/*", "www.pudim.com/*"}, f)

		require.ErrorContains(t, err, msg.ErrorTooManyUrls.Error())
	})

	t.Run("purge cache-keys", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("POST", "purge/cachekey"),
			httpmock.StatusStringResponse(201, ""),
		)

		f, _, _ := testutils.NewFactory(mock)

		err := purgeCacheKeys([]string{"www.example.com/@@cookie_name=cookie_value"}, f)

		require.NoError(t, err)
	})

}
