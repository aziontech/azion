package rulesengine

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeleteCmd struct {
	Io                 *iostreams.IOStreams
	ReadInput          func(string) (string, error)
	DeleteRuleRequest  func(context.Context, int64, int64) error
	DeleteRuleResponse func(context.Context, int64, int64) error
	AskInput           func(string) (string, error)
	ApplicationID      int64
	Phase              string
	RuleID             int64
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteRuleRequest: func(ctx context.Context, ruleID, appID int64) error {
			client := registry.RulesEngine(f)
			return client.DeleteRequest(ctx, appID, ruleID)
		},
		DeleteRuleResponse: func(ctx context.Context, ruleID, appID int64) error {
			client := registry.RulesEngine(f)
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

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdRule
				}

				delete.RuleID = num
			}

			if !cmd.Flags().Changed("application-id") {
				answer, err := delete.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				delete.ApplicationID = num
			}

			if !cmd.Flags().Changed("phase") {
				answer, err := delete.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}

				delete.Phase = answer
			}

			ctx := context.Background()

			switch delete.Phase {
			case "request":
				err = delete.DeleteRuleRequest(ctx, delete.RuleID, delete.ApplicationID)
				if err != nil {
					return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
				}
			case "response":
				err = delete.DeleteRuleResponse(ctx, delete.RuleID, delete.ApplicationID)
				if err != nil {
					return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
				}
			default:
				return msg.ErrorInvalidPhase

			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeleteOutputSuccess, delete.RuleID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&delete.RuleID, "rule-id", 0, msg.FlagRuleID)
	cobraCmd.Flags().Int64Var(&delete.ApplicationID, "application-id", 0, msg.FlagAppID)
	cobraCmd.Flags().StringVar(&delete.Phase, "phase", "request", msg.FlagPhase)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
