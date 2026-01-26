package firewallrules

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDescribeFirewallRule(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("describe firewall rule", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		describeCmd := NewDescribeCmd(f)
		describeCmd.GetFirewallRule = func(_ context.Context, _, _ int64) (sdk.FirewallRule, error) {
			b, err := os.ReadFile("./fixtures/response.json")
			if err != nil {
				return sdk.FirewallRule{}, err
			}
			var wrapper struct {
				State string          `json:"state"`
				Data  sdk.FirewallRule `json:"data"`
			}
			if err := json.Unmarshal(b, &wrapper); err != nil {
				return sdk.FirewallRule{}, err
			}
			return wrapper.Data, nil
		}

		cobraCmd := NewCobraCmd(describeCmd, f)
		cobraCmd.SetArgs([]string{"--firewall-id", "1234", "--rule-id", "111"})

		_, err := cobraCmd.ExecuteC()
		require.NoError(t, err)
	})
}
