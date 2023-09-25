package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var applicationID, ruleID int64
	var phase string

	cmd := &cobra.Command{
		Use:           msg.RulesEngineDeleteUsage,
		Short:         msg.RulesEngineDeleteShortDescription,
		Long:          msg.RulesEngineDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azion rules_engine delete --application-id 1673635839 --rule-id 12312
		  $ azion rules_engine delete -a 1673635839 -r 12312
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("phase") || !cmd.Flags().Changed("rule-id") {
				return msg.ErrorMissingArgumentsDelete
			}
			if err := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token")).
				DeleteRulesEngine(context.Background(), applicationID, phase, ruleID); err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.RulesEngineDeleteOutputSuccess, ruleID)
			return nil
		},
	}

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, msg.RulesEngineDeleteFlagApplicationID)
	cmd.Flags().Int64VarP(&ruleID, "rule-id", "r", 0, msg.RulesEngineDeleteFlagRuleID)
	cmd.Flags().StringVar(&phase, "phase", "", msg.RulesEngineDeleteFlagPhase)
	cmd.Flags().BoolP("help", "h", false, msg.RulesEngineDeleteHelpFlag)
	return cmd
}
