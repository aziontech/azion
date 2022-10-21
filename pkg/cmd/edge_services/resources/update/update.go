package update

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_services"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeservices"
	"github.com/spf13/cobra"
)

const SHELL_SCRIPT string = "Shell Script"

type Fields struct {
	ServiceId   int64
	ResourceId  int64
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
		Use:           msg.EdgeServiceResourceUpdateUsage,
		Short:         msg.EdgeServiceResourceUpdateShortDescription,
		Long:          msg.EdgeServiceResourceUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli edge_services resources update --service-id 1234 --resource-id 69420 --name '/tmp/hello.txt'
		$ azioncli edge_services resources update --service-id 1234 --resource-id 69420 --name "/tmp/my_script.sh" --content-type shellscript --content-file "./text.txt"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("service-id") || !cmd.Flags().Changed("resource-id") {
				return msg.ErrorMissingArgumentUpdateResource
			}

			request := UpdateRequestResource{}

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
			} else {

				replacer := strings.NewReplacer("shellscript", "Shell Script", "text", "Text", "install", "Install", "reload", "Reload", "uninstall", "Uninstall")

				valueHasChanged := false

				if cmd.Flags().Changed("name") {
					name, err := cmd.Flags().GetString("name")
					if err != nil {
						return msg.ErrorInvalidNameFlag
					}
					request.SetName(name)
					valueHasChanged = true
				}

				if cmd.Flags().Changed("trigger") {
					trigger, err := cmd.Flags().GetString("trigger")
					if err != nil {
						return msg.ErrorInvalidTriggerFlag
					}
					triggerConverted := replacer.Replace(trigger)
					request.SetTrigger(triggerConverted)
					request.SetContentType(SHELL_SCRIPT)
					valueHasChanged = true
				}

				if cmd.Flags().Changed("content-type") {
					contentType, err := cmd.Flags().GetString("content-type")
					if err != nil {
						return msg.ErrorInvalidContentTypeFlag
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

			if err := updateResource(client, f.IOStreams.Out, fields.ServiceId, fields.ResourceId, request); err != nil {
				return err
			}

			return nil
		},
	}

	updateCmd.Flags().Int64VarP(&fields.ServiceId, "service-id", "s", 0, msg.EdgeServiceFlagId)
	updateCmd.Flags().Int64VarP(&fields.ResourceId, "resource-id", "r", 0, msg.EdgeServiceResourceFlagId)
	updateCmd.Flags().StringVar(&fields.Name, "name", "", msg.EdgeServiceResourceUpdateFlagName)
	updateCmd.Flags().StringVar(&fields.Trigger, "trigger", "", msg.EdgeServiceResourceUpdateFlagTrigger)
	updateCmd.Flags().StringVar(&fields.ContentType, "content-type", "", msg.EdgeServiceResourceUpdateFlagContentType)
	updateCmd.Flags().StringVar(&fields.ContentFile, "content-file", "", msg.EdgeServiceResourceUpdateFlagContentFile)
	updateCmd.Flags().StringVar(&fields.InPath, "in", "", msg.EdgeServiceResourceUpdateFlagIn)
	updateCmd.Flags().BoolP("help", "h", false, msg.EdgeServiceResourceUpdateFlagHelp)

	return updateCmd
}

func updateResource(client *sdk.APIClient, out io.Writer, service_id int64, resource_id int64, update UpdateRequestResource) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.PatchServiceResource(c, service_id, resource_id).UpdateResourceRequest(update.UpdateResourceRequest).Execute()
	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf(msg.ErrorUpdateResource.Error(), message)
	}

	fmt.Fprintf(out, msg.EdgeServiceResourceUpdateOutputSuccess, resp.Id)

	return nil
}
