package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_functions"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/printer"
	"github.com/spf13/cobra"

	table "github.com/MaxwelMazur/tablecli"
	"github.com/fatih/color"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           msg.EdgeFunctionListUsage,
		Short:         msg.EdgeFunctionListShortDescription,
		Long:          msg.EdgeFunctionListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion edge_functions list --details
		$ azion edge_functions list --order_by "id"
		$ azion edge_functions list --page 1  
		$ azion edge_functions list --page_size 5
		$ azion edge_functions list --sort "asc" 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var numberPage int64 = opts.Page
			if !cmd.Flags().Changed("page") && !cmd.Flags().Changed("page_size") {
				for {
					pages, err := PrintTable(f, opts, &numberPage)
					if numberPage > pages {
						return nil
					}
					if err != nil {
						return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
					}
				}
			}

			if _, err := PrintTable(f, opts, &numberPage); err != nil {
				return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.EdgeFunctionListHelpFlag)
	return cmd
}

func PrintTable(f *cmdutil.Factory, opts *contracts.ListOptions, numberPage *int64) (int64, error) {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	table.DefaultWriter = f.IOStreams.Out
	tbl := table.New("ID", "NAME", "LANGUAGE", "ACTIVE")
	fields := []string{"GetId()", "GetName()", "GetLanguage()", "GetActive()"}

	functions, pages, err := client.List(ctx, opts)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
	}

	if opts.Details {
		fields = append(fields, "GetLastEditor()", "GetModified()", "GetReferenceCount()", "GetInitiatorType()")
		tbl = table.New("ID", "NAME", "LANGUAGE", "ACTIVE", "LAST EDITOR", "MODIFIED", "REFERENCE COUNT", "INITIATOR_TYPE")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	rows := printer.BuildRows(functions, fields)
	for _, row := range rows {
		if len(row) == 8 {
			tbl.AddRow(row[0], row[1], row[2], row[3], row[4], row[5], row[6], row[7])
		} else {
			tbl.AddRow(row[0], row[1], row[2], row[3])
		}
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	if *numberPage == 1 {
		tbl.PrintHeader(format)
	}

	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	*numberPage += 1
	opts.Page = *numberPage
	return pages, nil
}
