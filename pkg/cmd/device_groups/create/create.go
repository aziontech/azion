package create

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/device_groups"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID int64
	Name          string
	UserAgent     string
	Path          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DeviceGroupsCreateUsage,
		Short:         msg.DeviceGroupsCreateShortDescription,
		Long:          msg.DeviceGroupsCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion device_groups create --application-id 1673635839 --name "asdf" --user-agent "httpbin.org"
        $ azion device_groups create -a 1673635839 --name "asdf" --user-agent "httpbin.org"
        $ azion device_groups create -a 1673635839 --in "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.CreateDeviceGroupsRequest{}
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
				if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("name") ||
					!cmd.Flags().Changed("user-agent") { // flags requireds
					return msg.ErrorMandatoryCreateFlags
				}

				request.SetName(fields.Name)
				request.SetUserAgent(fields.UserAgent)
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.CreateDeviceGroups(context.Background(), &request, fields.ApplicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateDeviceGroups.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.DeviceGroupsCreateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.DeviceGroupsCreateFlagEdgeApplicationId)
	flags.StringVar(&fields.Name, "name", "", msg.DeviceGroupsCreateFlagName)
	flags.StringVar(&fields.UserAgent, "user-agent", "", msg.DeviceGroupsCreateFlagUserAgent)
	flags.StringVar(&fields.Path, "in", "", msg.DeviceGroupsCreateFlagIn)
	flags.BoolP("help", "h", false, msg.DeviceGroupsFlagHelp)
	return cmd
}
