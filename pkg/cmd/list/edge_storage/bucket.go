package edge_storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"

	table "github.com/MaxwelMazur/tablecli"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/messages/general"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
)

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
	tbl := table.New("NAME", "EDGE ACCESS")
	tbl.WithWriter(b.Factory.IOStreams.Out)
	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, v := range resp.Results {
		tbl.AddRow(v.GetName(), v.GetEdgeAccess())
	}
	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	// print the header only in the first flow
	if b.Options.Page == 1 {
		logger.PrintHeader(tbl, format)
	}
	for _, row := range tbl.GetRows() {
		logger.PrintRow(tbl, format, row)
	}
	return nil
}

func (b *Bucket) AddFlags(flags *pflag.FlagSet) {
	flags.Int64Var(&b.Options.Page, "page", 1, general.ApiListFlagPage)
	flags.Int64Var(&b.Options.PageSize, "page-size", 10, general.ApiListFlagPageSize)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_LIST_BUCKET)
}
