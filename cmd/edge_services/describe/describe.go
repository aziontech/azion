package describe

import (
	"context"
	"fmt"
	"io"

	"github.com/aziontech/azion-cli/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {

	// describeCmd represents the describe command
	describeCmd := &cobra.Command{
		Use:           "describe <service_id> [flags]",
		Short:         "Describes a service based on a given service_id",
		Long:          `Describes a service based on a given service_id`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return utils.ErrorMissingServiceIdArgument
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
				return err
			}

			if err := describeService(client, f.IOStreams.Out, ids[0], withVariables); err != nil {
				return err
			}

			return nil

		},
	}
	describeCmd.Flags().Bool("with-variables", false, "")

	return describeCmd

}

func describeService(client *sdk.APIClient, out io.Writer, service_id int64, withVariables bool) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.GetService(c, service_id).WithVars(withVariables).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		return err
	}

	fmt.Fprintf(out, "ID: %d\n", resp.Id)
	fmt.Fprintf(out, "Name: %s\n", resp.Name)
	fmt.Fprintf(out, "Updated at: %s\n", resp.UpdatedAt)
	fmt.Fprintf(out, "Last Editor: %s\n", resp.LastEditor)
	fmt.Fprintf(out, "Active: %t\n", resp.Active)
	fmt.Fprintf(out, "Bound Nodes: %d\n", resp.BoundNodes)
	fmt.Fprintf(out, "Permissions: %s\n", resp.Permissions)
	if withVariables {
		fmt.Fprint(out, "Variables:\n")
		for _, variable := range *resp.Variables {
			fmt.Fprintf(out, " Name: %s\tValue: %s\n", variable.Name, variable.Value)
		}
	}
	return nil
}
