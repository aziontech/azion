package rulesengine

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/rules_engine"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID string
	ruleID        string
)

type DescribeCmd struct {
	Io             *iostreams.IOStreams
	ReadInput      func(string) (string, error)
	GetRulesEngine func(context.Context, string, string) (api.RulesEngineResponse, error)
	AskInput       func(string) (string, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		GetRulesEngine: func(ctx context.Context, appID, ruleID string) (api.RulesEngineResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.GetRulesEngine(ctx, appID, ruleID)
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
			var err error

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
			rules, err := describe.GetRulesEngine(ctx, applicationID, ruleID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetRulesEngine.Error(), err)
			}

			fields := make(map[string]string)
			fields["Id"] = "Rules Engine ID"
			fields["Name"] = "Name"
			fields["Description"] = "Description"
			fields["Order"] = "Order"
			fields["Active"] = "Active"
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

	cobraCmd.Flags().StringVar(&applicationID, "application-id", "", msg.FlagAppID)
	cobraCmd.Flags().StringVar(&ruleID, "rule-id", "", msg.FlagRuleID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
