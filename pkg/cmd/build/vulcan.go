package build

import (
	"fmt"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

func vulcan(cmd *BuildCmd, conf *contracts.AzionApplicationOptions) error {
	const command string = "npx --yes edge-functions@1.5.0 build --preset %s --mode %s"

	err := runCommand(cmd, fmt.Sprintf(command, strings.ToLower(conf.Template), strings.ToLower(conf.Mode)))
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	envPath := conf.ProjectRoot + "/.edge/.env"
	fileEnv, err := cmd.FileReader(envPath)
	if err != nil {
		return msg.ErrorEnvFileVulcan
	}
	verIdSlice := strings.Split(string(fileEnv), "=")
	versionID := verIdSlice[1]

	conf.VersionID = versionID

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return utils.ErrorWritingAzionJsonFile
	}

	return nil
}
