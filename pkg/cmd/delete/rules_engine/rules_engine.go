package rulesengine

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	ruleID        int64
	applicationID int64
	phase         string
)

type DeleteCmd struct {
	Io         *iostreams.IOStreams
	ReadInput  func(string) (string, error)
	DeleteRule func(context.Context, int64, int64, string) error
	AskInput   func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteRule: func(ctx context.Context, ruleID, appID int64, phase string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Delete(ctx, appID, phase, ruleID)
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
			$ azion delete rules-engine --rule-id 1234 --application-id 99887766 --phase request
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
				ruleID = num
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
				applicationID = num
			}

			if !cmd.Flags().Changed("phase") {
				answer, err := delete.AskInput(msg.AskInputPhase)
				if err != nil {
					return err
				}
				phase = answer
			}

			ctx := context.Background()

			err = delete.DeleteRule(ctx, ruleID, applicationID, phase)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeleteOutputSuccess, ruleID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&ruleID, "rule-id", 0, msg.FlagRuleID)
	cobraCmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.FlagAppID)
	cobraCmd.Flags().StringVar(&phase, "phase", "", msg.FlagPhase)
	cobraCmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
