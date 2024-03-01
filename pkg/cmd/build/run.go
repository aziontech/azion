package build

import (
	"fmt"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (cmd *BuildCmd) run(fields *contracts.BuildInfo) error {
	logger.Debug("Running build command")

	err := RunBuildCmdLine(cmd, fields)
	if err != nil {
		return err
	}

	return nil
}

func RunBuildCmdLine(cmd *BuildCmd, fields *contracts.BuildInfo) error {
	var err error

	conf, err := cmd.GetAzionJsonContent()
	if err != nil {
		logger.Debug("Error while building your project", zap.Error(err))
		return msg.ErrorBuilding
	}

	if fields.Preset != "" {
		conf.Template = fields.Preset
	}

	if fields.Mode != "" {
		conf.Mode = fields.Mode
	}

	var vulcanParams string

	if fields.Entry != "" {
		vulcanParams = " --entry " + fields.Entry
	}

	if fields.NodePolyfills != "" {
		_, err := strconv.ParseBool(fields.NodePolyfills)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorPolyfills, fields.NodePolyfills)
		}
		vulcanParams += " --useNodePolyfills " + fields.NodePolyfills
	}

	if fields.OwnWorker != "" {
		_, err := strconv.ParseBool(fields.OwnWorker)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorWorker, fields.OwnWorker)
		}
		vulcanParams += " --useOwnWorker " + fields.OwnWorker
	}

	err = checkArgsJson(cmd)
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
		return vulcan(cmd, conf, vulcanParams)
	}

	if conf.Template == "nextjs" {
		return adapter(cmd, conf)
	}

	return utils.ErrorUnsupportedType
}

func checkArgsJson(cmd *BuildCmd) error {
	workingDir, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	workDirPath := workingDir + "/azion/args.json"
	_, err = cmd.FileReader(workDirPath)
	if err != nil {
		if err := cmd.WriteFile(workDirPath, []byte("{}"), 0644); err != nil {
			logger.Debug("Error while trying to create args.json file", zap.Error(err))
			return fmt.Errorf(utils.ErrorCreateFile.Error(), workDirPath)
		}
	}

	return nil
}
