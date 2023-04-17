package device_groups

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/device_groups"
	"github.com/aziontech/azion-cli/pkg/cmd/device_groups/list"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	originsCmd := &cobra.Command{
		Use:   msg.DeviceGroupsUsage,
		Short: msg.DeviceGroupsShortDescription,
		Long:  msg.DeviceGroupsLongDescription, Example: heredoc.Doc(`
		$ azioncli device_groups --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	originsCmd.AddCommand(list.NewCmd(f))
	originsCmd.Flags().BoolP("help", "h", false, msg.DeviceGroupsFlagHelp)
	return originsCmd
}
