package list

import (
	"context"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_functions/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/printer"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           "list [flags]",
		Short:         "Lists your account's Edge Functions",
		Long:          "Lists your account's Edge Functions",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions list [--details]
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()

			fields := []string{"GetId()", "GetName()", "GetLanguage()", "GetActive()"}
			headers := []string{"ID", "NAME", "LANGUAGE", "ACTIVE"}

			functions, err := client.List(ctx, opts)
			if err != nil {
				return errmsg.ErrorGetFunctions
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

	return cmd
}
