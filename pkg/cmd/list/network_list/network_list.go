package networklist

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/network_list"
	api "github.com/aziontech/azion-cli/pkg/api/network_list"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io              *iostreams.IOStreams
	ListNetworkList func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedNetworkListList, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListNetworkList: func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedNetworkListList, error) {
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
			$ azion list network-list --details
			$ azion list network-list --order_by "id"
			$ azion list network-list --page 1
			$ azion list network-list --page_size 5
			$ azion list network-list --sort "asc"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := PrintTable(cmd, f, list, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetNetworkLists.Error(), err)
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
	netlist, err := list.ListNetworkList(ctx, opts)
	if err != nil {
		return fmt.Errorf(msg.ErrorGetNetworkLists.Error(), err)
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "ACTIVE"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "ACTIVE", "TYPE"}
	}

	if netlist == nil || len(netlist.Results) == 0 {
		return output.Print(&listOut)
	}

	for _, v := range netlist.Results {
		var ln []string
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.GetId()),
				v.GetName(),
				fmt.Sprintf("%v", v.GetActive()),
				v.GetType(),
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
