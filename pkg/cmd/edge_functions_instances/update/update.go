package update

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/edge_functions_instances"
	"os"

	"github.com/MakeNowJust/heredoc"

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
	Args          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           edge_functions_instances.EdgeFuncInstanceUpdateUsage,
		Short:         edge_functions_instances.EdgeFuncInstanceUpdateShortDescription,
		Long:          edge_functions_instances.EdgeFuncInstanceUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_functions_instances update -a 1674767911 -i 43121 -f 1209
        $ azion edge_functions_instances update --application-id 1674767911 --instance-id 2121 --function-id 1212 --name updated
        $ azion edge_functions_instances update  -a 1674767911 -i 43121 --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("in") && (!cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("instance-id")) {
				return edge_functions_instances.ErrorMandatoryUpdateFlagsIn
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
					return edge_functions_instances.ErrorMandatoryUpdateFlags
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
				return fmt.Errorf(edge_functions_instances.ErrorUpdateFuncInstance.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, edge_functions_instances.EdgeFuncInstanceUpdateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&fields.ApplicationID, "application-id", "a", "0", edge_functions_instances.EdgeFuncInstanceUpdateFlagEdgeApplicationId)
	flags.StringVar(&fields.Name, "name", "", edge_functions_instances.EdgeFuncInstanceUpdateFlagName)
	flags.StringVar(&fields.Args, "args", "", edge_functions_instances.EdgeFuncInstanceUpdateFlagArgs)
	flags.Int64VarP(&fields.FunctionID, "function-id", "f", 0, edge_functions_instances.EdgeFuncInstanceUpdateFlagFunctionID)
	flags.StringVarP(&fields.InstanceID, "instance-id", "i", "0", edge_functions_instances.EdgeFuncInstanceUpdateFlagInstanceID)
	flags.StringVar(&fields.Path, "in", "", edge_functions_instances.EdgeFuncInstanceUpdateFlagIn)
	flags.BoolP("help", "h", false, edge_functions_instances.EdgeFuncInstanceUpdateHelpFlag)
	return cmd
}
