package vulcan

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"go.uber.org/zap"
)

var (
	currentMajor       = 1
	installBundler     = "npx --yes %s @aziontech/bundler%s %s"
	firstTimeExecuting = "@1.2.0"
	versionVulcan      = "@1.2.0"
	releaseChannel     = ""
	StagePkgURL        = "https://pkg.pr.new/aziontech/bundler/@aziontech/bundler@main"
)

var ErrConfigFileNotFound = errors.New("The azion.config file informed in --config-file was not found. Verify the path and try again")

const (
	ConfigVersionV3 = 3
	ConfigVersionV4 = 4
)

// configVersionCommands are the bundler subcommands that accept --config-version.
var configVersionCommands = map[string]bool{
	"build":    true,
	"dev":      true,
	"presets":  true,
	"manifest": true,
}

type VulcanPkg struct {
	Command          func(flags, params string, f *cmdutil.Factory) string
	CheckVulcanMajor func(currentVersion string, f *cmdutil.Factory, vulcan *VulcanPkg) error
	ReadSettings     func(string) (token.Settings, error)
}

func NewVulcan() *VulcanPkg {
	return newVulcan(ConfigVersionV4)
}

// NewVulcanV3 returns a bundler that targets the Azion API v3 config
func NewVulcanV3() *VulcanPkg {
	return newVulcan(ConfigVersionV3)
}

func newVulcan(configVersion int) *VulcanPkg {
	return &VulcanPkg{
		Command: func(flags, params string, f *cmdutil.Factory) string {
			return command(flags, withConfigVersion(params, configVersion), f)
		},
		CheckVulcanMajor: checkVulcanMajor,
		ReadSettings:     token.ReadSettings,
	}
}

// withConfigVersion appends --config-version to the subcommands that support it
func withConfigVersion(params string, configVersion int) string {
	// Callers may pass format verbs glued to the subcommand (e.g. "build%s")
	subcommand := strings.TrimLeft(params, " ")
	if i := strings.IndexFunc(subcommand, func(r rune) bool { return r < 'a' || r > 'z' }); i >= 0 {
		subcommand = subcommand[:i]
	}
	if !configVersionCommands[subcommand] {
		return params
	}
	return fmt.Sprintf("%s --config-version %d", params, configVersion)
}

func command(flags, params string, f *cmdutil.Factory) string {
	if releaseChannel == "stage" {
		stageCmd := "npx --yes %s " + StagePkgURL + " %s"
		if f.Logger.Debug {
			stageCmd = "DEBUG=true " + stageCmd
		}
		return fmt.Sprintf(stageCmd, flags, params)
	}

	selectedVersion := versionVulcan
	if f.Logger.Debug {
		installDebug := "DEBUG=true " + installBundler
		return fmt.Sprintf(installDebug, flags, selectedVersion, params)
	}
	return fmt.Sprintf(installBundler, flags, selectedVersion, params)
}

func checkVulcanMajor(currentVersion string, f *cmdutil.Factory, vulcan *VulcanPkg) error {
	activeProfile := f.GetActiveProfile()
	parts := strings.Split(currentVersion, ".")
	// strings.Split will always return at least one element, so parts will always be len>0
	// to avoid this, I am checking if version is empty. If so, I just use an empty slice
	if currentVersion == "" {
		parts = []string{}
	}

	// Extract the first part and convert it to a number
	if len(parts) > 0 {
		firstNumber, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		config, err := vulcan.ReadSettings(activeProfile)
		if err != nil {
			return err
		}

		if firstNumber > currentMajor {
			return nil
		}

		config.LastVulcanVersion = currentVersion
		err = token.WriteSettings(config, activeProfile)
		if err != nil {
			logger.Debug("Error while saving settings", zap.Error(err))
			return err
		}

	} else {
		logger.Debug("Failed to parse information on current Bundler version")
		versionVulcan = firstTimeExecuting
		return nil
	}
	return nil
}

// ConfigFileArg returns the --config-file option for the bundler, or an empty
// string when no config file was given. The path is single-quoted for sh -c.
func ConfigFileArg(configFile string) string {
	if configFile == "" {
		return ""
	}
	return fmt.Sprintf(" --config-file '%s'", strings.ReplaceAll(configFile, "'", `'\''`))
}

// CheckConfigFile fails early when the given azion.config file does not exist,
// before any command changes local or remote state. An empty path is valid.
func CheckConfigFile(configFile string) error {
	if configFile == "" {
		return nil
	}
	info, err := os.Stat(configFile)
	if err != nil || !info.Mode().IsRegular() {
		return fmt.Errorf("%w: %s", ErrConfigFileNotFound, configFile)
	}
	logger.Debug(fmt.Sprintf("Using %s config file", configFile))
	return nil
}
