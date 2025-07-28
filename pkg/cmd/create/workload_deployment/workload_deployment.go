package workloaddeployment

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create/workload_deployment"
	api "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var workloadId int64

type Fields struct {
	WorkloadId       string `json:"workload_id"`
	EdgeApplication  string `json:"edge_application"`
	EdgeFirewall     string `json:"edge_firewall"`
	Path             string
}

type CreateCmd struct {
	Io               *iostreams.IOStreams
	CreateDeployment func(context.Context, int64, *api.CreateDeploymentRequest) (api.DeploymentResponse, error)
	AskInput         func(string) (string, error)
}

func NewCreateCmd(f *cmdutil.Factory) *CreateCmd {
	return &CreateCmd{
		Io: f.IOStreams,
		CreateDeployment: func(ctx context.Context, workloadId int64, req *api.CreateDeploymentRequest) (api.DeploymentResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.CreateDeployment(ctx, workloadId, req)
		},
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
	}
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}
	createCmd := NewCreateCmd(f)

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion create workload-deployment --workload-id 1234 --edge-application 5678
			$ azion create workload-deployment --workload-id 1234 --edge-application 5678 --edge-firewall 9012
			$ azion create workload-deployment --file "create.json"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := &api.CreateDeploymentRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				// Get workload ID
				if !cmd.Flags().Changed("workload-id") {
					answer, err := createCmd.AskInput(msg.AskInputWorkloadId)
					if err != nil {
						return err
					}
					fields.WorkloadId = answer
				}

				workloadIdInt, err := strconv.ParseInt(fields.WorkloadId, 10, 64)
				if err != nil {
					logger.Debug("Error while converting workload ID to int64", zap.Error(err))
					return msg.ErrorConvertWorkloadId
				}
				workloadId = workloadIdInt

				// Get edge application ID
				if !cmd.Flags().Changed("edge-application") {
					answer, err := createCmd.AskInput(msg.AskInputEdgeApplication)
					if err != nil {
						return err
					}
					fields.EdgeApplication = answer
				}

				edgeAppId, err := strconv.ParseInt(fields.EdgeApplication, 10, 64)
				if err != nil {
					logger.Debug("Error while converting edge application ID to int64", zap.Error(err))
					return msg.ErrorConvertEdgeApplication
				}

				// Create deployment strategy attributes
				attributes := sdk.NewDefaultDeploymentStrategyAttrsRequest(edgeAppId)

				// Set edge firewall if provided
				if cmd.Flags().Changed("edge-firewall") && fields.EdgeFirewall != "" {
					edgeFirewallId, err := strconv.ParseInt(fields.EdgeFirewall, 10, 64)
					if err != nil {
						logger.Debug("Error while converting edge firewall ID to int64", zap.Error(err))
						return msg.ErrorConvertEdgeFirewall
					}
					attributes.SetEdgeFirewall(edgeFirewallId)
				}

				// Create deployment strategy
				strategy := sdk.NewDeploymentStrategyDefaultDeploymentStrategyRequest("default", *attributes)

				// Create deployment request
				deploymentName := fmt.Sprintf("deployment-%d", workloadId)
				request.WorkloadDeploymentRequest = *sdk.NewWorkloadDeploymentRequest(deploymentName, *strategy)
			}

			response, err := createCmd.CreateDeployment(context.Background(), workloadId, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}

			createOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&createOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.WorkloadId, "workload-id", "", msg.FlagWorkloadId)
	flags.StringVar(&fields.EdgeApplication, "edge-application", "", msg.FlagEdgeApplication)
	flags.StringVar(&fields.EdgeFirewall, "edge-firewall", "", msg.FlagEdgeFirewall)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}
