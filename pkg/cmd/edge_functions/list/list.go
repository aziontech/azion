package list

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_functions"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/printer"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           msg.EdgeFunctionListUsage,
		Short:         msg.EdgeFunctionListShortDescription,
		Long:          msg.EdgeFunctionListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions list --details
		$ azioncli edge_functions list --order_by "id"
		$ azioncli edge_functions list --page 1  
		$ azioncli edge_functions list --page_size 5
		$ azioncli edge_functions list --sort "asc" 

        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()

			fields := []string{"GetId()", "GetName()", "GetLanguage()", "GetActive()"}
			headers := []string{"ID", "NAME", "LANGUAGE", "ACTIVE"}

			functions, err := client.List(ctx, opts)
			if err != nil {
				return fmt.Errorf("%w: %s", msg.ErrorGetFunctions, err)
			}

			out := f.IOStreams.Out
			tp := printer.NewTab(out)
			if opts.Details {
				fields = append(fields, "GetLastEditor()", "GetModified()", "GetReferenceCount()", "GetInitiatorType()")
				headers = append(headers, "LAST EDITOR", "MODIFIED", "REFERENCE COUNT", "INITIATOR_TYPE")
			}

			tp.PrintWithHeaders(functions, fields, headers)

			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.EdgeFunctionListHelpFlag)

	return cmd
}
