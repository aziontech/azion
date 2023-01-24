package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/domains"
	api "github.com/aziontech/azion-cli/pkg/api/domains"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var function_id int64
	cmd := &cobra.Command{
		Use:           msg.DomainDeleteUsage,
		Short:         msg.DomainDeleteShortDescription,
		Long:          msg.DomainDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli domains delete --domain-id 1234
		$ azioncli domains delete -d 1234
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("domain-id") {
				return msg.ErrorMissingDomainIdArgumentDelete
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, function_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteDomain.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.DomainDeleteOutputSuccess, function_id)

			return nil
		},
	}

	cmd.Flags().Int64VarP(&function_id, "domain-id", "d", 0, msg.DomainFlagId)
	cmd.Flags().BoolP("help", "h", false, msg.DomainDeleteHelpFlag)

	return cmd
}
