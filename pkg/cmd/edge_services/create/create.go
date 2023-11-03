package create

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/edge_services"
	"io"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name   string
	InPath string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	createCmd := &cobra.Command{
		Use:           edgeservices.EdgeServiceCreateUsage,
		Short:         edgeservices.EdgeServiceCreateShortDescription,
		Long:          edgeservices.EdgeServiceCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion edge_services create --name "Hello"
		$ azion edge_services create --in "<path>/create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			serviceRequest := sdk.CreateServiceRequest{}

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
						return fmt.Errorf("%s %s", utils.ErrorOpeningFile, fields.InPath)
					}
				}

				err = cmdutil.UnmarshallJsonFromReader(file, &serviceRequest)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("name") {
					return edgeservices.ErrorMandatoryName
				}
				name, err := cmd.Flags().GetString("name")
				if err != nil {
					return edgeservices.ErrorInvalidNameFlag
				}
				serviceRequest.SetName(name)

			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := createNewService(client, f.IOStreams.Out, serviceRequest); err != nil {
				return err
			}

			return nil
		},
	}
	createCmd.Flags().StringVar(&fields.Name, "name", "", edgeservices.EdgeServiceCreateFlagName)
	createCmd.Flags().StringVar(&fields.InPath, "in", "", edgeservices.EdgeServiceCreateFlagIn)
	createCmd.Flags().BoolP("help", "h", false, edgeservices.EdgeServiceCreateFlagHelp)

	return createCmd
}

func createNewService(client *sdk.APIClient, out io.Writer, request sdk.CreateServiceRequest) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.NewService(c).CreateServiceRequest(request).Execute()
	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf(edgeservices.ErrorCreateService.Error(), message)
	}

	fmt.Fprintf(out, edgeservices.EdgeServiceCreateOutputSuccess, resp.Id)

	return nil
}
