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
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io                *iostreams.IOStreams
	ListDeviceGroups  func(context.Context, *contracts.ListOptions, int64) (*sdk.PaginatedDeviceGroupList, error)
	AskInput          func(string) (string, error)
	EdgeApplicationID int64
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListDeviceGroups: func(ctx context.Context, opts *contracts.ListOptions, appID int64) (*sdk.PaginatedDeviceGroupList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.List(ctx, opts, appID)
		},
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.DeviceGroupsUsage,
		Short:         msg.DeviceGroupsListShortDescription,
		Long:          msg.DeviceGroupsListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list device-group --application-id 1673635839
			$ azion list device-group --application-id 1673635839 --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := list.AskInput(msg.DeviceGroupsListAskInputApplicationID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				list.EdgeApplicationID = num
			}

			if err := PrintTable(cmd, f, opts, list); err != nil {
				return fmt.Errorf(msg.ErrorListDeviceGroups.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().Int64Var(&list.EdgeApplicationID, "application-id", 0, msg.DeviceGroupsListFlagEdgeApplicationID)
	cmd.Flags().BoolP("help", "h", false, msg.DeviceGroupsListHelpFlag)

	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, list *ListCmd) error {
	ctx := context.Background()

	response, err := list.ListDeviceGroups(ctx, opts, list.EdgeApplicationID)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "USER AGENT"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	for _, v := range response.GetResults() {
		ln := []string{
			fmt.Sprintf("%d", v.GetId()),
			v.GetName(),
			v.GetUserAgent(),
		}
		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}
