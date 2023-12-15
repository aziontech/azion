package build

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (cmd *BuildCmd) run() error {
	logger.Debug("Running build command")

	err := RunBuildCmdLine(cmd)
	if err != nil {
		return err
	}

	return nil
}

func RunBuildCmdLine(cmd *BuildCmd) error {
	var err error

	conf, err := cmd.GetAzionJsonContent()
	if err != nil {
		logger.Debug("Error while building your project", zap.Error(err))
		return msg.ErrorBuilding
	}

	if Preset != "" {
		conf.Template = Preset
	}

	if Mode != "" {
		conf.Mode = Mode
	}

	err = checkArgsJson(cmd, conf)
	if err != nil {
		return err
	}

	if conf.Template == "simple" {
		logger.FInfo(cmd.Io.Out, msg.BuildSimple)
		return nil
	}

	if conf.Template == "static" {
		versionID := cmd.VersionID()
		conf.Prefix = versionID

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			return nil
		}
		logger.FInfo(cmd.Io.Out, msg.BuildSimple)
		return nil
	}

	if conf.Template != "nextjs" {
		return vulcan(cmd, conf)
	}

	if conf.Template == "nextjs" {
		return adapter(cmd, conf)
	}

	return utils.ErrorUnsupportedType
}

func checkArgsJson(cmd *BuildCmd, conf *contracts.AzionApplicationOptions) error {

	workDirPath := conf.ProjectRoot + "/azion/args.json"
	_, err := cmd.FileReader(workDirPath)
	if err != nil {
		if err := cmd.WriteFile(workDirPath, []byte("{}"), 0644); err != nil {
			logger.Debug("Error while trying to create args.json file", zap.Error(err))
			return fmt.Errorf(utils.ErrorCreateFile.Error(), workDirPath)
		}
	}

	return nil
}
