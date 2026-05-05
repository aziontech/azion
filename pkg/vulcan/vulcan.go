package vulcan

import (
	"fmt"
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
	installVulcanV3    = "npx --yes %s edge-functions%s %s"
	firstTimeExecuting = "@1.0.0"
	versionVulcan      = "@1.0.0"
	releaseChannel     = ""
	StagePkgURL        = "https://pkg.pr.new/aziontech/bundler/@aziontech/bundler@main"
)

type VulcanPkg struct {
	Command          func(flags, params string, f *cmdutil.Factory) string
	CheckVulcanMajor func(currentVersion string, f *cmdutil.Factory, vulcan *VulcanPkg) error
	ReadSettings     func(string) (token.Settings, error)
}

func NewVulcan() *VulcanPkg {
	return &VulcanPkg{
		Command:          command,
		CheckVulcanMajor: checkVulcanMajor,
		ReadSettings:     token.ReadSettings,
	}
}

func NewVulcanV3() *VulcanPkg {
	versionVulcan = "@5.3.1"
	currentMajor = 5
	firstTimeExecuting = "@5.3.1"
	return &VulcanPkg{
		Command:          commandV3,
		CheckVulcanMajor: checkVulcanMajor,
		ReadSettings:     token.ReadSettings,
	}
}

func command(flags, params string, f *cmdutil.Factory) string {
	if releaseChannel == "stage" {
		stageCmd := "npx --yes %s " + StagePkgURL + " %s"
		if f.Logger.Debug {
			stageCmd = "DEBUG=true " + stageCmd
		}
		return fmt.Sprintf(stageCmd, flags, params)
	}

	// Production builds use @aziontech/bundler for V4
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
