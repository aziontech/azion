package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_services"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var service_id int64
	opts := &contracts.DescribeOptions{}
	// describeCmd represents the describe command
	describeCmd := &cobra.Command{
		Use:           msg.EdgeServiceDescribeUsage,
		Short:         msg.EdgeServiceDescribeShortDescription,
		Long:          msg.EdgeServiceDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion edge_services describe --service-id 4312
		$ azion edge_services describe --service-id 1337 --with-variables
		$ azion edge_services describe --service-id 1337 --format json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("service-id") {
				return msg.ErrorMissingServiceIdArgument
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			withVariables, err := cmd.Flags().GetBool("with-variables")
			if err != nil {
				return msg.ErrorWithVariablesFlag
			}

			service, err := describeService(client, service_id, withVariables)
			if err != nil {
				return err
			}

			out := f.IOStreams.Out
			formattedService, err := format(cmd, service)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedService, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, msg.EdgeServiceFileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedService[:])
				if err != nil {
					return err
				}
			}

			return nil

		},
	}
	describeCmd.Flags().Int64VarP(&service_id, "service-id", "s", 0, msg.EdgeServiceFlagId)
	describeCmd.Flags().Bool("with-variables", false, msg.EdgeServiceDescribeFlagWithVariable)
	describeCmd.Flags().StringVar(&opts.OutPath, "out", "", msg.EdgeServiceFlagOut)
	describeCmd.Flags().StringVar(&opts.Format, "format", "", msg.EdgeServiceFlagFormat)
	describeCmd.Flags().BoolP("help", "h", false, msg.EdgeServiceDescribeHelpFlag)

	return describeCmd

}

func describeService(client *sdk.APIClient, service_id int64, withVariables bool) (*sdk.ServiceResponse, error) {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetService(c, service_id).WithVars(withVariables).Execute()
	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return nil, fmt.Errorf(msg.ErrorGetSerivce.Error(), message)
	}
	return &resp, nil
}

func format(cmd *cobra.Command, service *sdk.ServiceResponse) ([]byte, error) {

	var b bytes.Buffer

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		file, err := json.MarshalIndent(service, "", " ")
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		b.Write([]byte(fmt.Sprintf("ID: %d\n", uint64(service.GetId()))))
		b.Write([]byte(fmt.Sprintf("Name: %s\n", service.GetName())))
		b.Write([]byte(fmt.Sprintf("Active: %t\n", service.GetActive())))
		b.Write([]byte(fmt.Sprintf("Updated at: %s\n", service.GetUpdatedAt())))
		b.Write([]byte(fmt.Sprintf("Last Editor: %s\n", service.GetLastEditor())))
		b.Write([]byte(fmt.Sprintf("Bound Nodes: %d\n", uint64(service.GetBoundNodes()))))
		b.Write([]byte(fmt.Sprintf("Permissions: %s\n", service.GetPermissions())))
		if cmd.Flags().Changed("with-variables") {
			b.Write([]byte("Variables:\n"))
			for _, variable := range *service.Variables {
				b.Write([]byte(fmt.Sprintf(" Name: %s\tValue: %s\n", variable.Name, variable.Value)))
			}
		}
		return b.Bytes(), nil
	}
}
