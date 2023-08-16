package build

import (
	"fmt"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func adapter(cmd *BuildCmd, conf *contracts.AzionApplicationOptions) error {
	const command = "npx --yes azion-framework-adapter@0.4.0 build --version-id %s"

	// pre-build version id. Used to check if there were changes to the project
	versionID := cmd.VersionID()

	err := runCommand(cmd, fmt.Sprintf(command, versionID))
	if err != nil {
		return err
	}

	conf.VersionID = versionID

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return utils.ErrorWritingAzionJsonFile
	}

	return err
}
