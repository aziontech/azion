package edge_applications

import (
	"context"
	"fmt"

	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azionapi-go-sdk/edgeapplications"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	api "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io           *iostreams.IOStreams
	ListEdgeApps func(context.Context, *contracts.ListOptions) (*edgeapplications.GetApplicationsResponse, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListEdgeApps: func(ctx context.Context, opts *contracts.ListOptions) (*edgeapplications.GetApplicationsResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.List(ctx, opts)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           "edge-application",
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list edge-application
			$ azion list edge-application --details
			$ azion list edge-application --page 1
			$ azion list edge-application --page-size 5
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := PrintTable(cmd, list, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetAll.Error(), err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolP("help", "h", false, msg.HelpFlag)
	cmdutil.AddAzionApiFlags(cmd, opts)

	return cmd
}

func PrintTable(cmd *cobra.Command, list *ListCmd, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	ctx := context.Background()

	resp, err := list.ListEdgeApps(ctx, opts)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "ACTIVE"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "ACTIVE", "LAST EDITOR", "LAST MODIFIED", "DEBUG RULES"}
	}

	for _, v := range resp.Results {
		var ln []string
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.Id),
				utils.TruncateString(v.Name),
				fmt.Sprintf("%v", v.Active),
				v.LastEditor,
				v.LastModified,
				fmt.Sprintf("%v", v.DebugRules),
			}
		} else {
			ln = []string{
				fmt.Sprintf("%d", v.Id),
				utils.TruncateString(v.Name),
				fmt.Sprintf("%v", v.Active),
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
