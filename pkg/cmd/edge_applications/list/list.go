package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

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
		Use:           msg.ListUsage,
		Short:         msg.ListShortDescription,
		Long:          msg.ListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		    $ azion edge_applications list --details
		    $ azion edge_applications list --order_by "id"
		    $ azion edge_applications list --page 1  
		    $ azion edge_applications list --page_size 5
		    $ azion edge_applications list --sort "asc" 
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			var numberPage int64 = opts.Page
			if !cmd.Flags().Changed("page") && !cmd.Flags().Changed("page_size") {
				for {
					pages, err := PrintTable(cmd, f, opts, &numberPage)
					if numberPage > pages && err == nil {
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
	cmd.Flags().BoolP("help", "h", false, msg.ListHelpFlag)
	return cmd
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
	if cmd.Flags().Changed("details") {
		tbl = table.New("ID", "NAME", "ACTIVE")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range applications.Results {
		if cmd.Flags().Changed("details") {
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
	f.IOStreams.Out = table.DefaultWriter
	return applications.TotalPages, nil
}
