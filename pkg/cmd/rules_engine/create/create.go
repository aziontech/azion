package create

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
	Phase         string
	Path          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.RulesEngineCreateUsage,
		Short:         msg.RulesEngineCreateShortDescription,
		Long:          msg.RulesEngineCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("phase") || !cmd.Flags().Changed("in") {
				return msg.ErrorMandatoryCreateFlags
			}

			request := api.CreateRulesEngineRequest{}
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

			if err := validRequest(request); err != nil {
				return err
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.CreateRulesEngine(context.Background(), fields.ApplicationID, fields.Phase, &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateRulesEngine.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.RulesEngineCreateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.RulesEngineCreateFlagEdgeApplicationID)
    flags.StringVarP(&fields.Phase, "phase", "p", "", msg.RulesEngineCreateFlagPhase)
	flags.StringVar(&fields.Path, "in", "", msg.RulesEngineCreateFlagIn)
	flags.BoolP("help", "h", false, msg.RulesEngineCreateHelpFlag)
	return cmd
}

func validRequest(request api.CreateRulesEngineRequest) error {
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

			if item.InputValue == nil {
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
	}

	return nil
}
