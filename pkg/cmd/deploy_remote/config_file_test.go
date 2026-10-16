package deploy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"go.uber.org/zap/zapcore"
)

// Every bundler call made by the local deploy must target the --config-file,
// otherwise the build and the azion.config rewrites would use different files.
func TestLocalDeployForwardsConfigFile(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	configFile := filepath.Join(t.TempDir(), "azion.staging.config.ts")
	if err := os.WriteFile(configFile, []byte("export default {}"), 0o644); err != nil {
		t.Fatal(err)
	}

	newState := func(commands *[]string) *DeployState {
		f, _, _ := testutils.NewFactory(&httpmock.Registry{})
		record := func(_ *cmdutil.Factory, comm string) error {
			*commands = append(*commands, comm)
			return nil
		}
		cmd := NewDeployCmd(f)
		cmd.ConfigFile = configFile
		cmd.commandRunInteractive = record
		cmd.commandRunnerOutput = func(*cmdutil.Factory, string, []string) (string, error) { return "", nil }
		cmd.BuildCmd = func(f *cmdutil.Factory) *build.BuildCmd {
			b := build.NewBuildCmd(f)
			b.GetAzionJsonContent = func(string) (*contracts.AzionApplicationOptions, error) {
				return &contracts.AzionApplicationOptions{Preset: "javascript"}, nil
			}
			b.WriteAzionJsonContent = func(*contracts.AzionApplicationOptions, string) error { return nil }
			b.CommandRunner = func(*cmdutil.Factory, string, []string) (string, error) { return "", nil }
			b.CommandRunInteractive = record
			return b
		}
		return &DeployState{
			Cmd: cmd, F: f, Msgs: &[]string{},
			Conf:      &contracts.AzionApplicationOptions{Name: "app", Preset: "javascript"},
			OldPrefix: "old", NewPrefix: "new",
		}
	}

	assertAll := func(t *testing.T, commands []string, want int) {
		t.Helper()
		if len(commands) != want {
			t.Fatalf("got %d bundler commands, want %d: %q", len(commands), want, commands)
		}
		for _, c := range commands {
			if !strings.Contains(c, "--config-file '"+configFile+"'") {
				t.Errorf("command %q does not forward --config-file", c)
			}
		}
	}

	t.Run("build step: prefix rotation and build", func(t *testing.T) {
		var commands []string
		if err := newState(&commands).build(); err != nil {
			t.Fatalf("build() error = %v", err)
		}
		assertAll(t, commands, 2)
	})

	t.Run("first run: name replacement and build", func(t *testing.T) {
		var commands []string
		if err := newState(&commands).bundlerInit(); err != nil {
			t.Fatalf("bundlerInit() error = %v", err)
		}
		assertAll(t, commands, 2)
	})
}

func TestGenerateManifestCommand(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	var commands []string
	f, _, _ := testutils.NewFactory(&httpmock.Registry{})
	cmd := NewDeployCmd(f)
	cmd.ConfigFile = "azion.staging.config.ts"
	cmd.AliasEnv = true
	cmd.commandRunInteractive = func(_ *cmdutil.Factory, comm string) error {
		commands = append(commands, comm)
		return nil
	}
	s := &DeployState{Cmd: cmd, F: f, Msgs: &[]string{}, Conf: &contracts.AzionApplicationOptions{Preset: "nuxt"}}

	if err := s.generateManifest(); err != nil {
		t.Fatalf("generateManifest() error = %v", err)
	}
	want := "build --preset 'nuxt' --only-generate-config --alias-env --config-file 'azion.staging.config.ts' --config-version 4"
	if len(commands) != 1 || !strings.HasSuffix(commands[0], want) {
		t.Fatalf("commands = %q, want one ending in %q", commands, want)
	}
}
