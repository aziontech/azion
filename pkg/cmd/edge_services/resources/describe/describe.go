package describe

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/MakeNowJust/heredoc"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {

	// describeCmd represents the describe command
	describeCmd := &cobra.Command{
		Use:           "describe <service_id> <resource_id> [flags]",
		Short:         "Describes a Resource",
		Long:          `Provides a long description of a Resource based on a given service_id and a resource_id`,
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

			if err := describeResource(client, f.IOStreams.Out, ids[0], ids[1]); err != nil {
				return err
			}

			return nil

		},
	}
	return describeCmd

}

func describeResource(client *sdk.APIClient, out io.Writer, service_id int64, resource_id int64) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetResource(c, service_id, resource_id).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s", errmsg.ErrorGetResource, string(body))
	}

	fmt.Fprintf(out, "ID: %d\n", resp.Id)
	fmt.Fprintf(out, "Name: %s\n", resp.Name)
	fmt.Fprintf(out, "Type: %s\n", resp.Type)
	fmt.Fprintf(out, "Content type: %s\n", resp.ContentType)
	fmt.Fprintf(out, "Content: \n")
	fmt.Fprintf(out, "%s", resp.Content)

	return nil
}
