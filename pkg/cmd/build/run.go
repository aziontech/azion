package build

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (cmd *BuildCmd) run() error {
	logger.Debug("Running build subcommand from edge_applications command tree")
	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}

	err = RunBuildCmdLine(cmd, path)
	if err != nil {
		return err
	}

	return nil
}

func RunBuildCmdLine(cmd *BuildCmd, path string) error {
	var err error

	conf, err := cmd.GetAzionJsonContent()
	if err != nil {
		return err
	}

	if conf.Template == "simple" {
		logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildSimple)
		return nil
	}

	if conf.Template == "static" {
		logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildSimple)
		return nil
	}

	err = checkArgsJson(cmd)
	if err != nil {
		return err
	}

	if conf.Template != "nextjs" {
		return vulcan(cmd, conf, path)
	}

	if conf.Template == "nextjs" {
		return adapter(cmd, conf)
	}

	return utils.ErrorUnsupportedType
}

func checkArgsJson(cmd *BuildCmd) error {
	workDirPath, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return utils.ErrorInternalServerError
	}

	workDirPath += "/azion/args.json"
	_, err = cmd.FileReader(workDirPath)
	if err != nil {
		if err := cmd.WriteFile(workDirPath, []byte("{}"), 0644); err != nil {
			logger.Debug("Error while trying to create args.json file", zap.Error(err))
			return fmt.Errorf(utils.ErrorCreateFile.Error(), workDirPath)
		}
	}

	return nil
}
