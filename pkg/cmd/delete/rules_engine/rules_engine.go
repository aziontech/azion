package rulesengine

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	ruleID        string
	applicationID string
	phase         string
)

type DeleteCmd struct {
	Io                 *iostreams.IOStreams
	ReadInput          func(string) (string, error)
	DeleteRuleRequest  func(context.Context, string, string) error
	DeleteRuleResponse func(context.Context, string, string) error
	AskInput           func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteRuleRequest: func(ctx context.Context, ruleID, appID string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.DeleteRequest(ctx, appID, ruleID)
		},
		DeleteRuleResponse: func(ctx context.Context, ruleID, appID string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.DeleteResponse(ctx, appID, ruleID)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(delete *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion delete rules-engine --rule-id 1234 --application-id 99887766
			$ azion delete rules-engine
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("rule-id") {
				answer, err := delete.AskInput(msg.AskInputRulesId)
				if err != nil {
					return err
				}

				ruleID = answer
			}

			if !cmd.Flags().Changed("application-id") {
				answer, err := delete.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				applicationID = answer
			}

			if !cmd.Flags().Changed("phase") {
				answer, err := delete.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				phase = answer
			}

			ctx := context.Background()

			switch phase {
			case "request":
				err = delete.DeleteRuleRequest(ctx, ruleID, applicationID)
				if err != nil {
					return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
				}
			case "response":
				err = delete.DeleteRuleResponse(ctx, ruleID, applicationID)
				if err != nil {
					return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
				}
			default:
				return msg.ErrorInvalidPhase

			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeleteOutputSuccess, ruleID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().StringVar(&ruleID, "rule-id", "", msg.FlagRuleID)
	cobraCmd.Flags().StringVar(&applicationID, "application-id", "", msg.FlagAppID)
	cobraCmd.Flags().StringVar(&phase, "phase", "request", msg.FlagPhase)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
