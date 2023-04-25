package device_groups

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/device_groups"
	"github.com/aziontech/azion-cli/pkg/cmd/device_groups/create"
	"github.com/aziontech/azion-cli/pkg/cmd/device_groups/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/device_groups/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/device_groups/list"
	"github.com/aziontech/azion-cli/pkg/cmd/device_groups/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	deviceGroupsCmd := &cobra.Command{
		Use:   msg.DeviceGroupsUsage,
		Short: msg.DeviceGroupsShortDescription,
		Long:  msg.DeviceGroupsLongDescription,
		Example: heredoc.Doc(`
		$ azioncli device_groups --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	deviceGroupsCmd.AddCommand(list.NewCmd(f))
	deviceGroupsCmd.AddCommand(delete.NewCmd(f))
	deviceGroupsCmd.AddCommand(create.NewCmd(f))
	deviceGroupsCmd.AddCommand(describe.NewCmd(f))
	deviceGroupsCmd.AddCommand(update.NewCmd(f))

	deviceGroupsCmd.Flags().BoolP("help", "h", false, msg.DeviceGroupsFlagHelp)
	return deviceGroupsCmd
}
