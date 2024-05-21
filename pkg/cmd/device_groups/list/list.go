package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/fatih/color"

	"github.com/MakeNowJust/heredoc"
	table "github.com/aziontech/tablecli"
	msg "github.com/aziontech/azion-cli/messages/device_groups"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var edgeApplicationID int64 = 0
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.DeviceGroupsListUsage,
		Short:         msg.DeviceGroupsListShortDescription,
		Long:          msg.DeviceGroupsListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true, Example: heredoc.Doc(`
        $ azion device_groups list -a 16736354321
        $ azion device_groups list --application-id 16736354321
        $ azion device_groups list --application-id 16736354321 --details
        $ azion device_groups list --application-id 16736354321 --order_by "id"
        $ azion device_groups list --application-id 16736354321 --page 1
        $ azion device_groups list --application-id 16736354321 --page_size 5
        $ azion device_groups list --application-id 16736354321 --sort "asc"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			var numberPage int64 = opts.Page
			if !cmd.Flags().Changed("application-id") {
				return msg.ErrorMissingApplicationIDArgument
			}
			if !cmd.Flags().Changed("page") && !cmd.Flags().Changed("page_size") {
				for {
					pages, err := PrintTable(cmd, f, opts, &edgeApplicationID, &numberPage)
					if numberPage > pages && err == nil {
						return nil
					}
					if err != nil {
						return fmt.Errorf(msg.ErrorListDeviceGroups.Error(), err)
					}
				}
			}

			if _, err := PrintTable(cmd, f, opts, &edgeApplicationID, &numberPage); err != nil {
				return fmt.Errorf(msg.ErrorGetDeviceGroups.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	flags := cmd.Flags()
	flags.Int64VarP(&edgeApplicationID, "application-id", "a", 0, msg.DeviceGroupsListFlagEdgeApplicationID)
	flags.BoolP("help", "h", false, msg.DeviceGroupsListHelpFlag)
	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, edgeApplicationID, numberPage *int64) (int64, error) {
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	applications, err := client.DeviceGroupsList(ctx, opts, *edgeApplicationID)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorGetDeviceGroups.Error(), err)
	}

	tbl := table.New("ID", "NAME")
	tbl.WithWriter(f.IOStreams.Out)
	table.DefaultWriter = f.IOStreams.Out
	if cmd.Flags().Changed("details") {
		tbl = table.New("ID", "NAME", "USER AGENT")
	}

	headerFmt := color.New(color.FgBlue, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgGreen).SprintfFunc()
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range applications.Results {
		if cmd.Flags().Changed("details") {
			tbl.AddRow(*v.Id, v.Name, v.UserAgent)
		} else {
			tbl.AddRow(*v.Id, v.Name)
		}
	}

	format := strings.Repeat("%s", len(tbl.GetHeader())) + "\n"
	tbl.CalculateWidths([]string{})
	if *numberPage == 1 {
		tbl.PrintHeader(format)
	}

	for _, row := range tbl.GetRows() {
		tbl.PrintRow(format, row)
	}

	*numberPage += 1
	opts.Page = *numberPage
	return applications.TotalPages, nil
}
