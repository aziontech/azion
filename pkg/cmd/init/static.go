package init

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	thoth "github.com/aziontech/go-thoth"
)

func initStatic(cmd *InitCmd, info *InitInfo) error {
	shouldFetchTemplates, err := shouldFetch(cmd, info)
	if err != nil {
		return err
	}

	if shouldFetchTemplates {
		if info.GlobalFlagAll {
			info.Name = thoth.GenerateName()
		} else {
			if info.Name == "" {
				projName, err := askForInput(msg.InitProjectQuestion, thoth.GenerateName())
				if err != nil {
					return err
				}

				info.Name = projName
			}
		}
		if err = cmd.createTemplateAzion(info); err != nil {
			return err
		}

		logger.FInfo(cmd.Io.Out, fmt.Sprintf(msg.EdgeApplicationsInitSuccessful+"\n", info.Name))
	}

	logger.FInfo(cmd.Io.Out, `  [ General Instructions ]
    [ Usage ]
    - Publish Command: publish page static`)

	return nil
}
