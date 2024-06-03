package edgefunction

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ListShortDescription,
		Long:          msg.ListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion list edge-function --details
		$ azion list edge-function --order_by "id"
		$ azion list edge-function --page 1  
		$ azion list edge-function --page_size 5
		$ azion list edge-function --sort "asc" 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := PrintTable(cmd, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.ListHelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	functions, err := client.List(ctx, opts)
	if err != nil {
		return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "LANGUAGE", "ACTIVE"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "LANGUAGE", "ACTIVE", "LAST EDITOR", "MODIFIED", "REFERENCE COUNT", "INITIATOR_TYPE"}
	}

	for _, v := range functions.Results {
		ln := []string{
			fmt.Sprintf("%d", v.GetId()),
			v.GetName(),
			v.GetLanguage(),
			fmt.Sprintf("%v", v.GetActive()),
			v.GetLastEditor(),
			v.GetModified(),
			fmt.Sprintf("%d", v.GetReferenceCount()),
			v.GetInitiatorType(),
		}
		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}
