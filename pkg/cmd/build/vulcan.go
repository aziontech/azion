package build

import (
	"fmt"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/build"
	"github.com/aziontech/azion-cli/pkg/contracts"
)

func vulcan(cmd *BuildCmd, conf *contracts.AzionApplicationOptions, path string) error {
	const command string = "npx --yes edge-functions@1.0.0 build --preset %s --mode %s"

	err := runCommand(cmd, fmt.Sprintf(command, strings.ToLower(conf.Template), strings.ToLower(conf.Mode)))
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	envPath := path + "/.edge/.env"
	fileEnv, err := cmd.FileReader(envPath)
	if err != nil {
		return msg.ErrorEnvFileVulcan
	}
	verIdSlice := strings.Split(string(fileEnv), "=")
	versionID := verIdSlice[1]

	conf.VersionID = versionID

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		return nil
	}

	return nil
}
