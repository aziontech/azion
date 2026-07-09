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
	Name          string
	UserAgent     string
	Path          string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DeviceGroupsUsage,
		Short:         msg.DeviceGroupsCreateShortDescription,
		Long:          msg.DeviceGroupsCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create device-group --application-id 1673635839 --name "mobile" --user-agent "(Mobile|iP(hone|od)|Android)"
        $ azion create device-group --application-id 1673635839 --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := sdk.DeviceGroupRequest{}

			if !cmd.Flags().Changed("application-id") {
				answers, err := utils.AskInput(msg.DeviceGroupsCreateAskInputApplicationID)
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

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("name") {
					answers, err := utils.AskInput(msg.DeviceGroupsCreateAskInputName)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					fields.Name = answers
				}

				if !cmd.Flags().Changed("user-agent") {
					answers, err := utils.AskInput(msg.DeviceGroupsCreateAskInputUserAgent)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					fields.UserAgent = answers
				}

				request.SetName(fields.Name)
				request.SetUserAgent(fields.UserAgent)
			}

			response, err := client.Create(context.Background(), request, fields.ApplicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateDeviceGroups.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeviceGroupsCreateOutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&creatOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)
	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.DeviceGroupsCreateFlagEdgeApplicationId)
	flags.StringVar(&fields.Name, "name", "", msg.DeviceGroupsCreateFlagName)
	flags.StringVar(&fields.UserAgent, "user-agent", "", msg.DeviceGroupsCreateFlagUserAgent)
	flags.StringVar(&fields.Path, "file", "", msg.DeviceGroupsCreateFlagIn)
	flags.BoolP("help", "h", false, msg.DeviceGroupsCreateHelpFlag)
}
