package devicegroups

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/device_groups"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type DeleteCmd struct {
	Io                *iostreams.IOStreams
	DeleteDeviceGroup func(context.Context, int64, int64) error
	AskInput          func(string) (string, error)
	ApplicationID     int64
	DeviceGroupID     int64
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		DeleteDeviceGroup: func(ctx context.Context, appID, groupID int64) error {
			client := registry.DeviceGroups(f)
			return client.Delete(ctx, appID, groupID)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(del *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.DeviceGroupsUsage,
		Short:         msg.DeviceGroupsDeleteShortDescription,
		Long:          msg.DeviceGroupsDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete device-group --application-id 1673635839 --group-id 107313
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := del.AskInput(msg.DeviceGroupsDeleteAskInputApplicationID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				del.ApplicationID = num
			}

			if !cmd.Flags().Changed("group-id") {
				answer, err := del.AskInput(msg.DeviceGroupsDeleteAskInputGroupID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdDeviceGroup
				}

				del.DeviceGroupID = num
			}

			ctx := context.Background()

			err := del.DeleteDeviceGroup(ctx, del.ApplicationID, del.DeviceGroupID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DeviceGroupsDeleteOutputSuccess, del.DeviceGroupID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&del.ApplicationID, "application-id", 0, msg.ApplicationFlagId)
	cobraCmd.Flags().Int64Var(&del.DeviceGroupID, "group-id", 0, msg.DeviceGroupFlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DeviceGroupsDeleteHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
