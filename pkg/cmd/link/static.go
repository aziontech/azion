package link

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/link"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	thoth "github.com/aziontech/go-thoth"
	"github.com/spf13/cobra"
)

func initStatic(cmd *LinkCmd, info *LinkInfo, options *contracts.AzionApplicationOptions, c *cobra.Command) error {
	shouldFetchTemplates, err := shouldFetch(cmd, info)
	if err != nil {
		return err
	}

	if shouldFetchTemplates {
		if info.GlobalFlagAll {
			info.Name = thoth.GenerateName()
		} else {
			if !c.Flags().Changed("name") {
				projName, err := askForInput(msg.LinkProjectQuestion, thoth.GenerateName())
				if err != nil {
					return err
				}

				info.Name = projName
			}
		}
		if err = cmd.createTemplateAzion(info); err != nil {
			return err
		}

		logger.FInfo(cmd.Io.Out, fmt.Sprintf(msg.EdgeApplicationsLinkSuccessful+"\n", info.Name))
	}

	logger.FInfo(cmd.Io.Out, `	[ General Instructions ]
    - Requirements:
        - Tools: npm
    [ Usage ]
    	- Dev Command: azion dev
    	- Build Command: azion build
    	- Deploy Command: azion deploy`) //nolint:all`)

	return nil
}
