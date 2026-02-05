package kv

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/kv"
	api "github.com/aziontech/azion-cli/pkg/api/kv"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/kv-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io             *iostreams.IOStreams
	ListNamespaces func(context.Context, *contracts.ListOptions) (*sdk.NamespaceList, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListNamespaces: func(ctx context.Context, opts *contracts.ListOptions) (*sdk.NamespaceList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.List(ctx, opts)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list kv
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := PrintTable(cmd, f, list, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetKv.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, list *ListCmd, opts *contracts.ListOptions) error {
	ctx := context.Background()
	namespaces, err := list.ListNamespaces(ctx, opts)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"NAME"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"NAME", "CREATED AT", "LAST MODIFIED"}
	}

	if namespaces == nil || len(namespaces.Results) == 0 {
		return output.Print(&listOut)
	}
	for _, v := range namespaces.Results {
		var ln []string
		if opts.Details {
			ln = []string{
				v.GetName(),
				v.GetCreatedAt().String(),
				v.GetLastModified().String(),
			}
		} else {
			ln = []string{
				v.GetName(),
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
