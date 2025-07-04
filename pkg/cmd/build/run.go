package build

import (
	"fmt"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"go.uber.org/zap"
)

func (b *BuildCmd) run(fields *contracts.BuildInfo, msgs *[]string) error {
	logger.Debug("Running build command")

	var err error

	conf, err := b.GetAzionJsonContent(fields.ProjectPath)
	if err != nil {
		logger.Debug("Error while building your project", zap.Error(err))
		return msg.ErrorBuilding
	}

	var vulcanParams string

	if fields.Preset != "" {
		vulcanParams = " --preset " + fields.Preset
		conf.Preset = fields.Preset
	} else {
		vulcanParams = " --preset " + conf.Preset
	}

	if fields.Entry != "" {
		vulcanParams += " --entry " + fields.Entry
	}

	if fields.NodePolyfills != "" {
		_, err := strconv.ParseBool(fields.NodePolyfills)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorPolyfills, fields.NodePolyfills)
		}
		vulcanParams += " --polyfills " + fields.NodePolyfills
	}

	if fields.OwnWorker != "" {
		_, err := strconv.ParseBool(fields.OwnWorker)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorWorker, fields.OwnWorker)
		}
		vulcanParams += " --worker " + fields.OwnWorker
	}

	if fields.SkipFramework {
		vulcanParams += " --skip-framework-build"
	}

	vul := vulcanPkg.NewVulcan()
	return b.vulcan(vul, conf, vulcanParams, fields, msgs)
}
