package list

import (
	"context"
	"io"
	"strings"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	"github.com/aziontech/azion-cli/messages/general"
	msg "github.com/aziontech/azion-cli/messages/variables"
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/printer"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	listCmd := &cobra.Command{
		Use:           msg.VariablesListUsage,
		Short:         msg.VariablesListShortDescription,
		Long:          msg.VariablesListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli variables list --details
		$ azioncli variables list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			if err := listAllVariables(client, f.IOStreams.Out, opts); err != nil {
				return err
			}
			return nil
		},
	}

	listCmd.Flags().BoolVar(&opts.Details, "details", false, general.ApiListFlagDetails)
	listCmd.Flags().BoolP("help", "h", false, msg.VariablesListHelpFlag)
	return listCmd
}

func listAllVariables(client *api.Client, out io.Writer, opts *contracts.ListOptions) error {
	c := context.Background()

	resp, err := client.List(c)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "KEY", "VALUE")
	tbl.WithWriter(out)
	fields := []string{"GetUuid()", "GetKey()", "GetValue()"}

	if opts.Details {
		fields = append(fields, "GetSecret()", "GetLastEditor()")
		tbl = table.New("ID", "KEY", "VALUE", "SECRET", "LAST EDITOR")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	rows := printer.BuildRows(resp, fields)
	for _, row := range rows {
		if len(row) == 5 {
			tbl.AddRow(row[0], row[1], row[2], row[3], row[4])
		} else {
			tbl.AddRow(row[0], row[1], row[2])
		}
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	tbl.PrintHeader(format)

	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	return nil
}
