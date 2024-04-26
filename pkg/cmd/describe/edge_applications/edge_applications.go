package edge_applications

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var applicationID string
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion describe edge-application --application-id 4312
        $ azion describe edge-application --application-id 1337 --out "./tmp/test.json"
        $ azion describe edge-application --application-id 1337 --format json
        `),
		RunE: func(cmd *cobra.Command, _ []string) error {

			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(msg.AskInputApplicationID)
				if err != nil {
					return err
				}

				applicationID = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			resp, err := client.Get(ctx, applicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetApplication.Error(), err)
			}

			fields := make(map[string]string, 0)

			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["Active"] = "Active"
			fields["ApplicationAcceleration"] = "Application Acceleration"
			fields["Caching"] = "Caching"
			fields["DeliveryProtocol"] = "Delivery Protocol"
			fields["DeviceDetection"] = "Device Detection"
			fields["EdgeFirewall"] = "Edge Firewall"
			fields["EdgeFunctions"] = "Edge Functions"
			fields["HttpPort"] = "Http Port"
			fields["HttpsPort"] = "HttpsPort"
			fields["ImageOptimization"] = "Image Optimization"
			fields["L2Caching"] = "L2 Caching"
			fields["LoadBalancer"] = "Load Balancer"
			fields["MinimumTlsVersion"] = "Minimum TLS Version"
			fields["RawLogs"] = "Raw Logs"
			fields["WebApplicationFirewall"] = "Web Application Firewall"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:         f.IOStreams.Out,
					Msg:         filepath.Clean(opts.OutPath),
					FlagOutPath: opts.OutPath,
					FlagFormat:  opts.Format,
				},
				Fields: fields,
				Values: resp,
			}
			return output.Print(&describeOut)

		},
	}

	cmd.Flags().StringVar(&applicationID, "application-id", "", msg.FlagId)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.FlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.FlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}
