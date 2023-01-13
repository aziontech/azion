package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var application_id string
	cmd := &cobra.Command{
		Use:           msg.EdgeApplicationDeleteUsage,
		Short:         msg.EdgeApplicationDeleteShortDescription,
		Long:          msg.EdgeApplicationDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_applications delete --application-id 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				return msg.ErrorMissingApplicationIdArgument
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, application_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.EdgeApplicationDeleteOutputSuccess, application_id)

			return nil
		},
	}

	cmd.Flags().StringVarP(&application_id, "application-id", "a", "", msg.EdgeApplicationFlagId)
	cmd.Flags().BoolP("help", "h", false, msg.EdgeApplicationDeleteHelpFlag)

	return cmd
}
