package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {

	opts := &contracts.DescribeOptions{}
	// describeCmd represents the describe command
	describeCmd := &cobra.Command{
		Use:           "describe <service_id> [flags]",
		Short:         "Describes an Edge Service",
		Long:          `Details an Edge Service based on the id given`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services describe 4312
        $ azioncli edge_services describe 1337 --with-variables
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errmsg.ErrorMissingServiceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			withVariables, err := cmd.Flags().GetBool("with-variables")
			if err != nil {
				return errmsg.ErrorWithVariablesFlag
			}

			service, err := describeService(client, ids[0], withVariables)
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
				fmt.Fprintf(out, "File successfully written to: %s\n", filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedService[:])
				if err != nil {
					return err
				}
			}

			return nil

		},
	}
	describeCmd.Flags().Bool("with-variables", false, "Displays the Edge Service's variables (disabled by default)")
	describeCmd.Flags().StringVar(&opts.OutPath, "out", "", "Exports the command result to the received file path")
	describeCmd.Flags().StringVar(&opts.Format, "format", "", "You can change the results format by passing json value to this flag")

	return describeCmd

}

func describeService(client *sdk.APIClient, service_id int64, withVariables bool) (*sdk.ServiceResponse, error) {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetService(c, service_id).WithVars(withVariables).Execute()
	if err != nil {
		errMsg := utils.ErrorPerStatusCode(httpResp, err)

		return nil, fmt.Errorf("%w: %s", errmsg.ErrorGetSerivce, errMsg)
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
