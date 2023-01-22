package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/domains"
	api "github.com/aziontech/azion-cli/pkg/api/domains"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.DomainsListUsage,
		Short:         msg.DomainsListShortDescription,
		Long:          msg.DomainsListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		$ azioncli domains list
		`),

		// $ azioncli domains list --details
		// $ azioncli domains list --order_by "id"
		// $ azioncli domains list --page 1
		// $ azioncli domains list --page_size 5
		// $ azioncli domains list --sort "asc"

		RunE: func(cmd *cobra.Command, args []string) error {
			var numberPage int64 = opts.Page
			if err := PrintTable(cmd, f, opts, &numberPage); err != nil {
				return fmt.Errorf(msg.ErrorGetDomains.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.DomainsListHelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, numberPage *int64) error {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	domains, err := client.List(ctx, opts)
	if err != nil {
		return fmt.Errorf(msg.ErrorGetDomains.Error(), err)
	}

	tbl := table.New("ID", "NAME")
	table.DefaultWriter = f.IOStreams.Out
	if cmd.Flags().Changed("details") {
		tbl = table.New("ID", "NAME", "ACTIVE")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range domains.Results {
		tbl.AddRow(v.Id, v.Name)
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
	return nil
}
