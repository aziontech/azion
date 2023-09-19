package update

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/device_groups"
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
		Use:           msg.DeviceGroupsUpdateUsage,
		Short:         msg.DeviceGroupsUpdateShortDescription,
		Long:          msg.DeviceGroupsUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion device_groups update --application-id 1673635839 --group-id 12312 --user-agent "(Mobile|iP(hone|od)|BlackBerry|IEMobile)"
        $ azion device_groups update -a 1673635839 -g 12312 --name "updated name"
        $ azion device_groups update -a 1673635839 -g 12312 --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Flags().Changed("in") && (!cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("group-id")) {
				return msg.ErrorMandatoryFlagsUpdate
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
					return msg.ErrorMandatoryFlagsUpdate
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
				return fmt.Errorf(msg.ErrorUpdateDeviceGroups.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.DeviceGroupsUpdateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.ApplicationFlagId)
	flags.Int64VarP(&fields.GroupID, "group-id", "g", 0, msg.DeviceGroupFlagId)
	flags.StringVar(&fields.Name, "name", "", msg.DeviceGroupsUpdateFlagName)
	flags.StringVar(&fields.UserAgent, "user-agent", "", msg.DeviceGroupsUpdateFlagUserAgent)
	flags.StringVar(&fields.Path, "in", "", msg.DeviceGroupsUpdateFlagIn)
	flags.BoolP("help", "h", false, msg.DeviceGroupsFlagHelp)
	return cmd
}
