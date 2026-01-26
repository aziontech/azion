package firewallrules

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestListFirewallRules(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("list firewall rules", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		listCmd.ListRules = func(_ context.Context, _ *contracts.ListOptions, _ int64) (*sdk.PaginatedFirewallRuleList, error) {
			b, err := os.ReadFile("./fixtures/response.json")
			if err != nil {
				return nil, err
			}
			var out sdk.PaginatedFirewallRuleList
			if err := json.Unmarshal(b, &out); err != nil {
				return nil, err
			}
			return &out, nil
		}

		cobraCmd := NewCobraCmd(listCmd, f)
		cobraCmd.SetArgs([]string{"--firewall-id", "1234"})

		_, err := cobraCmd.ExecuteC()
		require.NoError(t, err)
	})
}
