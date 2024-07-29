package vulcan

import (
	"fmt"
	"strconv"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/vulcan"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/pelletier/go-toml"
	"go.uber.org/zap"
)

const (
	currentMajor         = 3
	installEdgeFunctions = "npx --yes %s edge-functions%s %s"
	firstTimeExecuting   = "@3.0.0"
)

var versionVulcan = "@latest"

type VulcanPkg struct {
	Command          func(flags, params string, f *cmdutil.Factory) string
	CheckVulcanMajor func(currentVersion string, f *cmdutil.Factory, vulcan *VulcanPkg) error
	ReadSettings     func() (token.Settings, error)
}

func NewVulcan() *VulcanPkg {
	return &VulcanPkg{
		Command:          command,
		CheckVulcanMajor: checkVulcanMajor,
		ReadSettings:     token.ReadSettings,
	}
}

func command(flags, params string, f *cmdutil.Factory) string {
	if f.Logger.Debug {
		installDebug := "DEBUG=true " + installEdgeFunctions
		return fmt.Sprintf(installDebug, flags, versionVulcan, params)
	}
	return fmt.Sprintf(installEdgeFunctions, flags, versionVulcan, params)
}

func checkVulcanMajor(currentVersion string, f *cmdutil.Factory, vulcan *VulcanPkg) error {
	parts := strings.Split(currentVersion, ".")
	if currentVersion == "" {
		parts = []string{}
	}

	if len(parts) > 0 {
		firstNumber, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		config, err := vulcan.ReadSettings()
		if err != nil {
			return err
		}

		if firstNumber > currentMajor {
			logger.FInfo(f.IOStreams.Out, msg.NewMajorVersion)

			if config.LastVulcanVersion == "" {
				versionVulcan = firstTimeExecuting
				return nil
			}

			versionVulcan = "@" + strings.TrimRight(config.LastVulcanVersion, "\n")
			return nil
		}

		config.LastVulcanVersion = currentVersion
		client, err := token.New(&token.Config{Client: f.HttpClient})
		if err != nil {
			return err
		}

		byteSettings, err := toml.Marshal(config)
		if err != nil {
			logger.Debug("Error while marshalling settings.toml", zap.Error(err))
			return err
		}

		_, err = client.Save(byteSettings)
		if err != nil {
			logger.Debug("Error while saving settings", zap.Error(err))
			return err
		}

		versionVulcan = "@" + currentVersion

	} else {
		logger.Debug("Failed to parse information on current vulcan version")
		versionVulcan = firstTimeExecuting
		return nil
	}
	return nil
}
