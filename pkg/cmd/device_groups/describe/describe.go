package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/device_groups"
	"path/filepath"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID int64
	groupID       int64
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           device_groups.DeviceGroupsDescribeUsage,
		Short:         device_groups.DeviceGroupsDescribeShortDescription,
		Long:          device_groups.DeviceGroupsDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
      $ azion device_groups describe --application-id 1673635839 --group-id 31223
      $ azion device_groups describe -a 1673635839 -g 31223 --format json
      $ azion device_groups describe --application-id 1673635839 --group-id 31223 --out "./tmp/test.json"
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("group-id") {
				return device_groups.ErrorMandatoryFlags
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			groups, err := client.GetDeviceGroups(ctx, applicationID, groupID)
			if err != nil {
				return fmt.Errorf(device_groups.ErrorGetDeviceGroups.Error(), err)
			}

			out := f.IOStreams.Out
			formattedFuction, err := format(cmd, groups)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, device_groups.DeviceGroupsFileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedFuction[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, device_groups.ApplicationFlagId)
	cmd.Flags().Int64VarP(&groupID, "group-id", "g", 0, device_groups.DeviceGroupFlagId)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", device_groups.DeviceGroupsDescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", device_groups.DeviceGroupsDescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, device_groups.DeviceGroupsDescribeHelpFlag)

	return cmd
}

func format(cmd *cobra.Command, rules api.DeviceGroupsResponse) ([]byte, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(rules, "", " ")
	}

	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())
	tbl.AddRow("Device Group ID: ", rules.GetId())
	tbl.AddRow("Name: ", rules.GetName())
	tbl.AddRow("User Agent: ", rules.GetUserAgent())
	return tbl.GetByteFormat(), nil
}
