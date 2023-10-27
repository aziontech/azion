package edge_applications

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/MaxwelMazur/tablecli"
	msg "github.com/aziontech/azion-cli/messages/describe/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/fatih/color"
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
			application, err := client.Get(ctx, applicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetApplication.Error(), err)
			}

			out := f.IOStreams.Out
			formattedApp, err := format(cmd, application)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedApp, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, msg.FileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedApp[:])
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&applicationID, "application-id", "", msg.FlagId)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.FlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.FlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.HelpFlag)

	return cmd
}

func format(cmd *cobra.Command, application api.EdgeApplicationResponse) ([]byte, error) {
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		return json.MarshalIndent(application, "", " ")
	}

	tbl := tablecli.New("", "")
	tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())
	tbl.AddRow("ID: ", application.GetId())
	tbl.AddRow("Name: ", application.GetName())
	tbl.AddRow("Active: ", application.GetActive())
	tbl.AddRow("Application Acceleration: ", application.GetApplicationAcceleration())
	tbl.AddRow("Caching: ", application.GetCaching())
	tbl.AddRow("Delivery Protocol: ", application.GetDeliveryProtocol())
	tbl.AddRow("Device Detection: ", application.GetDeviceDetection())
	tbl.AddRow("Edge Firewall: ", application.GetEdgeFirewall())
	tbl.AddRow("Edge Functions: ", application.GetEdgeFunctions())
	tbl.AddRow("Http Port: ", application.GetHttpPort())
	tbl.AddRow("HttpsPort: ", application.GetHttpsPort())
	tbl.AddRow("Image Optimization: ", application.GetImageOptimization())
	tbl.AddRow("L2 Caching: ", application.GetL2Caching())
	tbl.AddRow("Load Balancer: ", application.GetLoadBalancer())
	tbl.AddRow("Minimum TLS Version: ", application.GetMinimumTlsVersion())
	tbl.AddRow("Raw Logs: ", application.GetRawLogs())
	tbl.AddRow("Web Application Firewall: ", application.GetWebApplicationFirewall())
	return tbl.GetByteFormat(), nil
}
