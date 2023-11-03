package update

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/device_groups"
	"os"

	"github.com/MakeNowJust/heredoc"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID int64
	GroupID       int64
	Name          string
	UserAgent     string
	Path          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           device_groups.DeviceGroupsUpdateUsage,
		Short:         device_groups.DeviceGroupsUpdateShortDescription,
		Long:          device_groups.DeviceGroupsUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion device_groups update --application-id 1673635839 --group-id 12312 --user-agent "(Mobile|iP(hone|od)|BlackBerry|IEMobile)"
        $ azion device_groups update -a 1673635839 -g 12312 --name "updated name"
        $ azion device_groups update -a 1673635839 -g 12312 --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("in") && (!cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("group-id")) {
				return device_groups.ErrorMandatoryFlagsUpdate
			}

			request := sdk.PatchDeviceGroupsRequest{}
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
				if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("group-id") {
					return device_groups.ErrorMandatoryFlagsUpdate
				}
				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				if cmd.Flags().Changed("user-agent") {
					request.SetUserAgent(fields.UserAgent)
				}

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := client.UpdateDeviceGroup(context.Background(), request, fields.ApplicationID, fields.GroupID)
			if err != nil {
				return fmt.Errorf(device_groups.ErrorUpdateDeviceGroups.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, device_groups.DeviceGroupsUpdateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, device_groups.ApplicationFlagId)
	flags.Int64VarP(&fields.GroupID, "group-id", "g", 0, device_groups.DeviceGroupFlagId)
	flags.StringVar(&fields.Name, "name", "", device_groups.DeviceGroupsUpdateFlagName)
	flags.StringVar(&fields.UserAgent, "user-agent", "", device_groups.DeviceGroupsUpdateFlagUserAgent)
	flags.StringVar(&fields.Path, "in", "", device_groups.DeviceGroupsUpdateFlagIn)
	flags.BoolP("help", "h", false, device_groups.DeviceGroupsFlagHelp)
	return cmd
}
