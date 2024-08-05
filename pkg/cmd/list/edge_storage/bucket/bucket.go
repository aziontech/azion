package bucket

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/messages/general"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/output"
)

type Bucket struct {
	Options *contracts.ListOptions
	Factory *cmdutil.Factory
}

func NewBucket(f *cmdutil.Factory) *cobra.Command {
	bucket := &Bucket{
		Factory: f,
		Options: &contracts.ListOptions{},
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_BUCKET,
		Short:         msg.SHORT_DESCRIPTION_LIST_BUCKET,
		Long:          msg.LONG_DESCRIPTION_LIST_BUCKET,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          bucket.RunE,
		Example:       heredoc.Doc(msg.EXAMPLE_LIST_BUCKET),
	}

	bucket.AddFlags(cmd.Flags())
	return cmd
}

func (b *Bucket) RunE(cmd *cobra.Command, args []string) error {
	client := api.NewClient(
		b.Factory.HttpClient,
		b.Factory.Config.GetString("storage_url"),
		b.Factory.Config.GetString("token"))
	return b.PrintTable(client)
}

func (b *Bucket) PrintTable(client *api.Client) error {
	c := context.Background()
	resp, err := client.ListBucket(c, b.Options)
	if err != nil {
		return fmt.Errorf(msg.ERROR_LIST_BUCKET, err)
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"NAME", "EDGE ACCESS"}
	listOut.Out = b.Factory.IOStreams.Out
	listOut.Flags = b.Factory.Flags

	for _, v := range resp.Results {
		ln := []string{
			v.GetName(),
			string(v.GetEdgeAccess()),
		}
		listOut.Lines = append(listOut.Lines, ln)
	}
	return output.Print(&listOut)
}

func (b *Bucket) AddFlags(flags *pflag.FlagSet) {
	flags.Int64Var(&b.Options.Page, "page", 1, general.ApiListFlagPage)
	flags.Int64Var(&b.Options.PageSize, "page-size", 50, general.ApiListFlagPageSize)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_LIST_BUCKET)
}
