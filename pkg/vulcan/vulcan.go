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
	currentMajor         = 2
	installEdgeFunctions = "npx --yes %s edge-functions%s %s"
	firstTimeExecuting   = "@v2.5.0"
)

var versionVulcan = "@2.6.0-stage.17"

func Command(flags, params string) string {
	return fmt.Sprintf(installEdgeFunctions, flags, versionVulcan, params)
}

func CheckVulcanMajor(currentVersion string, f *cmdutil.Factory) error {
	parts := strings.Split(currentVersion, ".")

	// Extract the first part and convert it to a number
	if len(parts) > 0 {
		firstNumber, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}

		config, err := token.ReadSettings()
		if err != nil {
			return err
		}

		if firstNumber > currentMajor {
			logger.FInfo(f.IOStreams.Out, msg.NewMajorVersion)

			if config.LastVulcanVersion == "" {
				// new version update please
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

	} else {
		logger.Debug("Failed to parse information on current vulcan version")
		versionVulcan = firstTimeExecuting
		return nil
	}
	return nil
}
