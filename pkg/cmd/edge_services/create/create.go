package create

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/MakeNowJust/heredoc"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
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
		Use:           "create [flags]",
		Short:         "Creates a new Edge Service",
		Long:          `Creates a new Edge Service`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services create --name "Hello"
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
					return errmsg.ErrorMandatoryName
				}
				name, err := cmd.Flags().GetString("name")
				if err != nil {
					return errmsg.ErrorInvalidNameFlag
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
	createCmd.Flags().StringVar(&fields.Name, "name", "", "Your Edge Service's name (Mandatory)")
	createCmd.Flags().StringVar(&fields.InPath, "in", "", "Uses provided file path to create an Edge Service. You can use - for reading from stdin")

	return createCmd
}

func createNewService(client *sdk.APIClient, out io.Writer, request sdk.CreateServiceRequest) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.NewService(c).CreateServiceRequest(request).Execute()
	if err != nil {
		if httpResp == nil || httpResp.StatusCode >= 500 {
			err := utils.CheckStatusCode500Error(err)
			return err
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s", errmsg.ErrorCreateService, string(body))
	}

	fmt.Fprintf(out, "Created Edge Service with ID %d\n", resp.Id)

	return nil
}
