package edge_applications

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	"github.com/aziontech/azion-cli/messages/general"
	msg "github.com/aziontech/azion-cli/messages/list/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           msg.ListUsage,
		Short:         msg.ListShortDescription,
		Long:          msg.ListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
            $ azion list edge_applications
            $ azion list edge_applications --details
            $ azion list edge_applications --page 1 
            $ azion list edge_applications --page_size 5
        `),
		RunE: func(cmd *cobra.Command, _ []string) error {
			var numberPage int64 = opts.Page

			client := api.NewClient(f.HttpClient,
				f.Config.GetString("api_url"),
				f.Config.GetString("token"),
			)

			if !cmd.Flags().Changed("page") && !cmd.Flags().Changed("page_size") {
				for {
					pages, err := PrintTable(client, f, opts, &numberPage)
					if numberPage > pages && err == nil {
						return nil
					}
					if err != nil {
						return fmt.Errorf(msg.ErrorGetAll.Error(), err)
					}
				}
			}

			if _, err := PrintTable(client, f, opts, &numberPage); err != nil {
				return fmt.Errorf(msg.ErrorGetAll.Error(), err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&opts.Page, "page", 1, general.ApiListFlagPage)
	flags.Int64Var(&opts.PageSize, "page_size", 10, general.ApiListFlagPageSize)
	flags.BoolVar(&opts.Details, "details", false, general.ApiListFlagDetails)
	flags.BoolP("help", "h", false, msg.ListHelpFlag)
	return cmd
}

func PrintTable(client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions, numberPage *int64) (int64, error) {
	c := context.Background()

	resp, err := client.List(c, opts)
	if err != nil {
		return 0, err
	}

	tbl := table.New("ID", "NAME", "ACTIVE")
	tbl.WithWriter(f.IOStreams.Out)

	if opts.Details {
		tbl = table.New("ID", "NAME", "DEBUG RULES", "LAST EDITOR", "LAST MODIFIED", "ACTIVE")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range resp.Results {
		tbl.AddRow(
			v.Id,
			utils.TruncateString(v.Name),
			v.DebugRules,
			v.LastEditor,
			v.LastModified,
			v.Active,
		)
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})

	// print the header only in the first flow
	if *numberPage == 1 {
		tbl.PrintHeader(format)
	}

	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	*numberPage += 1
	opts.Page = *numberPage

	return resp.TotalPages, nil
}
