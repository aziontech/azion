package build

import (
	"fmt"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
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

	conf, err := cmd.GetAzionJsonContent(fields.ProjectPath)
	if err != nil {
		logger.Debug("Error while building your project", zap.Error(err))
		return msg.ErrorBuilding
	}

	if fields.Preset != "" {
		conf.Preset = fields.Preset
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

	if fields.IsFirewall {
		vulcanParams += " --firewall "
	}

	return vulcan(cmd, conf, vulcanParams, fields)

}
