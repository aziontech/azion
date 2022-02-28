package create

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
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

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// createCmd represents the create command
	createCmd := &cobra.Command{
		Use:           "create <service_id> [flags]",
		Short:         "Creates a new Resource",
		Long:          `Creates a new Resource in an Edge Service based on a given service_id.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_services resources create 1234 --name "/tmp/test.txt" --content-type text --content-file "./text.txt"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errmsg.ErrorMissingResourceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			replacer := strings.NewReplacer("shellscript", "Shell Script", "text", "Text", "install", "Install", "reload", "Reload", "uninstall", "Uninstall")

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			trigger, err := cmd.Flags().GetString("trigger")
			triggerConverted := replacer.Replace(trigger)
			if err != nil {
				return err
			}

			contentType, err := cmd.Flags().GetString("content-type")
			if err != nil {
				return err
			}
			contentTypeConverted := replacer.Replace(contentType)
			if contentTypeConverted == SHELL_SCRIPT {
				if trigger == "" {
					return errmsg.ErrorInvalidResourceTrigger
				}
			}

			contentPath, err := cmd.Flags().GetString("content-file")
			if err != nil {
				return utils.ErrorHandlingFile
			}

			file, err := ioutil.ReadFile(contentPath)
			if err != nil {
				return utils.ErrorHandlingFile
			}

			stringFile := string(file)

			client, err := requests.CreateClient(f)
			if err != nil {
				return err
			}

			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return err
			}
			if err := createNewResource(client, f.IOStreams.Out, ids[0], name, triggerConverted, contentTypeConverted, stringFile, verbose); err != nil {
				return err
			}

			return nil
		},
	}

	createCmd.Flags().String("name", "", "Name of your Resource: <PATH>/<RESOURCE_NAME> (Mandatory)")
	_ = createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("trigger", "", "Trigger of your Resource: <Install|Reload|Uninstall>")
	createCmd.Flags().String("content-type", "", "Content-type of your Resource: <shellscript|text> (Mandatory)")
	_ = createCmd.MarkFlagRequired("content-type")
	createCmd.Flags().String("content-file", "", "Absolute path to where the file with the content is located at (Mandatory)")
	_ = createCmd.MarkFlagRequired("content-file")

	return createCmd
}

func createNewResource(client *sdk.APIClient, out io.Writer, service_id int64, name string, trigger string, contentType string, file string, verbose bool) error {
	c := context.Background()
	api := client.DefaultApi

	request := sdk.CreateResourceRequest{
		Name:        name,
		Trigger:     trigger,
		ContentType: contentType,
		Content:     file,
	}

	resp, httpResp, err := api.PostResource(c, service_id).CreateResourceRequest(request).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("%w: %s", errmsg.ErrorCreateResource, string(body))
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
