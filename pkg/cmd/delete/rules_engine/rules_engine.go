package rulesengine

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/delete/rules_engine"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var rule_id int64
	var app_id int64
	var phase string
	cmd := &cobra.Command{
		Use:           rulesengine.Usage,
		Short:         rulesengine.ShortDescription,
		Long:          rulesengine.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete rules-engine --rule-id 1234 --application-id 99887766 --phase request
		$ azion delete rules-engine
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("rule-id") {

				answer, err := utils.AskInput(rulesengine.AskInputRulesId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return rulesengine.ErrorConvertIdRule
				}

				rule_id = num
			}

			if !cmd.Flags().Changed("application-id") {

				answer, err := utils.AskInput(rulesengine.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return rulesengine.ErrorConvertIdRule
				}

				app_id = num
			}

			if !cmd.Flags().Changed("phase") {

				answer, err := utils.AskInput(rulesengine.AskInputPhase)
				if err != nil {
					return err
				}

				phase = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, app_id, phase, rule_id)
			if err != nil {
				return fmt.Errorf(rulesengine.ErrorFailToDelete.Error(), err)
			}

			out := f.IOStreams.Out
			fmt.Fprintf(out, rulesengine.DeleteOutputSuccess, rule_id)

			return nil
		},
	}

	cmd.Flags().Int64Var(&rule_id, "rule-id", 0, rulesengine.FlagRuleID)
	cmd.Flags().Int64Var(&app_id, "application-id", 0, rulesengine.FlagAppID)
	cmd.Flags().StringVar(&phase, "phase", "", rulesengine.FlagPhase)
	cmd.Flags().BoolP("help", "h", false, rulesengine.HelpFlag)

	return cmd
}
