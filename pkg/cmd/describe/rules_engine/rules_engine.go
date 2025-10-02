package rulesengine

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/rules_engine"

	api "github.com/aziontech/azion-cli/pkg/api/applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID string
	ruleID        string
	phase         string
)

type DescribeCmd struct {
	Io                     *iostreams.IOStreams
	ReadInput              func(string) (string, error)
	GetRulesEngineRequest  func(context.Context, string, string) (api.RulesEngineResponse, error)
	GetRulesEngineResponse func(context.Context, string, string) (api.RulesEngineResponse, error)
	AskInput               func(string) (string, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetRulesEngineRequest: func(ctx context.Context, appID, ruleID string) (api.RulesEngineResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.GetRulesEngineRequest(ctx, appID, ruleID)
		},
		GetRulesEngineResponse: func(ctx context.Context, appID, ruleID string) (api.RulesEngineResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
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

				ruleID = answer
			}

			if !cmd.Flags().Changed("application-id") {
				answer, err := describe.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				applicationID = answer
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

			switch phase {
			case "request":
				rules, err := describe.GetRulesEngineRequest(ctx, applicationID, ruleID)
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
				rules, err := describe.GetRulesEngineResponse(ctx, applicationID, ruleID)
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

	cobraCmd.Flags().StringVar(&applicationID, "application-id", "", msg.FlagAppID)
	cobraCmd.Flags().StringVar(&ruleID, "rule-id", "", msg.FlagRuleID)
	cobraCmd.Flags().StringVar(&phase, "phase", "", msg.FlagPhase)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
