package vulcan

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/vulcan"
	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"go.uber.org/zap"
)

var (
	currentMajor       = 1
	installBundler     = "npx --yes %s @aziontech/bundler%s %s"
	installVulcanV3    = "npx --yes %s edge-functions%s %s"
	firstTimeExecuting = "@1.3.0"
	versionVulcan      = "@1.3.0"
	releaseChannel     = ""
	StagePkgURL        = "https://pkg.pr.new/aziontech/bundler/@aziontech/bundler@main"
)

// The --config-version sent to the bundler for each API generation. The value
// is the generation's number: the bundler writes and reads the azion.config of
// that generation.
const (
	ConfigVersionV3 = 3
	ConfigVersionV4 = 4
	ConfigVersionV6 = 6
)

// ConfigVersion is the --config-version of an API generation. An unknown or
// unresolved generation gets v4's, the CLI's default generation.
func ConfigVersion(v apiversion.Version) int {
	switch v {
	case apiversion.V3:
		return ConfigVersionV3
	case apiversion.V6:
		return ConfigVersionV6
	default:
		return ConfigVersionV4
	}
}

// For returns the bundler for the account's API generation, with the matching
// --config-version: v3 runs edge-functions, and v4 and v6 the published
// bundler. Every command that runs the bundler gets it from here, so
// the config version always follows the account.
func For(f *cmdutil.Factory) *VulcanPkg {
	var v apiversion.Version
	if f != nil {
		v = f.APIVersion
	}
	switch v {
	case apiversion.V3:
		return NewVulcanV3()
	case apiversion.V6:
		return NewVulcanV6()
	default:
		return NewVulcan()
	}
}

// v3SendsConfigVersion is false while the v3 bundler (edge-functions@5.3.2)
// has no --config-version option.
const v3SendsConfigVersion = false

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

// NewVulcanV6 runs the bundler for an API v6 project: the same published
// bundler as v4 (versionVulcan), with --config-version 6.
//
// Setting AZIONCLI_BUNDLER_V6 to a package spec (for example a pkg.pr.new
// build) runs that instead, to try a bundler change before it is published.
func NewVulcanV6() *VulcanPkg {
	override := os.Getenv("AZIONCLI_BUNDLER_V6")
	if override == "" {
		return newVulcan(ConfigVersionV6)
	}
	return &VulcanPkg{
		Command: func(flags, params string, f *cmdutil.Factory) string {
			cmd := fmt.Sprintf("npx --yes %s %s %s", flags, override, withConfigVersion(params, ConfigVersionV6))
			if f.Logger.Debug {
				cmd = "DEBUG=true " + cmd
			}
			return cmd
		},
		CheckVulcanMajor: checkVulcanMajor,
		ReadSettings:     token.ReadSettings,
	}
}

// NewVulcanV3 runs edge-functions for API v3 projects. Its config version is
// ConfigVersionV3, but edge-functions@5.3.2 rejects --config-version ("unknown
// option"), so the flag is not sent; v3SendsConfigVersion turns it on once a
// v3 package accepts it.
func NewVulcanV3() *VulcanPkg {
	versionVulcan = "@5.3.2"
	currentMajor = 5
	firstTimeExecuting = "@5.3.2"
	return &VulcanPkg{
		Command: func(flags, params string, f *cmdutil.Factory) string {
			if v3SendsConfigVersion {
				params = withConfigVersion(params, ConfigVersionV3)
			}
			return commandV3(flags, params, f)
		},
		CheckVulcanMajor: checkVulcanMajor,
		ReadSettings:     token.ReadSettings,
	}
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

// commandV3 uses edge-functions for V3 compatibility
func commandV3(flags, params string, f *cmdutil.Factory) string {
	if releaseChannel == "stage" {
		stageCmd := "npx --yes %s " + StagePkgURL + " %s"
		if f.Logger.Debug {
			stageCmd = "DEBUG=true " + stageCmd
		}
		return fmt.Sprintf(stageCmd, flags, params)
	}

	// V3 uses edge-functions package
	selectedVersion := versionVulcan
	if f.Logger.Debug {
		installDebug := "DEBUG=true " + installVulcanV3
		return fmt.Sprintf(installDebug, flags, selectedVersion, params)
	}
	return fmt.Sprintf(installVulcanV3, flags, selectedVersion, params)
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
		return fmt.Errorf("%w: %s", msg.ErrConfigFileNotFound, configFile)
	}
	logger.Debug(fmt.Sprintf("Using %s config file", configFile))
	return nil
}
