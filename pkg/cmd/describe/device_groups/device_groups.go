package devicegroups

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/device_groups"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
)

type DescribeCmd struct {
	Io            *iostreams.IOStreams
	AskInput      func(string) (string, error)
	Get           func(context.Context, int64, int64) (sdk.DeviceGroup, error)
	ApplicationID int64
	DeviceGroupID int64
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, appID, groupID int64) (sdk.DeviceGroup, error) {
			client := registry.DeviceGroups(f)
			return client.Get(ctx, appID, groupID)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
		Use:           msg.DeviceGroupsUsage,
		Short:         msg.DeviceGroupsDescribeShortDescription,
		Long:          msg.DeviceGroupsDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion describe device-group --application-id 1673635839 --group-id 107313
        $ azion describe device-group --application-id 1673635839 --group-id 107313 --format json
        $ azion describe device-group --application-id 1673635839 --group-id 107313 --out "./tmp/test.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := describe.AskInput(msg.DeviceGroupsDescribeAskInputApplicationID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				describe.ApplicationID = num
			}

			if !cmd.Flags().Changed("group-id") {
				answer, err := describe.AskInput(msg.DeviceGroupsDescribeAskInputGroupID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdDeviceGroup
				}

				describe.DeviceGroupID = num
			}

			ctx := context.Background()
			resp, err := describe.Get(ctx, describe.ApplicationID, describe.DeviceGroupID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetDeviceGroups.Error(), err)
			}

			fields := make(map[string]string, 0)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["UserAgent"] = "User Agent"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: &resp,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&describe.ApplicationID, "application-id", 0, msg.ApplicationFlagId)
	cobraCmd.Flags().Int64Var(&describe.DeviceGroupID, "group-id", 0, msg.DeviceGroupFlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DeviceGroupsDescribeHelpFlag)
	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
