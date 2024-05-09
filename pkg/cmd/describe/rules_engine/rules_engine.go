package rulesengine

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/rules_engine"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
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
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223 --phase request
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223 --phase response --format json
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223 --phase request --out "./tmp/test.json"
    `),
		RunE: func(cmd *cobra.Command, _ []string) error {
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

			fields := make(map[string]string, 0)
			fields["Id"] = "Rules Engine ID"
			fields["Name"] = "Name"
			fields["Description"] = "Description"
			fields["Order"] = "Order"
			fields["IsActive"] = "Active"
			fields["Behaviors"] = "Behaviours"
			fields["Criteria"] = "Criteria"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Msg:         filepath.Clean(opts.OutPath),
					FlagOutPath: f.Out,
					FlagFormat:  f.Format,
					Out:         f.IOStreams.Out,
				},
				Fields: fields,
				Values: rules,
			}
			return output.Print(&describeOut)
		},
	}

	cmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.FlagAppID)
	cmd.Flags().Int64Var(&ruleID, "rule-id", 0, msg.FlagRuleID)
	cmd.Flags().StringVar(&phase, "phase", "request", msg.FlagPhase)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}
