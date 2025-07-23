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

			if fields.Phase == "" {
				if !cmd.Flags().Changed("phase") {
					answer, err := utils.AskInput(msg.AskInputPhase)
					if err != nil {
						return err
					}

					fields.Phase = answer
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			var id int64
			switch fields.Phase {
			case "request":
				request := api.CreateRulesEngineRequest{}

				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
				response, err := client.CreateRequest(context.Background(), fields.ApplicationID, request.EdgeApplicationRequestPhaseRuleEngineRequest)
				if err != nil {
					return fmt.Errorf(msg.ErrorCreateRulesEngine.Error(), err)
				}
				id = response.GetId()
			case "response":
				request := api.CreateRulesEngineResponse{}

				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
				response, err := client.CreateResponse(context.Background(), fields.ApplicationID, request.EdgeApplicationResponsePhaseRuleEngineRequest)
				if err != nil {
					return fmt.Errorf(msg.ErrorCreateRulesEngine.Error(), err)
				}
				id = response.GetId()
			default:
			}

			creatOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.OutputSuccess, id),
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
