package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/aziontech/tablecli"
	msg "github.com/aziontech/azion-cli/messages/edge_functions_instances"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	var edgeApplicationID int64 = 0
	cmd := &cobra.Command{
		Use:           msg.EdgeFunctionsInstancesListUsage,
		Short:         msg.EdgeFunctionsInstancesListLongDescription,
		Long:          msg.EdgeFunctionsInstancesListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		    $ azion edge_functions_instances list --application-id 1234123423 --details
		    $ azion edge_functions_instances list --application-id 1234123423 --order-by "id"
		    $ azion edge_functions_instances list --application-id 1234123423 --page 1  
		    $ azion edge_functions_instances list --application-id 1234123423 --page_size 5
		    $ azion edge_functions_instances list -a 1234123423 --sort "asc" 
 			$ azion edge_functions_instances list -a 1234123423" 	
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			var numberPage int64 = opts.Page
			if !cmd.Flags().Changed("application-id") {
				return msg.ErrorMandatoryListFlags
			}
			if !cmd.Flags().Changed("page") && !cmd.Flags().Changed("page_size") {
				for {
					pages, err := PrintTable(cmd, f, opts, &numberPage, edgeApplicationID)
					if numberPage > pages && err == nil {
						return nil
					}
					if err != nil {
						return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
					}
				}
			}

			if _, err := PrintTable(cmd, f, opts, &numberPage, edgeApplicationID); err != nil {
				return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().Int64VarP(&edgeApplicationID, "application-id", "a", 0, msg.EdgeApplicationFlagId)
	cmd.Flags().BoolP("help", "h", false, msg.EdgeFunctionsInstancesListHelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, numberPage *int64, edgeApplicationID int64) (int64, error) {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	applications, err := client.EdgeFuncInstancesList(ctx, opts, edgeApplicationID)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
	}

	tbl := table.New("ID", "NAME")
	table.DefaultWriter = f.IOStreams.Out
	if cmd.Flags().Changed("details") {
		tbl = table.New("ID", "EDGE FUNCTIONS ID", "NAME", "ARGS")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range applications.Results {
		if cmd.Flags().Changed("details") {
			tbl.AddRow(v.Id, v.EdgeFunctionId, v.Name, v.Args)
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
