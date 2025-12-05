package functioninstance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/function_instance"
	api "github.com/aziontech/azion-cli/pkg/api/function_instance"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	ApplicationID int64
)

type ListCmd struct {
	Io            *iostreams.IOStreams
	AskInput      func(string) (string, error)
	ListInstances func(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.PaginatedApplicationFunctionInstanceList, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		ListInstances: func(ctx context.Context, opts *contracts.ListOptions, edgeApplicationID int64) (*sdk.PaginatedApplicationFunctionInstanceList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.List(ctx, opts, edgeApplicationID)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list function-instance
			$ azion list function-instance --application-id 16736354321 --details
			$ azion list function-instance --application-id 16736354321 --page 1
			$ azion list function-instance --application-id 16736354321 --page-size 5
		`),
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !cmd.Flags().Changed("application-id") {
				answer, err := list.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertApplicationId
				}

				ApplicationID = num
			}

			if err := PrintTable(cmd, list, f, opts); err != nil {
				return fmt.Errorf(msg.ErrorGetAll, err)
			}
			return nil
		},
	}

	flags := cmd.Flags()
	flags.BoolP("help", "h", false, msg.HelpFlag)
	flags.Int64Var(&ApplicationID, "application-id", 0, msg.ApplicationIdFlag)
	cmdutil.AddAzionApiFlags(cmd, opts)

	return cmd
}

func PrintTable(cmd *cobra.Command, list *ListCmd, f *cmdutil.Factory, opts *contracts.ListOptions) error {
	ctx := context.Background()

	resp, err := list.ListInstances(ctx, opts, ApplicationID)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "ACTIVE", "FUNCTION"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "NAME", "ACTIVE", "FUNCTION", "LAST EDITOR", "LAST MODIFIED"}
	}

	if resp == nil || len(resp.Results) == 0 {
		return output.Print(&listOut)
	}

	for _, v := range resp.Results {
		var ln []string
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.Id),
				utils.TruncateString(v.Name),
				fmt.Sprintf("%v", *v.Active),
				fmt.Sprintf("%d", v.Function),
				v.LastEditor,
				v.LastModified.String(),
			}
		} else {
			ln = []string{
				fmt.Sprintf("%d", v.Id),
				utils.TruncateString(v.Name),
				fmt.Sprintf("%v", *v.Active),
				fmt.Sprintf("%d", v.Function),
			}
		}

		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	listCmd := NewListCmd(f)
	return NewCobraCmd(listCmd, f)
}
