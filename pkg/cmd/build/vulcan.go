package build

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/build"
)

func vulcan(cmd *BuildCmd, typeLang, mode string) error {
	const command string = "npx --yes edge-functions@1.0.0 build --preset %s --mode %s"

	err := runCommand(cmd, fmt.Sprintf(command, typeLang, mode))
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}

	return nil
}
