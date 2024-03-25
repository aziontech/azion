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
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
)

func NewObject(f *cmdutil.Factory) *cobra.Command {
	object := &Objects{
		Factory: f,
		Options: &contracts.ListOptions{},
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         msg.SHORT_DESCRIPTION_LIST_OBJECT,
		Long:          msg.LONG_DESCRIPTION_LIST_OBJECT,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          object.RunE,
		Example:       heredoc.Doc(msg.EXAMPLE_LIST_OBJECT),
	}
	object.AddFlags(cmd.Flags())
	return cmd
}

func (b *Objects) RunE(cmd *cobra.Command, args []string) error {
	client := api.NewClient(
		b.Factory.HttpClient,
		b.Factory.Config.GetString("storage_url"),
		b.Factory.Config.GetString("token"))
	return b.PrintTable(client)
}

func (b *Objects) PrintTable(client *api.Client) error {
	c := context.Background()
	resp, err := client.ListObject(c, b.BucketName, b.Options)
	if err != nil {
		return fmt.Errorf(msg.ERROR_LIST_BUCKET, err)
	}
	tbl := table.New("KEY", "LAST MODIFIED")
	tbl.WithWriter(b.Factory.IOStreams.Out)
	if b.Options.Details {
		tbl = table.New("KEY", "LAST MODIFIED", "SIZE", "ETAG")
	}
	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	for _, v := range resp.Results {
		tbl.AddRow(v.GetKey(), v.GetLastModified(), v.GetSize(), v.GetEtag())
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

func (b *Objects) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&b.BucketName, "bucket-name", "", msg.FLAG_NAME_BUCKET)
	flags.BoolVar(&b.Options.Details, "details", false, msg.FLAG_HELP_DETAILS_OBJECTS)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_LIST_OBJECT)
}
