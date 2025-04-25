package rules_engine

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Fields struct {
	ApplicationID string
	RuleID        string
	Phase         string
	Path          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update rules-engine -h"
		$ azion update rules-engine --rule-id 1234 --application-id 1673635839 --file ruleengine.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateUserInput(cmd, fields); err != nil {
				return err
			}

			request := api.UpdateRulesEngineRequest{}

			err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
			if err != nil {
				logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
				return utils.ErrorUnmarshalReader
			}

			if err := validateRequest(request); err != nil {
				return err
			}

			request.ApplicationID = fields.ApplicationID
			request.RulesID = fields.RuleID

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			response, err := client.Update(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdate.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.ApplicationID, "application-id", "", msg.FlagApplicationID)
	flags.StringVar(&fields.RuleID, "rule-id", "", msg.FlagRulesEngineID)
	flags.StringVar(&fields.Phase, "phase", "", msg.RulesEnginePhase)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.FlagHelp)
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

				if !item.Argument.IsSet() {
					return msg.ErrorInputValueEmpty
				}
			}
		}
	}

	if request.GetBehaviors() == nil {
		return msg.ErrorStructBehaviorsNil
	}

	for _, item := range request.GetBehaviors() {
		if item.Name == "" {
			return msg.ErrorNameBehaviorsEmpty
		}

		if !item.Argument.IsSet() {
			return msg.ErrorArgumentBehaviorsEmpty
		}
	}

	return nil
}

func validateUserInput(cmd *cobra.Command, fields *Fields) error {
	if !cmd.Flags().Changed("application-id") {
		answer, err := utils.AskInput(msg.AskInputApplicationID)
		if err != nil {
			return err
		}

		fields.ApplicationID = answer
	}

	if !cmd.Flags().Changed("rule-id") {
		answer, err := utils.AskInput(msg.AskInputRulesID)
		if err != nil {
			return err
		}

		fields.RuleID = answer
	}

	if !cmd.Flags().Changed("file") {
		answer, err := utils.AskInput(msg.AskInputPathFile)
		if err != nil {
			return err
		}

		fields.Path = answer
	}

	return nil
}
