package update

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID int64
	RuleID        int64
	Phase         string
	Path          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.RulesEngineUpdateUsage,
		Short:         msg.RulesEngineUpdateShortDescription,
		Long:          msg.RulesEngineUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli rules_engine update --rule-id 1234 --application-id 1673635839 --phase request --in ruleengine.json"
		$ azioncli rules_engine update --application-id 1673635839 --rule-id 1234 --phase request --in ruleengine.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("phase") || !cmd.Flags().Changed("rule-id") || !cmd.Flags().Changed("in") {
				return msg.ErrorMandatoryFlagsUpdate
			}

			request := api.UpdateRulesEngineRequest{}
			var (
				file *os.File
				err  error
			)
			if fields.Path == "-" {
				file = os.Stdin
			} else {
				file, err = os.Open(fields.Path)
				if err != nil {
					return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.Path)
				}
			}
			err = cmdutil.UnmarshallJsonFromReader(file, &request)
			if err != nil {
				return utils.ErrorUnmarshalReader
			}

			if err := validateRequest(request); err != nil {
				return err
			}
			request.IdApplication = fields.ApplicationID
			request.Phase = fields.Phase
			request.Id = fields.RuleID

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.UpdateRulesEngine(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateRulesengine.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.RulesEngineUpdateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.RulesEngineUpdateFlagEdgeApplicationId)
	flags.Int64VarP(&fields.RuleID, "rule-id", "r", 0, msg.RulesEngineFlagId)
	flags.StringVarP(&fields.Phase, "phase", "p", "", msg.RulesEnginePhase)
	flags.StringVar(&fields.Path, "in", "", msg.RulesEngineUpdateFlagIn)
	flags.BoolP("help", "h", false, msg.RulesEngineUpdateHelpFlag)
	return cmd
}

func validateRequest(request api.UpdateRulesEngineRequest) error {

	if request.GetCriteria() != nil {
		for _, itemCriteria := range request.GetCriteria() {
			for _, item := range itemCriteria {
				if item.Conditional == "" {
					return msg.ErrorConditionalEmpty
				}

				if item.Variable == "" {
					return msg.ErrorVariableEmpty
				}

				if item.Operator == "" {
					return msg.ErrorOperatorEmpty
				}

				if item.InputValue == nil {
					return msg.ErrorInputValueEmpty
				}
			}
		}
	}

	if request.GetBehaviors() != nil {
		for _, item := range request.GetBehaviors() {
			if item.Name == "" {
				return msg.ErrorNameBehaviorsEmpty
			}
		}
	}

	return nil
}
