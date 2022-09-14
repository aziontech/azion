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

	"github.com/MakeNowJust/heredoc"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

type Fields struct {
	Id        int64
	Name      string
	Active    string
	Variables string
	InPath    string
}

type UpdateRequestService struct {
	sdk.UpdateServiceRequest
	Id int64
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}
	// listCmd represents the list command
	updateCmd := &cobra.Command{
		Use:           "update --service-id <service_id> [flags]",
		Short:         "Updates an Edge Service",
		Long:          `Updates an Edge Service`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services update --service-id 1234 --name 'Hello'
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			// either id parameter or in path should be passed
			if !cmd.Flags().Changed("service-id") && !cmd.Flags().Changed("in") {
				return errmsg.ErrorMissingArgumentUpdate
			}

			request := UpdateRequestService{}

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
						return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.InPath)
					}
				}
				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
				fields.Id = request.Id
			} else {

				nameHasChanged := cmd.Flags().Changed("name")
				if nameHasChanged {
					name, err := cmd.Flags().GetString("name")
					if err != nil {
						return err
					}
					request.SetName(name)
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
					request.SetActive(active)
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
					request.SetVariables(v)

					if err := scanner.Err(); err != nil {
						return err
					}

				}
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := updateService(client, f.IOStreams.Out, fields.Id, cmd, request); err != nil {
				return err
			}

			return nil
		},
	}
	updateCmd.Flags().Int64VarP(&fields.Id, "service-id", "s", 0, "Unique identifier of the Edge Service")
	updateCmd.Flags().StringVar(&fields.Name, "name", "", "Your Edge Service's name")
	updateCmd.Flags().StringVar(&fields.Active, "active", "", "Whether your Edge Service should be active or not: <true|false>")
	updateCmd.Flags().StringVar(&fields.Variables, "variables-file", "", `Path to the file containing your Edge Service's Variables.
The accepted format for defining variables is one <KEY>=<VALUE> per line`)
	updateCmd.Flags().StringVar(&fields.InPath, "in", "", "Uses provided file path to update the fields. You can use - for reading from stdin")

	return updateCmd
}

func updateService(client *sdk.APIClient, out io.Writer, id int64, cmd *cobra.Command, request UpdateRequestService) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.PatchService(c, id).UpdateServiceRequest(request.UpdateServiceRequest).Execute()
	if err != nil {
		errMsg := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf("%w: %s", errmsg.ErrorUpdateService, errMsg)
	}

	fmt.Fprintf(out, "Updated Edge Service with ID %d\n", resp.Id)

	return nil
}
