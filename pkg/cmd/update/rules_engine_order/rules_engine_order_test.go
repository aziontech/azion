package rules_engine_order

import (
	"fmt"
	"net/http"
	"testing"

	msg "github.com/aziontech/azion-cli/messages/update/rules_engine_order"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

const orderedListResponse = `{"count":0,"results":[]}`

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	t.Run("order request phase rules", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPut, "workspace/applications/1234/request_rules/order"),
			httpmock.JSONFromString(orderedListResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1234",
			"--phase", "request",
			"--rule-ids", "10,20,30",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 1234), stdout.String())
	})

	t.Run("order response phase rules", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPut, "workspace/applications/1234/response_rules/order"),
			httpmock.JSONFromString(orderedListResponse),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1234",
			"--phase", "response",
			"--rule-ids", "30,20,10",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.OutputSuccess, 1234), stdout.String())
	})

	t.Run("invalid phase", func(t *testing.T) {
		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1234",
			"--phase", "invalid",
			"--rule-ids", "10,20,30",
		})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorInvalidPhase)
	})

	t.Run("invalid rule ids", func(t *testing.T) {
		mock := &httpmock.Registry{}

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1234",
			"--phase", "request",
			"--rule-ids", "10,abc,30",
		})

		err := cmd.Execute()

		require.ErrorIs(t, err, msg.ErrorConvertRuleIDs)
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPut, "workspace/applications/1234/request_rules/order"),
			httpmock.StatusStringResponse(http.StatusBadRequest, "Invalid"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1234",
			"--phase", "request",
			"--rule-ids", "10,20,30",
		})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("internal server error", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST(http.MethodPut, "workspace/applications/1234/response_rules/order"),
			httpmock.StatusStringResponse(http.StatusInternalServerError, "Internal Server Error"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)
		cmd.SetArgs([]string{
			"--application-id", "1234",
			"--phase", "response",
			"--rule-ids", "10,20,30",
		})

		err := cmd.Execute()

		require.Error(t, err)
	})
}
