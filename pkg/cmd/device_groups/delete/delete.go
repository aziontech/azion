package delete

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/device_groups"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var applicationID, groupID int64

	cmd := &cobra.Command{
		Use:           device_groups.DeviceGroupsDeleteUsage,
		Short:         device_groups.DeviceGroupsDeleteShortDescription,
		Long:          device_groups.DeviceGroupsDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azion device_groups delete --application-id 1234 --group-id 12312
		  $ azion device_groups delete -a 1234 -g 12312
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("group-id") {
				return device_groups.ErrorMandatoryFlags
			}
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.DeleteDeviceGroup(ctx, applicationID, groupID)
			if err != nil {
				return fmt.Errorf(device_groups.ErrorFailToDelete.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, device_groups.DeviceGroupsDeleteOutputSuccess, groupID)
			return nil
		},
	}

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, device_groups.ApplicationFlagId)
	cmd.Flags().Int64VarP(&groupID, "group-id", "g", 0, device_groups.DeviceGroupFlagId)
	cmd.Flags().BoolP("help", "h", false, device_groups.DeviceGroupsDeleteHelpFlag)
	return cmd
}
