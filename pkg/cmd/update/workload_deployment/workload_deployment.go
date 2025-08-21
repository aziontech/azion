package workloaddeployment

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/workload_deployment"
	api "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Fields struct {
	Name            string `json:"name"`
	Active          string `json:"active"`
	Current         string `json:"current"`
	Path            string
	WorkloadID      int64
	DeploymentID    int64
	StrategyType    string `json:"strategy_type"`
	EdgeApplication string `json:"edge_application"`
	EdgeFirewall    string `json:"edge_firewall"`
	CustomPage      string `json:"custom_page"`
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update workload-deployment --workload-id 1234 --deployment-id 5678 --name 'Hello'
		$ azion update workload-deployment --workload-id 1234 --deployment-id 5678 --active true --current true
		$ azion update workload-deployment --workload-id 1234 --deployment-id 5678 --strategy-type blue-green --edge-application 123
		$ azion update workload-deployment --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := sdk.PatchedWorkloadDeploymentRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("workload-id") {
					answer, err := utils.AskInput(msg.AskInputWorkloadID)
					if err != nil {
						return err
					}
					num, err := strconv.ParseInt(answer, 10, 64)
					if err != nil {
						logger.Debug("Error while converting answer to int64", zap.Error(err))
						return msg.ErrorConvertWorkloadId
					}
					logger.Debug("Converted Workload ID", zap.Any("Workload ID", num))
					fields.WorkloadID = num
				}

				if !cmd.Flags().Changed("deployment-id") {
					answer, err := utils.AskInput(msg.AskInputDeploymentID)
					if err != nil {
						return err
					}
					num, err := strconv.ParseInt(answer, 10, 64)
					if err != nil {
						logger.Debug("Error while converting answer to int64", zap.Error(err))
						return msg.ErrorConvertDeploymentId
					}
					logger.Debug("Converted Deployment ID", zap.Any("Deployment ID", num))
					fields.DeploymentID = num
				}

				// Set name if provided
				if cmd.Flags().Changed("name") {
					name := fields.Name
					request.Name = &name
				}

				// Handle Active flag as pointer
				if cmd.Flags().Changed("active") {
					isActive, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorIsActiveFlag, fields.Active)
					}
					request.Active = &isActive
				}

				// Handle Current flag as pointer
				if cmd.Flags().Changed("current") {
					isCurrent, err := strconv.ParseBool(fields.Current)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorConvertCurrent, fields.Current)
					}
					request.Current = &isCurrent
				}

				// Handle Strategy
				if cmd.Flags().Changed("strategy-type") {
					strategy := sdk.DeploymentStrategyDefaultDeploymentStrategyRequest{
						Type: fields.StrategyType,
					}

					// Create attributes object
					attributes := sdk.DefaultDeploymentStrategyAttrsRequest{}

					// Set EdgeApplication if provided
					if cmd.Flags().Changed("edge-application") {
						edgeApp, err := strconv.ParseInt(fields.EdgeApplication, 10, 64)
						if err != nil {
							return fmt.Errorf("%w: %q", msg.ErrorConvertEdgeApplication, fields.EdgeApplication)
						}
						attributes.EdgeApplication = edgeApp
					}

					// Set EdgeFirewall if provided
					if cmd.Flags().Changed("edge-firewall") {
						edgeFirewall, err := strconv.ParseInt(fields.EdgeFirewall, 10, 64)
						if err != nil {
							return fmt.Errorf("%w: %q", msg.ErrorConvertEdgeFirewall, fields.EdgeFirewall)
						}
						var nullableEdgeFirewall sdk.NullableInt64
						nullableEdgeFirewall.Set(&edgeFirewall)
						attributes.EdgeFirewall = nullableEdgeFirewall
					}

					// Set CustomPage if provided
					if cmd.Flags().Changed("custom-page") {
						customPage, err := strconv.ParseInt(fields.CustomPage, 10, 64)
						if err != nil {
							return fmt.Errorf("%w: %q", msg.ErrorConvertCustomPage, fields.CustomPage)
						}
						var nullableCustomPage sdk.NullableInt64
						nullableCustomPage.Set(&customPage)
						attributes.CustomPage = nullableCustomPage
					}

					strategy.Attributes = attributes
					request.Strategy = &strategy
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			ctx := context.Background()
			response, err := client.UpdateDeployment(ctx, request, fields.WorkloadID, fields.DeploymentID)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateWorkloadDeployment.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	flags.Int64Var(&fields.WorkloadID, "workload-id", 0, msg.FlagWorkloadID)
	flags.Int64Var(&fields.DeploymentID, "deployment-id", 0, msg.FlagDeploymentID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Active, "active", "", msg.FlagIsActive)
	flags.StringVar(&fields.Current, "current", "", msg.FlagIsCurrent)
	flags.StringVar(&fields.StrategyType, "strategy-type", "", msg.FlagStrategyType)
	flags.StringVar(&fields.EdgeApplication, "edge-application", "", msg.FlagEdgeApplicationId)
	flags.StringVar(&fields.EdgeFirewall, "edge-firewall", "", msg.FlagEdgeFirewallId)
	flags.StringVar(&fields.CustomPage, "custom-page", "", msg.FlagCustomPageId)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}
