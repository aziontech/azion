package edgeconnector

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_connector"
	api "github.com/aziontech/azion-cli/pkg/api/edge_connector"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io             *iostreams.IOStreams
	ListConnectors func(context.Context, *contracts.ListOptions) (*sdk.PaginatedEdgeConnectorPolymorphicList, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListConnectors: func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedEdgeConnectorPolymorphicList, error) {
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
			$ azion list edge-connector --details
			$ azion list edge-connector --order_by "id"
			$ azion list edge-connector --page 1
			$ azion list edge-connector --page_size 5
			$ azion list edge-connector --sort "asc"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := PrintTable(cmd, f, list, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetConnectors.Error(), err)
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
	functions, err := list.ListConnectors(ctx, opts)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "TYPE", "ACTIVE"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "TYPE", "ACTIVE", "LAST EDITOR", "LAST MODIFIED"}
	}

	for _, v := range functions.Results {
		var ln []string
		if v.EdgeConnectorHTTP != nil {
			vObj := v.EdgeConnectorHTTP
			if opts.Details {
				ln = []string{
					fmt.Sprintf("%d", vObj.Id),
					vObj.GetName(),
					vObj.GetType(),
					fmt.Sprintf("%v", vObj.GetActive()),
					vObj.GetLastEditor(),
					vObj.GetLastModified().String(),
				}
			} else {
				ln = []string{
					fmt.Sprintf("%d", vObj.Id),
					vObj.GetName(),
					vObj.GetType(),
					fmt.Sprintf("%v", vObj.GetActive()),
				}
			}
			listOut.Lines = append(listOut.Lines, ln)
		} else if v.EdgeConnectorLiveIngest != nil {
			vObj := v.EdgeConnectorLiveIngest
			if opts.Details {
				ln = []string{
					fmt.Sprintf("%d", vObj.Id),
					vObj.GetName(),
					vObj.GetType(),
					fmt.Sprintf("%v", vObj.GetActive()),
					vObj.GetLastEditor(),
					vObj.GetLastModified().String(),
				}
			} else {
				ln = []string{
					fmt.Sprintf("%d", vObj.Id),
					vObj.GetName(),
					vObj.GetType(),
					fmt.Sprintf("%v", vObj.GetActive()),
				}
			}
			listOut.Lines = append(listOut.Lines, ln)
		} else if v.EdgeConnectorStorage != nil {
			vObj := v.EdgeConnectorStorage
			if opts.Details {
				ln = []string{
					fmt.Sprintf("%d", vObj.Id),
					vObj.GetName(),
					vObj.GetType(),
					fmt.Sprintf("%v", vObj.GetActive()),
					vObj.GetLastEditor(),
					vObj.GetLastModified().String(),
				}
			} else {
				ln = []string{
					fmt.Sprintf("%d", vObj.Id),
					vObj.GetName(),
					vObj.GetType(),
					fmt.Sprintf("%v", vObj.GetActive()),
				}
			}
			listOut.Lines = append(listOut.Lines, ln)

		}
	}

	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	listCmd := NewListCmd(f)
	return NewCobraCmd(listCmd, f)
}
