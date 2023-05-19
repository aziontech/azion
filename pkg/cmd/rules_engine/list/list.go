package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	table "github.com/maxwelbm/tablecli"
	"github.com/spf13/cobra"
)

var edgeApplicationID int64
var phase string

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.RulesEngineListUsage,
		Short:         msg.RulesEngineListShortDescription,
		Long:          msg.RulesEngineListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
		$ azioncli rules_engine list -a 1673635839 -p request
		$ azioncli rules_engine list --application-id 1673635839 --phase response --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("phase") {
				return msg.ErrorMandatoryListFlags
			}
			if err := PrintTable(cmd, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetRulesEngines.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineListHelpFlag)
	cmd.Flags().Int64VarP(&edgeApplicationID, "application-id", "a", 0, msg.ApplicationFlagId)
	cmd.Flags().StringVarP(&phase, "phase", "p", "request", msg.RulesEnginePhase)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	rules, err := client.ListRulesEngine(ctx, opts, edgeApplicationID, phase)
	if err != nil {
		return err
	}

	tbl := table.New("ID", "NAME")
	table.DefaultWriter = f.IOStreams.Out
	if cmd.Flags().Changed("details") {
		tbl = table.New("ID", "NAME", "ORDER", "PHASE", "ACTIVE")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	if cmd.Flags().Changed("details") {
		for _, v := range rules.Results {
			tbl.AddRow(v.Id, v.Name, v.Order, v.Phase, v.IsActive)
		}
	} else {
		for _, v := range rules.Results {
			tbl.AddRow(v.Id, v.Name)
		}
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	tbl.PrintHeader(format)
	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	f.IOStreams.Out = table.DefaultWriter
	return nil
}
