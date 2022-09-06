package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
		Use:           "describe <service_id> <resource_id> [flags]",
		Short:         "Describes a Resource",
		Long:          `Provides a long description of a Resource based on a service_id and a resource_id given`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services resources describe 1234 80312
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errmsg.ErrorMissingResourceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0], args[1])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			resource, err := describeResource(client, f.IOStreams.Out, ids[0], ids[1])
			if err != nil {
				return err
			}

			out := f.IOStreams.Out
			formattedResource, err := format(cmd, resource)
			if err != nil {
				return utils.ErrorFormatOut
			}

			if cmd.Flags().Changed("out") {
				err := cmdutil.WriteDetailsToFile(formattedResource, opts.OutPath, out)
				if err != nil {
					return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
				}
				fmt.Fprintf(out, "File successfully written to: %s\n", filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedResource[:])
				if err != nil {
					return err
				}
			}

			return nil

		},
	}

	describeCmd.Flags().StringVar(&opts.OutPath, "out", "", "Exports the command result to the received file path")
	describeCmd.Flags().StringVar(&opts.Format, "format", "", "You can change the results format by passing json value to this flag")

	return describeCmd

}

func describeResource(client *sdk.APIClient, out io.Writer, service_id int64, resource_id int64) (*sdk.ResourceDetail, error) {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetResource(c, service_id, resource_id).Execute()
	if err != nil {
		errMsg := utils.ErrorPerStatusCode(httpResp, err)

		return nil, fmt.Errorf("%w: %s", errmsg.ErrorGetResource, errMsg)
	}

	return &resp, nil
}

func format(cmd *cobra.Command, resource *sdk.ResourceDetail) ([]byte, error) {

	var b bytes.Buffer

	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}

	if format == "json" || cmd.Flags().Changed("out") {
		file, err := json.MarshalIndent(resource, "", " ")
		if err != nil {
			return nil, err
		}
		return file, nil
	} else {
		b.Write([]byte(fmt.Sprintf("ID: %d\n", uint64(resource.GetId()))))
		b.Write([]byte(fmt.Sprintf("Name: %s\n", resource.GetName())))
		b.Write([]byte(fmt.Sprintf("Trigger: %s\n", resource.GetTrigger())))
		b.Write([]byte(fmt.Sprintf("Content type: %s\n", resource.GetContentType())))
		b.Write([]byte("Content: \n"))
		b.Write([]byte(resource.GetContent()))
		return b.Bytes(), nil

	}
}
