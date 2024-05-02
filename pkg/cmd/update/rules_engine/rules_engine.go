package rules_engine

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
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
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update rules-engine -h"
		$ azion update rules-engine --rule-id 1234 --application-id 1673635839 --phase request --file ruleengine.json"
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

			reqSdk := dtoStructRequest(request)

			if err := validateRequest(reqSdk); err != nil {
				return err
			}

			reqSdk.ApplicationID = fields.ApplicationID
			reqSdk.RulesID = fields.RuleID
			reqSdk.Phase = fields.Phase

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Update(context.Background(), &reqSdk)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdate.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:         fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:         f.IOStreams.Out,
				FlagOutPath: f.Out,
				FlagFormat:  f.Format,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.FlagApplicationID)
	flags.Int64Var(&fields.RuleID, "rule-id", 0, msg.FlagRulesEngineID)
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

				if item.InputValue == nil {
					return msg.ErrorInputValueEmpty
				}
			}
		}
	}

	if request.GetBehaviors() != nil {
		for _, item := range request.GetBehaviors() {
			if item.RulesEngineBehaviorString != nil {
				if item.RulesEngineBehaviorString.Name == "" {
					return msg.ErrorNameBehaviorsEmpty

				}
			}
			if item.RulesEngineBehaviorObject != nil && (item.RulesEngineBehaviorObject.Target.CapturedArray == nil || item.RulesEngineBehaviorObject.Target.Regex == nil || item.RulesEngineBehaviorObject.Target.Subject == nil) {
				if item.RulesEngineBehaviorObject.Name == "" {
					return msg.ErrorNameBehaviorsEmpty
				}
			}
		}
	}

	return nil
}

func dtoStructRequest(request api.UpdateRulesEngineRequest) api.UpdateRulesEngineRequest {
	var req api.UpdateRulesEngineRequest

	req.Name = request.Name
	req.Description = request.Description

	var rulesEngineCriteria [][]sdk.RulesEngineCriteria
	for _, itemCriterias := range request.Criteria {
		var criterias []sdk.RulesEngineCriteria
		for _, itemCriteria := range itemCriterias {
			var criteria sdk.RulesEngineCriteria

			criteria.Conditional = itemCriteria.Conditional
			criteria.Variable = itemCriteria.Variable
			criteria.Operator = itemCriteria.Operator
			criteria.InputValue = itemCriteria.InputValue

			criterias = append(criterias, criteria)
		}
		rulesEngineCriteria = append(rulesEngineCriteria, criterias)
	}

	req.Criteria = rulesEngineCriteria
	var behaviors []sdk.RulesEngineBehaviorEntry
	for _, v := range request.Behaviors {
		if v.RulesEngineBehaviorObject != nil {
			if v.RulesEngineBehaviorObject.Target.CapturedArray != nil && v.RulesEngineBehaviorObject.Target.Regex != nil && v.RulesEngineBehaviorObject.Target.Subject != nil {
				var behaviorObject sdk.RulesEngineBehaviorObject
				behaviorObject.SetName(v.RulesEngineBehaviorObject.Name)
				behaviorObject.SetTarget(v.RulesEngineBehaviorObject.Target)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorObject: &behaviorObject,
				})
			} else {
				var behaviorString sdk.RulesEngineBehaviorString
				behaviorString.SetName(v.RulesEngineBehaviorObject.Name)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorString: &behaviorString,
				})
			}
		} else {
			if v.RulesEngineBehaviorString != nil {
				var behaviorString sdk.RulesEngineBehaviorString
				behaviorString.SetName(v.RulesEngineBehaviorString.Name)
				behaviorString.SetTarget(v.RulesEngineBehaviorString.Target)
				behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
					RulesEngineBehaviorString: &behaviorString,
				})
			}
		}
	}

	req.Behaviors = behaviors

	return req
}

func validateUserInput(cmd *cobra.Command, fields *Fields) error {
	if !cmd.Flags().Changed("application-id") {
		answer, err := utils.AskInput(msg.AskInputApplicationID)
		if err != nil {
			return err
		}

		num, err := strconv.ParseInt(answer, 10, 64)
		if err != nil {
			logger.Debug("Error while converting answer to int64", zap.Error(err))
			return msg.ErrorConvertApplicationID
		}

		fields.ApplicationID = num
	}

	if !cmd.Flags().Changed("rule-id") {
		answer, err := utils.AskInput(msg.AskInputRulesID)
		if err != nil {
			return err
		}

		num, err := strconv.ParseInt(answer, 10, 64)
		if err != nil {
			logger.Debug("Error while converting answer to int64", zap.Error(err))
			return msg.ErrorConvertRulesID
		}

		fields.RuleID = num
	}

	if !cmd.Flags().Changed("phase") {
		answer, err := utils.AskInput(msg.AskInputPhase)
		if err != nil {
			return err
		}

		fields.Phase = answer
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
