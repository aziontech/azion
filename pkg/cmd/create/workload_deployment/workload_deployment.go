package workloaddeployment

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/workload_deployment"
	api "github.com/aziontech/azion-cli/pkg/api/workloads"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge-api"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name            string `json:"name"`
	Active          string `json:"active"`
	Current         string `json:"current"`
	Path            string
	WorkloadID      int64
	ApplicationID   int64
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
        $ azion create workload-deployment --name workloadName
        $ azion create workload-deployment --name withargs --active true --current true
        $ azion create workload-deployment --name withstrategy --strategy-type blue-green --edge-application 123
        $ azion create workload-deployment --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := sdk.WorkloadDeploymentRequest{}
			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {

				if !cmd.Flags().Changed("name") {
					answer, err := utils.AskInput(msg.AskInputName)
					if err != nil {
						return err
					}

					fields.Name = answer
				}

				// Check if workload-id is provided, if not ask for it
				if fields.WorkloadID == 0 {
					workloadIDStr, err := utils.AskInput(msg.AskInputWorkloadID)
					if err != nil {
						return err
					}

					workloadID, err := strconv.ParseInt(workloadIDStr, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid workload ID: %q", workloadIDStr)
					}
					fields.WorkloadID = workloadID
				}

				// Check if application-id is provided, if not ask for it
				if fields.ApplicationID == 0 {
					appIDStr, err := utils.AskInput(msg.AskInputApplicationID)
					if err != nil {
						return err
					}

					appID, err := strconv.ParseInt(appIDStr, 10, 64)
					if err != nil {
						return fmt.Errorf("invalid application ID: %q", appIDStr)
					}
					fields.ApplicationID = appID
				}

				request.SetName(fields.Name)

				// Handle Active flag as pointer
				if cmd.Flags().Changed("active") {
					isActive, err := strconv.ParseBool(fields.Active)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorIsActiveFlag, fields.Active)
					}
					active := isActive
					request.Active = &active
				}

				// Handle Current flag as pointer
				if cmd.Flags().Changed("current") {
					isCurrent, err := strconv.ParseBool(fields.Current)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorIsActiveFlag, fields.Current)
					}
					current := isCurrent
					request.Current = &current
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
					request.Strategy = strategy
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			response, err := client.CreateDeployment(context.Background(), request, fields.WorkloadID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateWorkloadDeployment.Error(), err)
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
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.Active, "active", "", msg.FlagIsActive)
	flags.StringVar(&fields.Current, "current", "", msg.FlagIsCurrent)
	flags.StringVar(&fields.StrategyType, "strategy-type", "", msg.FlagStrategyType)
	flags.StringVar(&fields.EdgeApplication, "edge-application", "", "Edge Application ID for the deployment strategy")
	flags.StringVar(&fields.EdgeFirewall, "edge-firewall", "", "Edge Firewall ID for the deployment strategy")
	flags.StringVar(&fields.CustomPage, "custom-page", "", "Custom Page ID for the deployment strategy")
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.Int64Var(&fields.WorkloadID, "workload-id", 0, "Workload ID")
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, "Application ID")
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}
