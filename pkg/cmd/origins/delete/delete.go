package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/origins"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/spf13/cobra"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
)



func NewCmd(f *cmdutil.Factory) *cobra.Command {
  var applicationID int64
  var originKey     string
	cmd := &cobra.Command{
		Use:           msg.OriginsDeleteUsage,
		Short:         msg.OriginsDeleteShortDescription,
		Long:          msg.OriginsDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azioncli origins delete --application-id 1234 --origin-key 03a6e7bf-8e26-49c7-a66e-ab8eaa425086
		  $ azioncli origins delete -a 1234 -o 03a6e7bf-8e26-49c7-a66e-ab8eaa425086
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("origin-key") {
				return msg.ErrorMissingArgumentsDelete
			}
			if err := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token")).
        DeleteOrigins(context.Background(), applicationID, originKey); err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.OriginsDeleteOutputSuccess, originKey)
      return nil
		},
	}

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, msg.OriginsDeleteFlagApplicationID)
	cmd.Flags().StringVarP(&originKey, "origin-key", "o", "", msg.OriginsDeleteFlagOriginKey)
	cmd.Flags().BoolP("help", "h", false, msg.OriginsDeleteHelpFlag)
	return cmd
}
