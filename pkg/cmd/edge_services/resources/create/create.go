package create

import (
	"context"
	"fmt"
	"io"
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
	Name        string
	Trigger     string
	ContentType string
	ContentFile string
	InPath      string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	// createCmd represents the create command
	createCmd := &cobra.Command{
		Use:           msg.EdgeServiceResourceCreateUsage,
		Short:         msg.EdgeServiceResourceCreateShortDescription,
		Long:          msg.EdgeServiceResourceCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion edge_services resources create --service-id 1234 --name "/tmp/test.txt" --content-type text --content-file "./text.txt"
		$ azion edge_services resources create --service-id 1234 --name "/tmp/my_script.sh" --content-type shellscript --content-file "./text.txt"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("service-id") {
				return msg.ErrorMissingServiceIdArgument
			}

			request := sdk.CreateResourceRequest{}

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

				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("name") || !cmd.Flags().Changed("content-file") || !cmd.Flags().Changed("content-type") {
					return msg.ErrorMandatoryFlagsResource
				}

				replacer := strings.NewReplacer("shellscript", "Shell Script", "text", "Text", "install", "Install", "reload", "Reload", "uninstall", "Uninstall")

				name, err := cmd.Flags().GetString("name")
				if err != nil {
					return err
				}
				request.SetName(name)

				trigger, err := cmd.Flags().GetString("trigger")
				triggerConverted := replacer.Replace(trigger)
				if err != nil {
					return err
				}
				request.SetTrigger(triggerConverted)

				contentType, err := cmd.Flags().GetString("content-type")
				if err != nil {
					return err
				}
				contentTypeConverted := replacer.Replace(contentType)
				if contentTypeConverted == SHELL_SCRIPT {
					if trigger == "" {
						return msg.ErrorInvalidResourceTrigger
					}
				}
				request.SetContentType(contentTypeConverted)

				contentPath, err := cmd.Flags().GetString("content-file")
				if err != nil {
					return utils.ErrorHandlingFile
				}

				file, err := os.ReadFile(contentPath)
				if err != nil {
					return utils.ErrorHandlingFile
				}

				stringFile := string(file)
				if stringFile == "" {
					return utils.ErrorEmptyFile
				}
				request.SetContent(stringFile)
			}

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			if err := createNewResource(client, f.IOStreams.Out, fields.ServiceId, request); err != nil {
				return err
			}

			return nil
		},
	}

	createCmd.Flags().Int64VarP(&fields.ServiceId, "service-id", "s", 0, msg.EdgeServiceFlagId)
	createCmd.Flags().StringVar(&fields.Name, "name", "", msg.EdgeServiceResourceCreateFlagName)
	createCmd.Flags().StringVar(&fields.Trigger, "trigger", "", msg.EdgeServiceResourceCreateFlagTrigger)
	createCmd.Flags().StringVar(&fields.ContentType, "content-type", "", msg.EdgeServiceResourceCreateFlagContentType)
	createCmd.Flags().StringVar(&fields.ContentFile, "content-file", "", msg.EdgeServiceResourceCreateFlagContentFile)
	createCmd.Flags().StringVar(&fields.InPath, "in", "", msg.EdgeServiceResourceCreateFlagIn)
	createCmd.Flags().BoolP("help", "h", false, msg.EdgeServiceResourceCreateFlagHelp)

	return createCmd
}

func createNewResource(client *sdk.APIClient, out io.Writer, service_id int64, request sdk.CreateResourceRequest) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.PostResource(c, service_id).CreateResourceRequest(request).Execute()
	if err != nil {
		message := utils.ErrorPerStatusCode(httpResp, err)

		return fmt.Errorf(msg.ErrorCreateResource.Error(), message)
	}

	fmt.Fprintf(out, msg.EdgeServiceResourceCreateOutputSuccess, resp.Id)

	return nil
}
