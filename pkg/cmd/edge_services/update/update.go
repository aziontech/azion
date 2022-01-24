package update

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// listCmd represents the list command
	updateCmd := &cobra.Command{
		Use:           "update <service_id> [flags]",
		Short:         "Updates parameters of an edge service",
		Long:          `Updates parameters of an edge service`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return utils.ErrorMissingServiceIdArgument
			}

			id, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := updateService(client, f.IOStreams.Out, id[0], cmd, args); err != nil {
				return err
			}

			return nil
		},
	}
	updateCmd.Flags().String("name", "", "<EDGE_SERVICE_NAME>")
	updateCmd.Flags().String("active", "", "<true|false>")
	updateCmd.Flags().String("variables-file", "", `<VARIABLES_FILE_PATH>
The format accepted for variables definition is one <KEY>=<VALUE> per line`)

	return updateCmd
}

func updateService(client *sdk.APIClient, out io.Writer, id int64, cmd *cobra.Command, args []string) error {
	c := context.Background()
	api := client.DefaultApi

	serviceRequest := sdk.UpdateServiceRequest{}

	nameHasChanged := cmd.Flags().Changed("name")
	if nameHasChanged {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		serviceRequest.SetName(name)
	}

	activeHasChanged := cmd.Flags().Changed("active")
	if activeHasChanged {
		activeStr, err := cmd.Flags().GetString("active")
		if err != nil {
			return err
		}

		active, err := strconv.ParseBool(activeStr)
		if err != nil {
			return utils.ErrorConvertingStringToBool
		}
		serviceRequest.SetActive(active)
	}

	variablesHasChanged := cmd.Flags().Changed("variables-file")
	if variablesHasChanged {
		variablesPath, err := cmd.Flags().GetString("variables-file")
		if err != nil {
			return utils.ErrorHandlingFile
		}

		file, err := os.Open(variablesPath)
		if err != nil {
			return utils.ErrorHandlingFile
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		reName := regexp.MustCompile("^([^= ]+) *= *(.*)")
		reValue := regexp.MustCompile("^([^= ]+) *= *(.+)")
		v := []sdk.Variable{}
		for scanner.Scan() {
			line := scanner.Text()
			entry := strings.Split(line, "=")
			if len(entry) < 2 {
				return utils.ErrorInvalidVariablesFileFormat
			}
			varName := reName.FindStringSubmatch(strings.Trim(line, " "))[1]
			varValue := reValue.FindStringSubmatch(strings.Trim(line, " "))[2]
			variable := sdk.NewVariable(varName, varValue)
			v = append(v, *variable)
		}
		serviceRequest.SetVariables(v)

		if err := scanner.Err(); err != nil {
			return err
		}

	}

	resp, httpResp, err := api.PatchService(c, id).UpdateServiceRequest(serviceRequest).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}

		return err
	}

	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		return err
	}

	if verbose {
		fmt.Fprintf(out, "ID: %d\n", resp.Id)
		fmt.Fprintf(out, "Name: %s\n", resp.Name)
		fmt.Fprintf(out, "Updated at: %s\n", resp.UpdatedAt)
		fmt.Fprintf(out, "Last Editor: %s\n", resp.LastEditor)
		fmt.Fprintf(out, "Active: %t\n", resp.Active)
		fmt.Fprintf(out, "Bound Nodes: %d\n", resp.BoundNodes)
		fmt.Fprintf(out, "Permissions: %s\n", resp.Permissions)
	}

	return nil
}
