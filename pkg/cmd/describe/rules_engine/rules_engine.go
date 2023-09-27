package rulesengine

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/describe/rules_engine"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID int64
	ruleID        int64
	phase         string
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
      $ azion describe rule-engine  --application-id 1673635839 --rule-id 31223 --phase request
      $ azion describe rule-engine --application-id 1673635839 --rule-id 31223 --phase response --format json
      $ azion describe rule-engine --application-id 1673635839 --rule-id 31223 --phase request --out "./tmp/test.json"
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("rule-id") {

				answer, err := utils.AskInput(msg.AskInputRulesId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdRule
				}

				ruleID = num
			}

			if !cmd.Flags().Changed("application-id") {

				answer, err := utils.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdRule
				}

				applicationID = num
			}

			if !cmd.Flags().Changed("phase") {

				answer, err := utils.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				phase = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			rules, err := client.GetRulesEngine(ctx, applicationID, ruleID, phase)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetRulesEngine.Error(), err)
			}

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, rules)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, msg.FileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedFuction[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.FlagAppID)
	cmd.Flags().Int64Var(&ruleID, "rule-id", 0, msg.FlagRuleID)
	cmd.Flags().StringVar(&phase, "phase", "request", msg.FlagPhase)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.DescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.DescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}

func format(cmd *cobra.Command, rules api.RulesEngineResponse) ([]byte, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(rules, "", " ")
	}

	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())
	tbl.AddRow("Rules Engine ID: ", rules.GetId())
	tbl.AddRow("Name: ", rules.GetName())
	tbl.AddRow("Description: ", rules.GetDescription())
	tbl.AddRow("Order: ", rules.GetOrder())
	tbl.AddRow("Active: ", rules.GetIsActive())
	tbl.AddRow("")
	tbl.AddRow("Behaviours: ")
	for _, b := range rules.GetBehaviors() {
		tbl.AddRow("  Name: ", b.GetName())
		tbl.AddRow("  Target: ", b.GetTarget())
		tbl.AddRow("")
	}
	tbl.AddRow("Criteria: ")
	for _, c := range rules.GetCriteria() {
		for _, c2 := range c {
			tbl.AddRow("  Conditional: ", c2.GetConditional())
			tbl.AddRow("  Variable: ", c2.GetVariable())
			tbl.AddRow("  Operator: ", c2.GetOperator())
			tbl.AddRow("  Input Value: ", c2.GetInputValue())
			tbl.AddRow("")
		}
	}
	return tbl.GetByteFormat(), nil
}
