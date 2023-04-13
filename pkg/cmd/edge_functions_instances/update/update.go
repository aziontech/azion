package update

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/edge_functions_instances"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID string
	Name          string
	FunctionID    int64
	InstanceID    string
	Path          string
	Args          interface{}
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.EdgeFuncInstanceUpdateUsage,
		Short:         msg.EdgeFuncInstanceUpdateShortDescription,
		Long:          msg.EdgeFuncInstanceUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions_instances update -a 12345 -i 43121 -f 1209
        $ azioncli edge_functions_instances update --application-id 12 --instance-id 2121 --function-id 1212 --name updated
        $ azioncli edge_functions_instances update  -a 12345 -i 43121 --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("in") && (!cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("instance-id")) {
				return msg.ErrorMandatoryUpdateFlagsIn
			}
			request := api.UpdateInstanceRequest{}
			if cmd.Flags().Changed("in") {
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
			} else {
				if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("instance-id") || !cmd.Flags().Changed("function-id") {
					return msg.ErrorMandatoryUpdateFlags
				}
				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}
				if cmd.Flags().Changed("function-id") {
					request.SetEdgeFunctionId(fields.FunctionID)
				}
				if cmd.Flags().Changed("args") {
					request.SetArgs(fields.Args)
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.UpdateInstance(context.Background(), &request, fields.ApplicationID, fields.InstanceID)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateFuncInstance.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.EdgeFuncInstanceUpdateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&fields.ApplicationID, "application-id", "a", "0", msg.EdgeFuncInstanceUpdateFlagEdgeApplicationId)
	flags.StringVar(&fields.Name, "name", "", msg.EdgeFuncInstanceUpdateFlagName)
	flags.Int64VarP(&fields.FunctionID, "function-id", "f", 0, msg.EdgeFuncInstanceUpdateFlagFunctionID)
	flags.StringVarP(&fields.InstanceID, "instance-id", "i", "0", msg.EdgeFuncInstanceUpdateFlagInstanceID)
	flags.StringVar(&fields.Path, "in", "", msg.EdgeFuncInstanceUpdateFlagIn)
	flags.BoolP("help", "h", false, msg.EdgeFuncInstanceUpdateHelpFlag)
	return cmd
}
