package devicegroups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/device_groups"
	api "github.com/aziontech/azion-cli/pkg/api/device_groups"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		Use:           msg.DeviceGroupsUsage,
		Short:         msg.DeviceGroupsUpdateShortDescription,
		Long:          msg.DeviceGroupsUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion update device-group --application-id 1673635839 --group-id 12312 --name "mobile"
        $ azion update device-group --application-id 1673635839 --group-id 12312 --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := sdk.PatchedDeviceGroupRequest{}

			if !cmd.Flags().Changed("application-id") {
				answers, err := utils.AskInput(msg.DeviceGroupsUpdateAskInputApplicationID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				applicationID, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return utils.ErrorConvertingStringToInt
				}

				fields.ApplicationID = int64(applicationID)
			}

			if !cmd.Flags().Changed("group-id") {
				answers, err := utils.AskInput(msg.DeviceGroupsUpdateAskInputGroupID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				groupID, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return utils.ErrorConvertingStringToInt
				}

				fields.GroupID = int64(groupID)
			}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				// The device group update is a full replacement (PUT), and both
				// name and user_agent are required. Fetch the current values for
				// any attribute the user didn't provide so they are preserved.
				if !cmd.Flags().Changed("name") || !cmd.Flags().Changed("user-agent") {
					current, err := client.Get(context.Background(), fields.ApplicationID, fields.GroupID)
					if err != nil {
						return fmt.Errorf(msg.ErrorGetDeviceGroups.Error(), err)
					}

					if !cmd.Flags().Changed("name") {
						fields.Name = current.GetName()
					}

					if !cmd.Flags().Changed("user-agent") {
						fields.UserAgent = current.GetUserAgent()
					}
				}

				request.SetName(fields.Name)
				request.SetUserAgent(fields.UserAgent)
			}

			response, err := client.Update(context.Background(), request, fields.ApplicationID, fields.GroupID)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateDeviceGroups.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeviceGroupsUpdateOutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)
	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.ApplicationFlagId)
	flags.Int64Var(&fields.GroupID, "group-id", 0, msg.DeviceGroupFlagId)
	flags.StringVar(&fields.Name, "name", "", msg.DeviceGroupsUpdateFlagName)
	flags.StringVar(&fields.UserAgent, "user-agent", "", msg.DeviceGroupsUpdateFlagUserAgent)
	flags.StringVar(&fields.Path, "file", "", msg.DeviceGroupsUpdateFlagIn)
	flags.BoolP("help", "h", false, msg.DeviceGroupsUpdateHelpFlag)
}
