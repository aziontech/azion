package build

import (
	"fmt"
	"strconv"
	"strings"

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

	var paramsBuilder strings.Builder

	if fields.Preset != "" {
		paramsBuilder.WriteString(" --preset ")
		paramsBuilder.WriteString(fields.Preset)
		conf.Preset = fields.Preset
	} else {
		paramsBuilder.WriteString(" --preset ")
		paramsBuilder.WriteString(conf.Preset)
	}

	if fields.Entry != "" {
		paramsBuilder.WriteString(" --entry ")
		paramsBuilder.WriteString(fields.Entry)
	}

	if fields.NodePolyfills != "" {
		_, err := strconv.ParseBool(fields.NodePolyfills)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorPolyfills, fields.NodePolyfills)
		}
		paramsBuilder.WriteString(" --polyfills ")
		paramsBuilder.WriteString(fields.NodePolyfills)
	}

	if fields.OwnWorker != "" {
		_, err := strconv.ParseBool(fields.OwnWorker)
		if err != nil {
			return fmt.Errorf("%w: %s", msg.ErrorWorker, fields.OwnWorker)
		}
		paramsBuilder.WriteString(" --worker ")
		paramsBuilder.WriteString(fields.OwnWorker)
	}

	if fields.SkipFramework {
		paramsBuilder.WriteString(" --skip-framework-build")
	}

	vul := vulcanPkg.NewVulcan()
	return b.vulcan(vul, conf, paramsBuilder.String(), fields, msgs)
}
