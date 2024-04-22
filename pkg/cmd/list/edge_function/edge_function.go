package edgefunction

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_function"
	api "github.com/aziontech/azion-cli/pkg/api/edge_function"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/spf13/cobra"

	"github.com/fatih/color"
	table "github.com/maxwelbm/tablecli"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ListShortDescription,
		Long:          msg.ListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion list edge-function --details
		$ azion list edge-function --order_by "id"
		$ azion list edge-function --page 1  
		$ azion list edge-function --page_size 5
		$ azion list edge-function --sort "asc" 
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			if err := PrintTable(cmd, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.ListHelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	for {
		functions, err := client.List(ctx, opts)
		if err != nil {
			return fmt.Errorf(msg.ErrorGetFunctions.Error(), err)
		}

		tbl := table.New("ID", "NAME", "LANGUAGE", "ACTIVE")
		tbl.WithWriter(f.IOStreams.Out)

		if opts.Details {
			tbl = table.New("ID", "NAME", "LANGUAGE", "ACTIVE", "LAST EDITOR", "MODIFIED", "REFERENCE COUNT", "INITIATOR_TYPE")
		}

		headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgGreen).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, v := range functions.Results {
			tbl.AddRow(v.GetId(), v.GetName(), v.GetLanguage(), v.GetActive(), v.GetLastEditor(), v.GetModified(), v.GetReferenceCount(), v.GetInitiatorType())
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

		if opts.Page >= *functions.TotalPages {
			break
		}

		if cmd.Flags().Changed("page") || cmd.Flags().Changed("page-size") {
			break
		}

		opts.Page++
	}

	return nil
}
