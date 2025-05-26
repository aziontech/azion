package rulesengine

import (
	"context"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/rules_engine"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID int64
	ruleID        int64
	phase         string
)

type DescribeCmd struct {
	Io             *iostreams.IOStreams
	ReadInput      func(string) (string, error)
	GetRulesEngine func(context.Context, int64, int64, string) (api.RulesEngineResponse, error)
	AskInput       func(string) (string, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetRulesEngine: func(ctx context.Context, appID, ruleID int64, phase string) (api.RulesEngineResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.GetRulesEngine(ctx, appID, ruleID, phase)
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
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223 --phase request
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223 --phase response --format json
      $ azion describe rules-engine --application-id 1673635839 --rule-id 31223 --phase request --out "./tmp/test.json"
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

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

				ruleID = num
			}

			if !cmd.Flags().Changed("application-id") {
				answer, err := describe.AskInput(msg.AskInputApplicationId)
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
				answer, err := describe.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				phase = answer
			}

			ctx := context.Background()
			rules, err := describe.GetRulesEngine(ctx, applicationID, ruleID, phase)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetRulesEngine.Error(), err)
			}

			fields := make(map[string]string)
			fields["Id"] = "Rules Engine ID"
			fields["Name"] = "Name"
			fields["Description"] = "Description"
			fields["Order"] = "Order"
			fields["IsActive"] = "Active"
			fields["Behaviors"] = "Behaviours"
			fields["Criteria"] = "Criteria"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Flags: f.Flags,
					Out:   f.IOStreams.Out,
				},
				Fields: fields,
				Values: rules,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.FlagAppID)
	cobraCmd.Flags().Int64Var(&ruleID, "rule-id", 0, msg.FlagRuleID)
	cobraCmd.Flags().StringVar(&phase, "phase", "request", msg.FlagPhase)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
