package connector

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/connector"
	api "github.com/aziontech/azion-cli/pkg/api/connector"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io             *iostreams.IOStreams
	ListConnectors func(context.Context, *contracts.ListOptions) (*sdk.PaginatedConnectorPolymorphicList, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListConnectors: func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedConnectorPolymorphicList, error) {
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
			$ azion list connector --details
			$ azion list connector --order_by "id"
			$ azion list connector --page 1
			$ azion list connector --page_size 5
			$ azion list connector --sort "asc"
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
	connectors, err := list.ListConnectors(ctx, opts)
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

	if connectors == nil || len(connectors.Results) == 0 {
		return output.Print(&listOut)
	}
	for _, v := range connectors.Results {
		var ln []string
		if v.ConnectorHTTP != nil {
			vObj := v.ConnectorHTTP
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
		} else if v.ConnectorLiveIngest != nil {
			vObj := v.ConnectorLiveIngest
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
		} else if v.ConnectorStorage != nil {
			vObj := v.ConnectorStorage
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
