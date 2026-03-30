package waf

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap/zapcore"

	msg "github.com/aziontech/azion-cli/messages/update/waf"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	t.Run("update WAF with name and active", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--waf-id", "1337", "--name", "Updated WAF", "--active", "true"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, int64(666)), stdout.String())
	})

	t.Run("update WAF engine settings", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--waf-id", "1337",
			"--engine-version", "2021-Q3",
			"--type", "score",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, int64(666)), stdout.String())
	})

	t.Run("update WAF with rulesets", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--waf-id", "1337",
			"--rulesets", "1,2,3",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, int64(666)), stdout.String())
	})

	t.Run("update WAF with thresholds", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{
			"--waf-id", "1337",
			"--thresholds", "cross_site_scripting=highest,sql_injection=high",
		})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, int64(666)), stdout.String())
	})

	t.Run("bad request", func(t *testing.T) {
		mock := &httpmock.Registry{}
		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1234"),
			httpmock.StatusStringResponse(http.StatusBadRequest, `{"details": "invalid field active"}`),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--waf-id", "1234", "--active", "invalid"})

		err := cmd.Execute()

		require.Error(t, err)
	})

	t.Run("update with file", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, stdout, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--waf-id", "1337", "--file", "./fixtures/update.json"})

		err := cmd.Execute()

		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(msg.UpdateOutputSuccess, int64(666)), stdout.String())
	})

	t.Run("Error Field active not is boolean", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--waf-id", "1337", "--active", "not_a_bool"})

		err := cmd.Execute()
		stringErr := fmt.Sprintf("%s: %s", msg.ErrorActiveFlag, "not_a_bool")
		if stringErr == err.Error() {
			return
		}
		t.Fatalf("Error: %q", err)
	})

	t.Run("Error Invalid rulesets format", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--waf-id", "1337", "--rulesets", "not_an_int"})

		err := cmd.Execute()
		stringErr := fmt.Sprintf("%s: %s", msg.ErrorRulesetsFlag, "not_an_int")
		if stringErr == err.Error() {
			return
		}
		t.Fatalf("Error: %q", err)
	})

	t.Run("Error Invalid thresholds format", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--waf-id", "1337", "--thresholds", "invalid_format"})

		err := cmd.Execute()
		if err == msg.ErrorThresholdsFlag {
			return
		}
		t.Fatalf("Error: %q", err)
	})

	t.Run("Error Invalid threat value in thresholds", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--waf-id", "1337", "--thresholds", "invalid_threat=highest"})

		err := cmd.Execute()
		if err == msg.ErrorThresholdThreat {
			return
		}
		t.Fatalf("Error: %q", err)
	})

	t.Run("Error Invalid sensitivity value in thresholds", func(t *testing.T) {
		mock := &httpmock.Registry{}

		mock.Register(
			httpmock.REST("PATCH", "workspace/wafs/1337"),
			httpmock.JSONFromFile("./fixtures/response.json"),
		)

		f, _, _ := testutils.NewFactory(mock)

		cmd := NewCmd(f)

		cmd.SetArgs([]string{"--waf-id", "1337", "--thresholds", "sql_injection=invalid_sensitivity"})

		err := cmd.Execute()
		if err == msg.ErrorThresholdSensitivity {
			return
		}
		t.Fatalf("Error: %q", err)
	})
}
