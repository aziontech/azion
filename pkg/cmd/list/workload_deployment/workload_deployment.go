package workloaddeployment

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list/workload_deployment"
	api "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var workloadId int64

type ListCmd struct {
	Io                      *iostreams.IOStreams
	ReadInput               func(string) (string, error)
	ListWorkloadDeployments func(context.Context, *contracts.ListOptions, int64) (*sdk.PaginatedWorkloadDeploymentList, error)
	AskInput                func(string) (string, error)
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		ListWorkloadDeployments: func(ctx context.Context, opts *contracts.ListOptions, id int64) (*sdk.PaginatedWorkloadDeploymentList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.ListDeployments(ctx, opts, id)
		},
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
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
			$ azion list workload-deployment
			$ azion list workload-deployment --details
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("workload-id") {
				answer, err := list.ReadInput(msg.AskListInput)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertId
				}

				workloadId = num
			}

			if err := PrintTable(cmd, f, list, opts); err != nil {
				return msg.ErrorGetWorkloadDeployments
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)
	cmd.Flags().Int64Var(&workloadId, "workload-id", 0, msg.WorkloadIdFlag)

	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, list *ListCmd, opts *contracts.ListOptions) error {
	ctx := context.Background()

	response, err := list.ListWorkloadDeployments(ctx, opts, workloadId)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "CURRENT"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	if opts.Details {
		listOut.Columns = []string{"ID", "CURRENT", "EDGE APPLICATION", "EDGE FIREWALL"}
	}

	for _, v := range response.Results {
		var ln []string
		stratety := v.GetStrategy()
		if opts.Details {
			ln = []string{
				fmt.Sprintf("%d", v.GetId()),
				fmt.Sprintf("%v", v.GetCurrent()),
				fmt.Sprintf("%d", stratety.Attributes.GetApplication()),
				fmt.Sprintf("%d", stratety.Attributes.GetFirewall()),
			}
		} else {
			ln = []string{
				fmt.Sprintf("%d", v.GetId()),
				fmt.Sprintf("%v", v.GetCurrent()),
			}
		}

		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}
