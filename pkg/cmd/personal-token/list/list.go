package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	"github.com/aziontech/azion-cli/messages/general"
	msg "github.com/aziontech/azion-cli/messages/personal-token"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var details bool

	cmd := &cobra.Command{
		Use:           msg.ListUsage,
		Short:         msg.ListShortDescription,
		Long:          msg.ListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
        $ azion personal_token list 
        $ azion personal_token list --details
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient,
				f.Config.GetString("api_url"),
				f.Config.GetString("token"),
			)

			if err := PrintTable(client, f, details); err != nil {
				return fmt.Errorf(msg.ErrorGet.Error(), err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&details, "details", false, general.ApiListFlagDetails)
	flags.BoolP("help", "h", false, msg.ListHelpFlag)
	return cmd
}

func PrintTable(client *api.Client, f *cmdutil.Factory, details bool) error {
	c := context.Background()

	resp, err := client.List(c)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "NAME", "DESCRIPTION")
	tbl.WithWriter(f.IOStreams.Out)

	if details {
		tbl = table.New("ID", "NAME", "CREATED", "EXPIRES AT", "DESCRIPTION")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range resp {
		tbl.AddRow(
			*v.Uuid,
			utils.TruncateString(*v.Name),
			*v.Created,
			*v.ExpiresAt,
			utils.TruncateString(*v.Description.Get()))
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	tbl.PrintHeader(format)

	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	return nil
}
