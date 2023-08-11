package build

import (
	"fmt"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

func adapter(cmd *BuildCmd, path, pathAzionJson string, file []byte) error {
	const command = "npx --yes azion-framework-adapter@0.4.0 build --version-id %s"

	// pre-build version id. Used to check if there were changes to the project
	versionID := cmd.VersionID(path)

	err := runCommand(cmd, fmt.Sprintf(command, versionID))
	if err != nil {
		return err
	}

	strJson, err := sjson.Set(string(file), "version-id", versionID)
	if err != nil {
		logger.Debug("Error while writing version-id to azion.json file: ", zap.Error(err))
		return utils.ErrorWritingAzionJsonFile
	}

	err = cmd.WriteFile(pathAzionJson, []byte(strJson), 0644)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return utils.ErrorWritingAzionJsonFile
	}

	return err
}
