package edge_application

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update/edge_application"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

// Fields struct of inputs
type Fields struct {
	ID                      int64
	Name                    string
	DeliveryProtocol        string
	HTTPPort                int64
	HTTPSPort               int64
	MinimumTLSVersion       string
	ApplicationAcceleration string
	DeviceDetection         string
	EdgeFirewall            string
	EdgeFunctions           string
	ImageOptimization       string
	L2Caching               string
	LoadBalancer            string
	RawLogs                 string
	WebApplicationFirewall  string
	InPath                  string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           "edge-application",
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update edge-application --application-id 1234 --name 'Hello'
		$ azion update edge-application --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") && !cmd.Flags().Changed("file") {

				answer, err := utils.AskInput(msg.AskInputApplicationId)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				fields.ID = num
			}

			if !returnAnyField(cmd) {
				return msg.ErrorNoFieldInformed
			}

			request := api.UpdateRequest{}
			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.InPath+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {

				request.Id = fields.ID

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				if cmd.Flags().Changed("http-port") {
					request.SetHttpPort(fields.HTTPPort)
				}

				if cmd.Flags().Changed("https-port") {
					request.SetHttpsPort(fields.HTTPSPort)
				}

				if cmd.Flags().Changed("delivery-protocol") {
					request.SetDeliveryProtocol(fields.DeliveryProtocol)
				}

				if cmd.Flags().Changed("min-tsl-ver") {
					request.SetMinimumTlsVersion(fields.MinimumTLSVersion)
				}

				if cmd.Flags().Changed("application-acceleration") {
					converted, err := strconv.ParseBool(fields.ApplicationAcceleration)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorApplicationAccelerationFlag, fields.ApplicationAcceleration)
					}
					request.SetApplicationAcceleration(converted)
				}

				if cmd.Flags().Changed("device-detection") {
					converted, err := strconv.ParseBool(fields.DeviceDetection)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorDeviceDetectionFlag, fields.DeviceDetection)
					}
					request.SetDeviceDetection(converted)
				}

				if cmd.Flags().Changed("edge-firewall") {
					converted, err := strconv.ParseBool(fields.EdgeFirewall)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorEdgeFirewallFlag, fields.EdgeFirewall)
					}
					request.SetEdgeFirewall(converted)
				}

				if cmd.Flags().Changed("edge-functions") {
					converted, err := strconv.ParseBool(fields.EdgeFunctions)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorEdgeFunctionsFlag, fields.EdgeFunctions)
					}
					request.SetEdgeFunctions(converted)
				}

				if cmd.Flags().Changed("image-optimization") {
					converted, err := strconv.ParseBool(fields.ImageOptimization)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorImageOptimizationFlag, fields.ImageOptimization)
					}
					request.SetImageOptimization(converted)
				}

				if cmd.Flags().Changed("l2-caching") {
					converted, err := strconv.ParseBool(fields.L2Caching)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorL2CachingFlag, fields.L2Caching)
					}
					request.SetL2Caching(converted)
				}

				if cmd.Flags().Changed("load-balancer") {
					converted, err := strconv.ParseBool(fields.LoadBalancer)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorLoadBalancerFlag, fields.LoadBalancer)
					}
					request.SetLoadBalancer(converted)
				}

				if cmd.Flags().Changed("raw-logs") {
					converted, err := strconv.ParseBool(fields.RawLogs)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorRawLogsFlag, fields.RawLogs)
					}
					request.SetRawLogs(converted)
				}

				if cmd.Flags().Changed("webapp-firewall") {
					converted, err := strconv.ParseBool(fields.WebApplicationFirewall)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorWebApplicationFirewallFlag, fields.WebApplicationFirewall)
					}
					request.SetWebApplicationFirewall(converted)
				}

			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err.Error())
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
	flags.Int64Var(&fields.ID, "application-id", 0, msg.FlagID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.DeliveryProtocol, "delivery-protocol", "", msg.FlagDeliveryProtocol)
	flags.Int64Var(&fields.HTTPPort, "http-port", 80, msg.FlagHttpPort)
	flags.Int64Var(&fields.HTTPSPort, "https-port", 443, msg.FlagHttpsPort)
	flags.StringVar(&fields.MinimumTLSVersion, "min-tsl-ver", "", msg.FlagMinimumTlsVersion)
	flags.StringVar(&fields.ApplicationAcceleration, "application-acceleration", "", msg.FlagApplicationAcceleration)
	flags.StringVar(&fields.DeviceDetection, "device-detection", "", msg.FlagDeviceDetection)
	flags.StringVar(&fields.EdgeFirewall, "edge-firewall", "", msg.FlagFirewall)
	flags.StringVar(&fields.EdgeFunctions, "edge-functions", "", msg.FlagFunctions)
	flags.StringVar(&fields.ImageOptimization, "image-optimization", "", msg.FlagImageOptimization)
	flags.StringVar(&fields.L2Caching, "l2-caching", "", msg.FlagL2Caching)
	flags.StringVar(&fields.LoadBalancer, "load-balancer", "", msg.FlagLoadBalancer)
	flags.StringVar(&fields.RawLogs, "raw-logs", "", msg.RawLogs)
	flags.StringVar(&fields.WebApplicationFirewall, "webapp-firewall", "", msg.WebApplicationFirewall)
	flags.StringVar(&fields.InPath, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.HelpFlag)
	return cmd
}

func returnAnyField(cmd *cobra.Command) bool {
	anyFlagChanged := false
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Changed && flag.Name != "application-id" {
			anyFlagChanged = true
		}
	})
	return anyFlagChanged
}
