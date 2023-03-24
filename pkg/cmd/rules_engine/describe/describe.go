package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/rules_engine"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
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
		Use:           msg.RulesEngineDescribeUsage,
		Short:         msg.RulesEngineDescribeShortDescription,
		Long:          msg.RulesEngineDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
      $ azioncli rules_engine describe --application-id 4312 --rule-id 31223 --phase request
      $ azioncli rules_engine describe --application-id 1337 --rule-id 31223 --phase response --format json
      $ azioncli rules_engine describe --application-id 1337 --rule-id 31223 --phase request --out "./tmp/test.json" --format json
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("phase") || !cmd.Flags().Changed("rule-id") {
				return msg.ErrorMandatoryFlags
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
				fmt.Fprintf(out, msg.RulesEngineFileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedFuction[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, msg.ApplicationFlagId)
	cmd.Flags().Int64VarP(&ruleID, "rule-id", "r", 0, msg.RulesEngineFlagId)
	cmd.Flags().StringVarP(&phase, "phase", "p", "request", msg.RulesEngineListHelpPhase)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.RulesEngineDescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.RulesEngineDescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineDescribeHelpFlag)

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
