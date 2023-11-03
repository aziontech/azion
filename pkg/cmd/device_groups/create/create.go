package create

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/device_groups"
	"os"

	"github.com/MakeNowJust/heredoc"

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
		Use:           device_groups.DeviceGroupsCreateUsage,
		Short:         device_groups.DeviceGroupsCreateShortDescription,
		Long:          device_groups.DeviceGroupsCreateLongDescription,
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
					return device_groups.ErrorMandatoryCreateFlags
				}

				request.SetName(fields.Name)
				request.SetUserAgent(fields.UserAgent)
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.CreateDeviceGroups(context.Background(), &request, fields.ApplicationID)
			if err != nil {
				return fmt.Errorf(device_groups.ErrorCreateDeviceGroups.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, device_groups.DeviceGroupsCreateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, device_groups.DeviceGroupsCreateFlagEdgeApplicationId)
	flags.StringVar(&fields.Name, "name", "", device_groups.DeviceGroupsCreateFlagName)
	flags.StringVar(&fields.UserAgent, "user-agent", "", device_groups.DeviceGroupsCreateFlagUserAgent)
	flags.StringVar(&fields.Path, "in", "", device_groups.DeviceGroupsCreateFlagIn)
	flags.BoolP("help", "h", false, device_groups.DeviceGroupsFlagHelp)
	return cmd
}
