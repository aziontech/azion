package rulesengine

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/rules_engine"

	api "github.com/aziontech/azion-cli/pkg/api/applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DescribeCmd struct {
	Io                     *iostreams.IOStreams
	ReadInput              func(string) (string, error)
	GetRulesEngineRequest  func(context.Context, int64, int64) (api.RulesEngineResponse, error)
	GetRulesEngineResponse func(context.Context, int64, int64) (api.RulesEngineResponse, error)
	AskInput               func(string) (string, error)
	ApplicationID          int64
	Phase                  string
	RuleID                 int64
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetRulesEngineRequest: func(ctx context.Context, appID, ruleID int64) (api.RulesEngineResponse, error) {
			client := registry.Applications(f)
			return client.GetRulesEngineRequest(ctx, appID, ruleID)
		},
		GetRulesEngineResponse: func(ctx context.Context, appID, ruleID int64) (api.RulesEngineResponse, error) {
			client := registry.Applications(f)
			return client.GetRulesEngineResponse(ctx, appID, ruleID)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223 --format json
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223 --out "./tmp/test.json"
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("rule-id") {
				answer, err := describe.AskInput(msg.AskInputRulesId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdRule
				}

				describe.RuleID = num
			}

			if !cmd.Flags().Changed("application-id") {
				answer, err := describe.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				describe.ApplicationID = num
			}

			if !cmd.Flags().Changed("phase") {
				answer, err := describe.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				describe.Phase = answer
			}

			ctx := context.Background()
			fields := make(map[string]string)
			fields["Id"] = "Rules Engine ID"
			fields["Name"] = "Name"
			fields["Description"] = "Description"
			fields["Order"] = "Order"
			fields["Active"] = "Active"
			fields["Behaviors"] = "Behaviours"
			fields["Criteria"] = "Criteria"

			switch describe.Phase {
			case "request":
				rules, err := describe.GetRulesEngineRequest(ctx, describe.ApplicationID, describe.RuleID)
				if err != nil {
					return fmt.Errorf(msg.ErrorGetRulesEngine.Error(), err)
				}

				describeOut := output.DescribeOutput{
					GeneralOutput: output.GeneralOutput{
						Flags: f.Flags,
						Out:   f.IOStreams.Out,
					},
					Fields: fields,
					Values: rules,
				}
				return output.Print(&describeOut)
			case "response":
				rules, err := describe.GetRulesEngineResponse(ctx, describe.ApplicationID, describe.RuleID)
				if err != nil {
					return fmt.Errorf(msg.ErrorGetRulesEngine.Error(), err)
				}

				describeOut := output.DescribeOutput{
					GeneralOutput: output.GeneralOutput{
						Flags: f.Flags,
						Out:   f.IOStreams.Out,
					},
					Fields: fields,
					Values: rules,
				}
				return output.Print(&describeOut)
			default:
				return msg.ErrorInvalidPhase
			}
		},
	}

	cobraCmd.Flags().Int64Var(&describe.ApplicationID, "application-id", 0, msg.FlagAppID)
	cobraCmd.Flags().Int64Var(&describe.RuleID, "rule-id", 0, msg.FlagRuleID)
	cobraCmd.Flags().StringVar(&describe.Phase, "phase", "", msg.FlagPhase)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
