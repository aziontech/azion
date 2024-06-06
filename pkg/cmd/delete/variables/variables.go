package variables

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/variables"
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var variableID string

	deleteCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.DeleteShortDescription,
		Long:          msg.DeleteLongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
		$ azion delete variables -h
		$ azion delete variables
		$ azion delete variables --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("variable-id") {
				answers, err := utils.AskInput(msg.AskVariableID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				variableID = answers
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			err := client.Delete(context.Background(), variableID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteVariable.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeleteOutputSuccess, variableID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)

		},
	}

	deleteCmd.Flags().StringVar(&variableID, "variable-id", "", msg.FlagVariableID)
	deleteCmd.Flags().BoolP("help", "h", false, msg.DeleteHelpFlag)

	return deleteCmd
}
