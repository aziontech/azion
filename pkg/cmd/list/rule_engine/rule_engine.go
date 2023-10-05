package ruleengine

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	table "github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/list/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
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
		$ azion list rules-engine --application-id 1673635839 --phase request
		$ azion list rules-engine --application-id 1673635839 --phase response --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("application-id") {

				answer, err := utils.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				edgeApplicationID = num
			}

			if !cmd.Flags().Changed("phase") {

				answer, err := utils.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				phase = answer
			}

			if err := PrintTable(cmd, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetRulesEngines.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineListHelpFlag)
	cmd.Flags().Int64Var(&edgeApplicationID, "application-id", 0, msg.ApplicationFlagId)
	cmd.Flags().StringVar(&phase, "phase", "request", msg.RulesEnginePhase)
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

	for _, v := range rules.Results {
		tbl.AddRow(v.Id, v.Name, v.Order, v.Phase, v.IsActive)
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	logger.PrintHeader(tbl, format)
	for _, row := range tbl.GetRows() {
		logger.PrintRow(tbl, format, row)
	}

	f.IOStreams.Out = table.DefaultWriter
	return nil
}
