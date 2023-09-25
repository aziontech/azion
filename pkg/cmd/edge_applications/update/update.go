package update

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	Id                      int64
	Name                    string
	DeliveryProtocol        string
	HttpPort                int64
	HttpsPort               int64
	MinimumTlsVersion       string
	Active                  string
	ApplicationAcceleration string
	Caching                 string
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
		Use:           msg.UpdateUsage,
		Short:         msg.UpdateShortDescription,
		Long:          msg.UpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion edge_applications update --application-id 1234 --name 'Hello'
		$ azion edge_applications update -a 9123 --active true
		$ azion edge_applications update -a 9123 --active false
		$ azion edge_applications update --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			// either function-id or in path should be passed
			if !cmd.Flags().Changed("application-id") && !cmd.Flags().Changed("in") {
				return msg.ErrorMissingApplicationIdArgument
			}

			if !returnAnyField(cmd) {
				return msg.ErrorNoFieldInformed
			}

			request := api.UpdateRequest{}
			if cmd.Flags().Changed("in") {
				var (
					file *os.File
					err  error
				)
				if fields.InPath == "-" {
					file = os.Stdin
				} else {
					file, err = os.Open(fields.InPath)
					if err != nil {
						return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.InPath)
					}
				}
				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {

				request.Id = fields.Id

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				if cmd.Flags().Changed("http-port") {
					request.SetHttpPort(fields.HttpPort)
				}

				if cmd.Flags().Changed("https-port") {
					request.SetHttpsPort(fields.HttpsPort)
				}

				if cmd.Flags().Changed("delivery-protocol") {
					request.SetDeliveryProtocol(fields.DeliveryProtocol)
				}

				if cmd.Flags().Changed("min-tsl-ver") {
					request.SetMinimumTlsVersion(fields.MinimumTlsVersion)
				}

				if cmd.Flags().Changed("application-acceleration") {
					converted, err := strconv.ParseBool(fields.ApplicationAcceleration)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorApplicationAccelerationFlag, fields.Active)
					}
					request.SetApplicationAcceleration(converted)
				}

				if cmd.Flags().Changed("caching") {
					converted, err := strconv.ParseBool(fields.Caching)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorCachingFlag, fields.Active)
					}
					request.SetCaching(converted)
				}

				if cmd.Flags().Changed("device-detection") {
					converted, err := strconv.ParseBool(fields.DeviceDetection)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorDeviceDetectionFlag, fields.Active)
					}
					request.SetDeviceDetection(converted)
				}

				if cmd.Flags().Changed("edge-firewall") {
					converted, err := strconv.ParseBool(fields.EdgeFirewall)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorEdgeFirewallFlag, fields.Active)
					}
					request.SetEdgeFirewall(converted)
				}

				if cmd.Flags().Changed("edge-functions") {
					converted, err := strconv.ParseBool(fields.EdgeFunctions)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorEdgeFunctionsFlag, fields.Active)
					}
					request.SetEdgeFunctions(converted)
				}

				if cmd.Flags().Changed("image-optimization") {
					converted, err := strconv.ParseBool(fields.ImageOptimization)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorImageOptimizationFlag, fields.Active)
					}
					request.SetImageOptimization(converted)
				}

				if cmd.Flags().Changed("l2-caching") {
					converted, err := strconv.ParseBool(fields.L2Caching)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorL2CachingFlag, fields.Active)
					}
					request.SetL2Caching(converted)
				}

				if cmd.Flags().Changed("load-balancer") {
					converted, err := strconv.ParseBool(fields.LoadBalancer)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorLoadBalancerFlag, fields.Active)
					}
					request.SetLoadBalancer(converted)
				}

				if cmd.Flags().Changed("raw-logs") {
					converted, err := strconv.ParseBool(fields.RawLogs)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorRawLogsFlag, fields.Active)
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
				return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
			}

			fmt.Fprintf(f.IOStreams.Out, "Updated Edge Application with ID %d\n", response.GetId())

			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.Id, "application-id", "a", 0, msg.FlagId)
	flags.StringVar(&fields.Name, "name", "", msg.UpdateFlagName)
	flags.StringVar(&fields.DeliveryProtocol, "delivery-protocol", "", msg.UpdateFlagDeliveryProtocol)
	flags.Int64Var(&fields.HttpPort, "http-port", 80, msg.UpdateFlagHttpPort)
	flags.Int64Var(&fields.HttpsPort, "https-port", 443, msg.UpdateFlagHttpsPort)
	flags.StringVar(&fields.MinimumTlsVersion, "min-tsl-ver", "", msg.UpdateFlagMinimumTlsVersion)
	flags.StringVar(&fields.ApplicationAcceleration, "application-acceleration", "", msg.UpdateFlagApplicationAcceleration)
	flags.StringVar(&fields.Caching, "caching", "", msg.UpdateFlagCaching)
	flags.StringVar(&fields.DeviceDetection, "device-detection", "", msg.UpdateFlagDeviceDetection)
	flags.StringVar(&fields.EdgeFirewall, "edge-firewall", "", msg.UpdateFlagEdgeFirewall)
	flags.StringVar(&fields.EdgeFunctions, "edge-functions", "", msg.UpdateFlagEdgeFunctions)
	flags.StringVar(&fields.ImageOptimization, "image-optimization", "", msg.UpdateFlagImageOptimization)
	flags.StringVar(&fields.L2Caching, "l2-caching", "", msg.UpdateFlagL2Caching)
	flags.StringVar(&fields.LoadBalancer, "load-balancer", "", msg.UpdateFlagLoadBalancer)
	flags.StringVar(&fields.RawLogs, "raw-logs", "", msg.UpdateRawLogs)
	flags.StringVar(&fields.WebApplicationFirewall, "webapp-firewall", "", msg.UpdateWebApplicationFirewall)
	flags.StringVar(&fields.InPath, "in", "", msg.UpdateFlagIn)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
	return cmd
}

func returnAnyField(cmd *cobra.Command) bool {
	anyFlagChanged := false
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Changed && flag.Name != "application-id" && flag.Shorthand != "a" {
			anyFlagChanged = true
		}
	})
	return anyFlagChanged
}
