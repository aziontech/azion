package rules_engine

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/rules_engine"
	api "github.com/aziontech/azion-cli/pkg/api/rules_engine"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID string
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
        $ azion create rules-engine --application-id 1679423488 --phase "response" --file ./file.json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				fields.ApplicationID = answer
			}

			if !cmd.Flags().Changed("file") {
				answer, err := utils.AskInput(msg.AskInputPathFile)
				if err != nil {
					return err
				}

				fields.Path = answer
			}

			request := api.CreateRulesEngineRequest{}

			err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
			if err != nil {
				logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
				return utils.ErrorUnmarshalReader
			}

			if request.Phase == "" {
				if !cmd.Flags().Changed("phase") {
					answer, err := utils.AskInput(msg.AskInputPhase)
					if err != nil {
						return err
					}

					fields.Phase = answer
				}
				request.SetPhase(fields.Phase)
			}

			if err := validateRequest(request.EdgeApplicationRuleEngineRequest); err != nil {
				return err
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			response, err := client.Create(context.Background(), fields.ApplicationID, request.EdgeApplicationRuleEngineRequest)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateRulesEngine.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out: f.IOStreams.Out,
			}
			return output.Print(&creatOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.ApplicationID, "application-id", "", msg.FlagEdgeApplicationID)
	flags.StringVar(&fields.Phase, "phase", "", msg.FlagPhase)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func validateRequest(request sdk.EdgeApplicationRuleEngineRequest) error {
	if request.GetName() == "" {
		return msg.ErrorNameEmpty
	}

	if request.GetCriteria() == nil {
		return msg.ErrorStructCriteriaNil
	}

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
