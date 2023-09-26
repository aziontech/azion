package describe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type Fields struct {
	ServiceId  int64
	ResourceId int64
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}
	opts := &contracts.DescribeOptions{}
	// describeCmd represents the describe command
	describeCmd := &cobra.Command{
		Use:           msg.EdgeServiceResourceDescribeUsage,
		Short:         msg.EdgeServiceResourceDescribeShortDescription,
		Long:          msg.EdgeServiceResourceDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion edge_services resources describe --service-id 1234 --resource-id 80312
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("service-id") || !cmd.Flags().Changed("resource-id") {
				return msg.ErrorMissingResourceIdArgument
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			resource, err := describeResource(client, f.IOStreams.Out, fields.ServiceId, fields.ResourceId)
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
				fmt.Fprintf(out, msg.EdgeServiceFileWritten, filepath.Clean(opts.OutPath))
			} else {
				_, err := out.Write(formattedResource[:])
				if err != nil {
					return err
				}
			}

			return nil

		},
	}

	describeCmd.Flags().Int64VarP(&fields.ServiceId, "service-id", "s", 0, msg.EdgeServiceFlagId)
	describeCmd.Flags().Int64VarP(&fields.ResourceId, "resource-id", "r", 0, msg.EdgeServiceResourceFlagId)
	describeCmd.Flags().StringVar(&opts.OutPath, "out", "", msg.EdgeServiceFlagOut)
	describeCmd.Flags().StringVar(&opts.Format, "format", "", msg.EdgeServiceFlagFormat)
	describeCmd.Flags().BoolP("help", "h", false, msg.EdgeServiceResourceDescribeFlagHelp)

	return describeCmd

}

func describeResource(client *sdk.APIClient, out io.Writer, service_id int64, resource_id int64) (*sdk.ResourceDetail, error) {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetResource(c, service_id, resource_id).Execute()
	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return nil, fmt.Errorf(msg.ErrorGetResource.Error(), message)
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
