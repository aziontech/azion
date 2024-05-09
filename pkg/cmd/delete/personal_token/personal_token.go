package personaltoken

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/personal_token"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var id string

	deleteCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
		$ azion delete personal-token --id 1234-123-321
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("id") {
				answer, err := utils.AskInput(msg.AskDeleteInput)
				if err != nil {
					return err
				}

				id = answer
			}

			if utils.IsEmpty(id) {
				return utils.ErrorArgumentIsEmpty
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			err := client.Delete(context.Background(), id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:         fmt.Sprintf(msg.OutputSuccess, id),
				Out:         f.IOStreams.Out,
				FlagOutPath: f.Out,
				FlagFormat:  f.Format,
			}
			return output.Print(&deleteOut)

		},
	}

	deleteCmd.Flags().StringVar(&id, "id", "", msg.FlagID)
	deleteCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return deleteCmd
}
