package build

import (
	"fmt"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	vul "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func vulcan(cmd *BuildCmd, conf *contracts.AzionApplicationOptions, vulcanParams string) error {
	// checking if vulcan major is correct
	vulcanVer, err := cmd.CommandRunner(cmd.f, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	err = vul.CheckVulcanMajor(vulcanVer, cmd.f)
	if err != nil {
		return err
	}

	command := vul.Command("", "build --preset %s --mode %s%s")

	err = runCommand(cmd, fmt.Sprintf(command, strings.ToLower(conf.Template), strings.ToLower(conf.Mode), vulcanParams))
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	versionID := cmd.VersionID()

	conf.Prefix = versionID

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return utils.ErrorWritingAzionJsonFile
	}

	return nil
}
