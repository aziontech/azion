package list

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"strings"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.EdgeApplicationsListUsage,
		Short:         msg.EdgeApplicationsListShortDescription,
		Long:          msg.EdgeApplicationsListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli edge_functions list --details
		$ azioncli edge_functions list --order_by "id"
		$ azioncli edge_functions list --page 1  
		$ azioncli edge_functions list --page_size 5
		$ azioncli edge_functions list --sort "asc" 
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			var numberPage int64 = opts.Page
			if !insertTheFlagPage(cmd) {
				for {
					pages, err := PrintTable(cmd, f, opts, &numberPage)
					if numberPage > pages {
						return nil
					}
					if err != nil {
						return fmt.Errorf(msg.ErrorGetApplication.Error(), err)
					}
				}
			}

			if _, err := PrintTable(cmd, f, opts, &numberPage); err != nil {
				return fmt.Errorf(msg.ErrorGetApplication.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.EdgeApplicationsListHelpFlag)
	return cmd
}

func insertTheFlagDetails(cmd *cobra.Command) bool {
	return cmd.Flags().Changed("details")
}

func insertTheFlagPage(cmd *cobra.Command) bool {
	return cmd.Flags().Changed("page")
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, numberPage *int64) (int64, error) {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	applications, err := client.List(ctx, opts)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorGetApplication.Error(), err)
	}

	tbl := table.New("ID", "NAME")
	table.DefaultWriter = f.IOStreams.Out
	if insertTheFlagDetails(cmd) {
		tbl = table.New("ID", "NAME", "ACTIVE")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range applications.Results {
		if insertTheFlagDetails(cmd) {
			tbl.AddRow(v.Id, v.Name, v.Active)
		} else {
			tbl.AddRow(v.Id, v.Name)
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
	return applications.TotalPages, nil
}
