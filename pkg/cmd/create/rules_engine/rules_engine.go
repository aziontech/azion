package rules_engine

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/create/rules_engine"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/logger"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID int64
	Phase         string
	Path          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           rules_engine.Usage,
		Short:         rules_engine.ShortDescription,
		Long:          rules_engine.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create rules-engine --application-id 1679423488 --phase "response" --in ./file.json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(rules_engine.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return rules_engine.ErrorConvertIdApplication
				}

				fields.ApplicationID = num
			}

			if !cmd.Flags().Changed("phase") {
				answer, err := utils.AskInput(rules_engine.AskInputPhase)
				if err != nil {
					return err
				}

				fields.Phase = answer
			}

			if !cmd.Flags().Changed("in") {
				answer, err := utils.AskInput(rules_engine.AskInputPathFile)
				if err != nil {
					return err
				}

				fields.Path = answer
			}

			request := api.CreateRulesEngineRequest{}

			err := utils.FlagINUnmarshalFileJSON(fields.Path, &request)
			if err != nil {
				logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
				return utils.ErrorUnmarshalReader
			}

			reqSdk := dtoStructRequest(request.CreateRulesEngineRequest)

			if err := validateRequest(reqSdk); err != nil {
				return err
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), fields.ApplicationID, fields.Phase, reqSdk)

			if err != nil {
				return fmt.Errorf(rules_engine.ErrorCreateRulesEngine.Error(), err)
			}

			logger.FInfo(f.IOStreams.Out, fmt.Sprintf(rules_engine.OutputSuccess, response.GetId()))
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, rules_engine.FlagEdgeApplicationID)
	flags.StringVar(&fields.Phase, "phase", "", rules_engine.FlagPhase)
	flags.StringVar(&fields.Path, "in", "", rules_engine.FlagIn)
	flags.BoolP("help", "h", false, rules_engine.HelpFlag)
	return cmd
}

func validateRequest(request sdk.CreateRulesEngineRequest) error {
	if request.GetName() == "" {
		return rules_engine.ErrorNameEmpty
	}

	if request.GetCriteria() == nil {
		return rules_engine.ErrorStructCriteriaNil
	}

	for _, itemCriteria := range request.GetCriteria() {
		for _, item := range itemCriteria {
			if item.Conditional == "" {
				return rules_engine.ErrorConditionalEmpty
			}

			if item.Variable == "" {
				return rules_engine.ErrorVariableEmpty
			}

			if item.Operator == "" {
				return rules_engine.ErrorOperatorEmpty
			}

			if item.InputValue == nil {
				return rules_engine.ErrorInputValueEmpty
			}
		}
	}

	if request.GetBehaviors() == nil {
		return rules_engine.ErrorStructBehaviorsNil
	}

	for _, item := range request.GetBehaviors() {
		if item.RulesEngineBehaviorString != nil {
			if item.RulesEngineBehaviorString.Name == "" {
				return rules_engine.ErrorNameBehaviorsEmpty

			}
		}
		if item.RulesEngineBehaviorObject != nil && (item.RulesEngineBehaviorObject.Target.CapturedArray == nil || item.RulesEngineBehaviorObject.Target.Regex == nil || item.RulesEngineBehaviorObject.Target.Subject == nil) {
			if item.RulesEngineBehaviorObject.Name == "" {
				return rules_engine.ErrorNameBehaviorsEmpty
			}
		}
	}

	return nil
}

func dtoStructRequest(request sdk.CreateRulesEngineRequest) sdk.CreateRulesEngineRequest {
	var req sdk.CreateRulesEngineRequest

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
