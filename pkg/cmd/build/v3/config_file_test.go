package build

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	msgVulcan "github.com/aziontech/azion-cli/messages/vulcan"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
)

func TestBuildConfigFile(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	configFile := filepath.Join(t.TempDir(), "azion.staging.config.ts")
	if err := os.WriteFile(configFile, []byte("export default {}"), 0o644); err != nil {
		t.Fatal(err)
	}

	newBuild := func(captured *string) *BuildCmd {
		f, _, _ := testutils.NewFactory(&httpmock.Registry{})
		b := NewBuildCmd(f)
		b.GetAzionJsonContent = func(string) (*contracts.AzionApplicationOptionsV3, error) {
			return &contracts.AzionApplicationOptionsV3{}, nil
		}
		b.WriteAzionJsonContent = func(*contracts.AzionApplicationOptionsV3, string) error { return nil }
		b.CommandRunner = func(*cmdutil.Factory, string, []string) (string, error) { return "", nil }
		b.CommandRunInteractive = func(_ *cmdutil.Factory, comm string) error {
			*captured = comm
			return nil
		}
		return b
	}

	t.Run("forwards --config-file to the bundler", func(t *testing.T) {
		var captured string
		err := newBuild(&captured).run(&contracts.BuildInfoV3{Preset: "javascript", ConfigFile: configFile}, &[]string{})
		if err != nil {
			t.Fatalf("run() error = %v", err)
		}
		if !strings.Contains(captured, "build --preset javascript --config-file '"+configFile+"'") {
			t.Errorf("command = %q, want it to forward --config-file", captured)
		}
	})

	t.Run("omits --config-file when not sent", func(t *testing.T) {
		var captured string
		err := newBuild(&captured).run(&contracts.BuildInfoV3{Preset: "javascript"}, &[]string{})
		if err != nil {
			t.Fatalf("run() error = %v", err)
		}
		if strings.Contains(captured, "--config-file") {
			t.Errorf("command = %q, want no --config-file", captured)
		}
	})

	t.Run("fails before running the bundler when the file is missing", func(t *testing.T) {
		var captured string
		err := newBuild(&captured).run(&contracts.BuildInfoV3{ConfigFile: "missing.config.ts"}, &[]string{})
		if !errors.Is(err, msgVulcan.ErrConfigFileNotFound) {
			t.Fatalf("run() error = %v, want ErrConfigFileNotFound", err)
		}
		if captured != "" {
			t.Errorf("bundler ran with %q, want it not to run", captured)
		}
	})
}
