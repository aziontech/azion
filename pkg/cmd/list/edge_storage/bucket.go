package edge_storage

import (
	"context"
	"strings"

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

var bucketName string

func NewBucket(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           msg.USAGE_BUCKET,
		Short:         msg.SHORT_DESCRIPTION_LIST_BUCKET,
		Long:          msg.LONG_DESCRIPTION_LIST_BUCKET,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          runE(f, opts),
	}

	flags := cmd.Flags()
	addFlags(flags, opts)
	return cmd
}

func runE(f *cmdutil.Factory, opts *contracts.ListOptions) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
		return PrintTable(client, f, opts)
	}
}

func PrintTable(client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	c := context.Background()

	resp, err := client.ListBucket(c, opts)
	if err != nil {
		return err
	}

	tbl := table.New("NAME", "EDGE ACCESS")
	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range resp.Results {
		tbl.AddRow(v.GetName(), v.GetEdgeAccess())
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})

	// print the header only in the first flow
	if opts.Page == 1 {
		logger.PrintHeader(tbl, format)
	}

	for _, row := range tbl.GetRows() {
		logger.PrintRow(tbl, format, row)
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, opts *contracts.ListOptions) {
	flags.Int64Var(&opts.Page, "page", 1, general.ApiListFlagPage)
	flags.Int64Var(&opts.PageSize, "page-size", 10, general.ApiListFlagPageSize)
	flags.BoolP("help", "h", false, "")
}
