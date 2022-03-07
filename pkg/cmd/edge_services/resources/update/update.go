package update

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_services/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

const SHELL_SCRIPT string = "Shell Script"

type Fields struct {
	Name        string
	Trigger     string
	ContentType string
	ContentFile string
	InPath      string
}

type UpdateRequestResource struct {
	sdk.UpdateResourceRequest
	Id int64
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}
	// updateCmd represents the update command
	updateCmd := &cobra.Command{
		Use:           "update <service_id> <resource_id> [flags]",
		Short:         "Updates a Resource",
		Long:          `Updates a Resource based on a resource_id`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services resources update 1234 69420 --name '/tmp/hello.txt'
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) < 2 && !cmd.Flags().Changed("in") {
				return errmsg.ErrorMissingArgumentUpdateResource
			}

			convertedIds := make([]int64, 2)
			request := UpdateRequestResource{}

			if cmd.Flags().Changed("in") {
				var (
					file *os.File
					err  error
				)
				ids, err := utils.ConvertIdsToInt(args[0])
				if err != nil {
					return utils.ErrorConvertingIdArgumentToInt
				}
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
				convertedIds[0] = ids[0]
				convertedIds[1] = request.Id
			} else {

				ids, err := utils.ConvertIdsToInt(args[0], args[1])
				if err != nil {
					return utils.ErrorConvertingIdArgumentToInt
				}

				convertedIds[0] = ids[0]
				convertedIds[1] = ids[1]

				replacer := strings.NewReplacer("shellscript", "Shell Script", "text", "Text", "install", "Install", "reload", "Reload", "uninstall", "Uninstall")

				valueHasChanged := false

				if cmd.Flags().Changed("name") {
					name, err := cmd.Flags().GetString("name")
					if err != nil {
						return errmsg.ErrorInvalidNameFlag
					}
					request.SetName(name)
					valueHasChanged = true
				}

				if cmd.Flags().Changed("trigger") {
					trigger, err := cmd.Flags().GetString("trigger")
					if err != nil {
						return errmsg.ErrorInvalidTriggerFlag
					}
					triggerConverted := replacer.Replace(trigger)
					request.SetTrigger(triggerConverted)
					request.SetContentType(SHELL_SCRIPT)
					valueHasChanged = true
				}

				if cmd.Flags().Changed("content-type") {
					contentType, err := cmd.Flags().GetString("content-type")
					if err != nil {
						return errmsg.ErrorInvalidContentTypeFlag
					}
					contentTypeConverted := replacer.Replace(contentType)
					request.SetContentType(contentTypeConverted)
					valueHasChanged = true
				}

				if cmd.Flags().Changed("content-file") {

					contentPath, err := cmd.Flags().GetString("content-file")
					if err != nil {
						return utils.ErrorHandlingFile
					}

					file, err := ioutil.ReadFile(contentPath)
					if err != nil {
						return utils.ErrorHandlingFile
					}

					stringFile := string(file)

					request.SetContent(stringFile)
					valueHasChanged = true
				}

				if !valueHasChanged {
					return utils.ErrorUpdateNoFlagsSent
				}
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}

			if err := updateResource(client, f.IOStreams.Out, convertedIds[0], convertedIds[1], request, verbose); err != nil {
				return err
			}

			return nil
		},
	}

	updateCmd.Flags().StringVar(&fields.Name, "name", "", "Your Resource's name: <PATH>/<RESOURCE_NAME>")
	updateCmd.Flags().StringVar(&fields.Trigger, "trigger", "", "Your Resource's trigger: <Install|Reload|Uninstall>")
	updateCmd.Flags().StringVar(&fields.ContentType, "content-type", "", "Your Resource's content-type: <shellscript|text>")
	updateCmd.Flags().StringVar(&fields.ContentFile, "content-file", "", "Path to the file containing your Resource's content")
	updateCmd.Flags().StringVar(&fields.InPath, "in", "", "Uses provided file path to update the fields. You can use - for reading from stdin")

	return updateCmd
}

func updateResource(client *sdk.APIClient, out io.Writer, service_id int64, resource_id int64, update UpdateRequestResource, verbose bool) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.PatchServiceResource(c, service_id, resource_id).UpdateResourceRequest(update.UpdateResourceRequest).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s", errmsg.ErrorUpdateResource, string(body))
	}

	if verbose {
		fmt.Fprintf(out, "ID: %d\n", resp.Id)
		fmt.Fprintf(out, "Name: %s\n", resp.Name)
		fmt.Fprintf(out, "Type: %s\n", resp.Type)
		fmt.Fprintf(out, "Content type: %s\n", resp.ContentType)
		fmt.Fprintf(out, "Content: \n")
		fmt.Fprintf(out, "%s", resp.Content)
	}

	return nil
}
