package build

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

func runCommand(cmd *BuildCmd, command string) error {
	logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildStart)

	logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildRunningCmd)
	logger.FInfo(cmd.Io.Out, fmt.Sprintf("$ %s\n", command))

	err := cmd.CommandRunnerStream(cmd.Io.Out, command, []string{})
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrFailedToRunBuildCommand
	}

	logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildSuccessful)
	return nil
}

func vulcan(cmd *BuildCmd, typeLang, mode string) error {
	const command string = "npx --yes edge-functions@1.0.0 build --preset %s --mode %s"

	err := runCommand(cmd, fmt.Sprintf(command, typeLang, mode))
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	return nil
}

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
