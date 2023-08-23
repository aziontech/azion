package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/personal-token"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var id string

	deleteCmd := &cobra.Command{
		Use:           msg.DeleteUsage,
		Short:         msg.DeleteShortDescription,
		Long:          msg.DeleteLongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
		$ azioncli personal_token delete --id 7a187044-4a00-4a4a-93ed-d230900421f3
		$ azioncli personal_token delete -i 7a187044-4a00-4a4a-93ed-d230900421f3
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("id") {
				return msg.ErrorMissingIDArgumentDelete
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			err := client.Delete(context.Background(), id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.DeleteOutputSuccess, id)

			return nil
		},
	}

	deleteCmd.Flags().StringVarP(&id, "id", "i", "", msg.FlagID)
	deleteCmd.Flags().BoolP("help", "h", false, msg.DeleteHelpFlag)

	return deleteCmd
}
