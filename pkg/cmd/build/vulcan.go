package build

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func (b *BuildCmd) vulcan(vul *vulcanPkg.VulcanPkg, conf *contracts.AzionApplicationOptions, vulcanParams string, fields *contracts.BuildInfo, msgs *[]string) error {
	// checking if vulcan major is correct
	vulcanVer, err := b.CommandRunner(b.f, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	err = vul.CheckVulcanMajor(vulcanVer, b.f, vul)
	if err != nil {
		return err
	}

	command := vul.Command("", "build%s", b.f)
	err = b.runCommand(fmt.Sprintf(command, vulcanParams), msgs)
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	err = b.WriteAzionJsonContent(conf, fields.ProjectPath)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return utils.ErrorWritingAzionJsonFile
	}

	return nil
}
