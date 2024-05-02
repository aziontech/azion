package rulesengine

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var rule_id int64
	var app_id int64
	var phase string
	cmd := &cobra.Command{
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

				rule_id = num
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

				app_id = num
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

			err := client.Delete(ctx, app_id, phase, rule_id)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:         fmt.Sprintf(msg.DeleteOutputSuccess, rule_id),
				Out:         f.IOStreams.Out,
				FlagOutPath: f.Out,
				FlagFormat:  f.Format,
			}
			return output.Print(&deleteOut)

		},
	}

	cmd.Flags().Int64Var(&rule_id, "rule-id", 0, msg.FlagRuleID)
	cmd.Flags().Int64Var(&app_id, "application-id", 0, msg.FlagAppID)
	cmd.Flags().StringVar(&phase, "phase", "", msg.FlagPhase)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}
