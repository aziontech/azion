package edgefunction

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io            *iostreams.IOStreams
	ListFunctions func(context.Context, *contracts.ListOptions) (*sdk.PaginatedEdgeFunctionsList, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListFunctions: func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedEdgeFunctionsList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.List(ctx, opts)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
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
			if err := PrintTable(cmd, f, list, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.ListHelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, list *ListCmd, opts *contracts.ListOptions) error {
	ctx := context.Background()
	functions, err := list.ListFunctions(ctx, opts)
	if err != nil {
		return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "ACTIVE"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "ACTIVE", "LAST EDITOR", "REFERENCE COUNT", "EXECUTION ENVIRONMENT"}
	}

	for _, v := range functions.Results {
		var ln []string
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.GetId()),
				v.GetName(),
				fmt.Sprintf("%v", v.GetActive()),
				v.GetLastEditor(),
				fmt.Sprintf("%d", v.GetReferenceCount()),
				v.GetExecutionEnvironment(),
			}
		} else {
			ln = []string{
				fmt.Sprintf("%d", v.GetId()),
				v.GetName(),
				fmt.Sprintf("%v", v.GetActive()),
			}
		}
		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	listCmd := NewListCmd(f)
	return NewCobraCmd(listCmd, f)
}
