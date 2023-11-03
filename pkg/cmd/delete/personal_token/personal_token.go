package personaltoken

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/delete/personal_token"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var id string

	deleteCmd := &cobra.Command{
		Use:           personaltoken.Usage,
		Short:         personaltoken.ShortDescription,
		Long:          personaltoken.LongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
		$ azion delete personal-token --id 1234-123-321
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("id") {
				answer, err := utils.AskInput(personaltoken.AskDeleteInput)
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
				return fmt.Errorf(personaltoken.ErrorFailToDelete.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, personaltoken.OutputSuccess, id)

			return nil
		},
	}

	deleteCmd.Flags().StringVar(&id, "id", "", personaltoken.FlagID)
	deleteCmd.Flags().BoolP("help", "h", false, personaltoken.HelpFlag)

	return deleteCmd
}
