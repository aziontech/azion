package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var application_id string
	opts := &contracts.DescribeOptions{}
	cmd := &cobra.Command{
		Use:           msg.DescribeUsage,
		Short:         msg.DescribeShortDescription,
		Long:          msg.DescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_applications describe --application-id 4312
        $ azion edge_applications describe --application-id 1337 --out "./tmp/test.json" --format json
        $ azion edge_applications describe --application-id 1337 --format json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") {
				return msg.ErrorMissingApplicationIdArgument
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()
			application, err := client.Get(ctx, application_id)
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

	cmd.Flags().StringVarP(&application_id, "application-id", "a", "", msg.FlagId)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.DescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.DescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)

	return cmd
}

func format(cmd *cobra.Command, application api.EdgeApplicationResponse) ([]byte, error) {

	var b bytes.Buffer

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		file, err := json.MarshalIndent(application, "", " ")
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		b.Write([]byte(fmt.Sprintf("ID: %d\n", uint64(application.GetId()))))
		b.Write([]byte(fmt.Sprintf("Name: %s\n", application.GetName())))
		b.Write([]byte(fmt.Sprintf("Active: %t\n", application.GetActive())))
		b.Write([]byte(fmt.Sprintf("Application Acceleration: %t\n", application.GetApplicationAcceleration())))
		b.Write([]byte(fmt.Sprintf("Caching: %t\n", application.GetCaching())))
		b.Write([]byte(fmt.Sprintf("Delivery Protocol: %s\n", application.GetDeliveryProtocol())))
		b.Write([]byte(fmt.Sprintf("Device Detection: %t\n", application.GetDeviceDetection())))
		b.Write([]byte(fmt.Sprintf("Edge Firewall: %t\n", application.GetEdgeFirewall())))
		b.Write([]byte(fmt.Sprintf("Edge Functions: %t\n", application.GetEdgeFunctions())))
		b.Write([]byte(fmt.Sprintf("Http Port: %d\n", application.GetHttpPort())))
		b.Write([]byte(fmt.Sprintf("Https Port: %d\n", application.GetHttpsPort())))
		b.Write([]byte(fmt.Sprintf("Image Optimization: %t\n", application.GetImageOptimization())))
		b.Write([]byte(fmt.Sprintf("L2 Caching: %t\n", application.GetL2Caching())))
		b.Write([]byte(fmt.Sprintf("Load Balancer: %t\n", application.GetLoadBalancer())))
		b.Write([]byte(fmt.Sprintf("Minimum TLS Version: %s\n", application.GetMinimumTlsVersion())))
		b.Write([]byte(fmt.Sprintf("Raw Logs: %t\n", application.GetRawLogs())))
		b.Write([]byte(fmt.Sprintf("Web Application Firewall: %t\n", application.GetWebApplicationFirewall())))
		return b.Bytes(), nil
	}
}
