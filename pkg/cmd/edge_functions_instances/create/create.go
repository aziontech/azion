package create

import (
	"context"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_functions_instances"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"os"
)

type Fields struct {
	ApplicationID  int64
	Name           string
	EdgeFunctionId int64
	Args           string
	Path           string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.EdgeFuncInstanceCreateUsage,
		Short:         msg.EdgeFuncInstanceCreateShortDescription,
		Long:          msg.EdgeFuncInstanceCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
 		$ azioncli edge_functions_instances create --application-id 1673635839 --instance-id 12314 --name "ffcafe222sdsdffdf"
		$ azioncli edge_functions_instances create -a 1673635839 -i 12314 --name "ffcafe222sdsdffdf"
        $ azioncli edge_functions_instances create -a 1673635839 --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.CreateFuncInstancesRequest{}
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
				if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("instance-id") || !cmd.Flags().Changed("name") {
					return msg.ErrorMandatoryCreateFlags
				}
				request.SetName(fields.Name)
				request.SetEdgeFunctionId(fields.EdgeFunctionId)
				request.SetArgs(fields.Args)
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.CreateFuncInstances(context.Background(), &request, fields.ApplicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.EdgeFuncInstanceCreateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.EdgeFuncInstanceCreateFlagEdgeApplicationId)
	flags.Int64VarP(&fields.EdgeFunctionId, "instance-id", "i", 0, msg.EdgeFuncInstanceCreateFlagEdgeFunctionID)
	flags.StringVar(&fields.Name, "name", "", msg.EdgeFuncInstanceCreateFlagName)
	flags.StringVar(&fields.Args, "args", "", msg.EdgeFuncInstanceCreateFlagArgs)
	flags.StringVar(&fields.Path, "in", "", msg.EdgeFuncInstanceCreateFlagIn)
	flags.BoolP("help", "h", false, msg.EdgeFuncInstanceCreateHelpFlag)
	return cmd
}
