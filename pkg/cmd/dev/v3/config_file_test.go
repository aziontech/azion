package dev

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	msgVulcan "github.com/aziontech/azion-cli/messages/vulcan"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
)

func TestDevConfigFile(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	configFile := filepath.Join(t.TempDir(), "azion.staging.config.ts")
	if err := os.WriteFile(configFile, []byte("export default {}"), 0o644); err != nil {
		t.Fatal(err)
	}

	newDev := func(captured *string) *DevCmd {
		f, _, _ := testutils.NewFactory(&httpmock.Registry{})
		d := NewDevCmd(f)
		d.CommandRunInteractive = func(_ *cmdutil.Factory, comm string) error {
			*captured = comm
			return nil
		}
		return d
	}

	t.Run("forwards --config-file to the bundler", func(t *testing.T) {
		var captured string
		d := newDev(&captured)
		d.ConfigFile = configFile
		if err := d.Run(d.F); err != nil {
			t.Fatalf("Run() error = %v", err)
		}
		if !strings.Contains(captured, " dev --config-version ") || !strings.HasSuffix(captured, "--config-file '"+configFile+"'") {
			t.Errorf("command = %q, want it to forward --config-file", captured)
		}
	})

	t.Run("fails before running the bundler when the file is missing", func(t *testing.T) {
		var captured string
		d := newDev(&captured)
		d.ConfigFile = "missing.config.ts"
		err := d.Run(d.F)
		if !errors.Is(err, msgVulcan.ErrConfigFileNotFound) {
			t.Fatalf("Run() error = %v, want ErrConfigFileNotFound", err)
		}
		if captured != "" {
			t.Errorf("bundler ran with %q, want it not to run", captured)
		}
	})
}
