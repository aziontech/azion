package streams

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/streams"
	api "github.com/aziontech/azion-cli/pkg/api/data_stream"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io          *iostreams.IOStreams
	ListStreams func(context.Context, *contracts.ListOptions) (*sdk.PaginatedDataStreamList, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListStreams: func(ctx context.Context, opts *contracts.ListOptions) (*sdk.PaginatedDataStreamList, error) {
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
			$ azion list streams --details
			$ azion list streams --order_by "id"
			$ azion list streams --page 1
			$ azion list streams --page_size 5
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := PrintTable(cmd, f, list, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetStreams.Error(), err)
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
	streams, err := list.ListStreams(ctx, opts)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "ACTIVE"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "ACTIVE", "LAST EDITOR", "LAST MODIFIED"}
	}

	if streams == nil || len(streams.Results) == 0 {
		return output.Print(&listOut)
	}
	for _, v := range streams.Results {
		var ln []string
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.GetId()),
				v.GetName(),
				fmt.Sprintf("%v", v.GetActive()),
				v.GetLastEditor(),
				v.GetLastModified().String(),
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
