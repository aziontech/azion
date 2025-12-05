package ruleengine

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/stretchr/testify/require"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("list all rules engines", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)

		listCmd := NewListCmd(f)
		// Override the API call to deserialize our fixture into the SDK struct
		listCmd.ListRulesEngineRequest = func(_ context.Context, _ *contracts.ListOptions, _ int64) (*sdk.PaginatedApplicationRequestPhaseRuleEngineList, error) {
			b, err := os.ReadFile("./fixtures/rules.json")
			if err != nil {
				return nil, err
			}
			var out sdk.PaginatedApplicationRequestPhaseRuleEngineList
			if err := json.Unmarshal(b, &out); err != nil {
				return nil, err
			}
			return &out, nil
		}
		cobraCmd := NewCobraCmd(listCmd, f)
		cobraCmd.SetArgs([]string{"--application-id", "1678743802", "--phase", "request"})

		_, err := cobraCmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("list - page 0 should generate an error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("GET", "edge_application/applications/1678743802/rules"),
			httpmock.StatusStringResponse(404, "Not Found"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--application-id", "1673635839", "--page", "0"})

		_, err := cmd.ExecuteC()
		require.Error(t, err)
	})

	t.Run("list all rules engines - ask for input", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)
		listCmd := NewListCmd(f)
		// Stub API to use fixture
		listCmd.ListRulesEngineRequest = func(_ context.Context, _ *contracts.ListOptions, _ int64) (*sdk.PaginatedApplicationRequestPhaseRuleEngineList, error) {
			b, err := os.ReadFile("./fixtures/rules.json")
			if err != nil {
				return nil, err
			}
			var out sdk.PaginatedApplicationRequestPhaseRuleEngineList
			if err := json.Unmarshal(b, &out); err != nil {
				return nil, err
			}
			return &out, nil
		}
		listCmd.AskInput = func(s string) (string, error) {
			return "1678743802", nil
		}
		cmd := NewCobraCmd(listCmd, f)

		cmd.SetArgs([]string{"--phase", "request"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})

	t.Run("no itens", func(t *testing.T) {
		mock := &httpmock.Registry{}
		f, _, _ := testutils.NewFactory(mock)
		listCmd := NewListCmd(f)
		// Stub API to use empty fixture
		listCmd.ListRulesEngineRequest = func(_ context.Context, _ *contracts.ListOptions, _ int64) (*sdk.PaginatedApplicationRequestPhaseRuleEngineList, error) {
			b, err := os.ReadFile("./fixtures/norules.json")
			if err != nil {
				return nil, err
			}
			var out sdk.PaginatedApplicationRequestPhaseRuleEngineList
			if err := json.Unmarshal(b, &out); err != nil {
				return nil, err
			}
			return &out, nil
		}
		cmd := NewCobraCmd(listCmd, f)
		cmd.SetArgs([]string{"--application-id", "1678743802", "--phase", "request"})

		_, err := cmd.ExecuteC()
		require.NoError(t, err)
	})
}
