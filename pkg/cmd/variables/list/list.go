package list

import (
	"context"
	"strings"

	"github.com/aziontech/azion-cli/utils"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	"github.com/aziontech/azion-cli/messages/general"
	msg "github.com/aziontech/azion-cli/messages/variables"
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	listCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.VariablesListShortDescription,
		Long:          msg.VariablesListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion variables list --details
		$ azion variables list
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			if err := listAllVariables(client, f, opts); err != nil {
				return err
			}
			return nil
		},
	}

	listCmd.Flags().BoolVar(&opts.Details, "details", false, general.ApiListFlagDetails)
	listCmd.Flags().BoolP("help", "h", false, msg.VariablesListHelpFlag)
	return listCmd
}

func listAllVariables(client *api.Client, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	c := context.Background()

	resp, err := client.List(c)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "KEY", "VALUE")
	tbl.WithWriter(f.IOStreams.Out)

	if opts.Details {
		tbl = table.New("ID", "KEY", "VALUE", "SECRET", "LAST EDITOR")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range resp {
		tbl.AddRow(v.GetUuid(), v.GetKey(), utils.TruncateString(v.GetValue()), v.GetSecret(), v.GetLastEditor())
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	tbl.PrintHeader(format)

	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	return nil
}
