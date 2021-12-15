package create

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/aziontech/azion-cli/cmd/edge_services/requests"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/edgeservices-go-sdk"
	"github.com/spf13/cobra"
)

const SHELL_SCRIPT string = "Shell Script"

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// createCmd represents the create command
	createCmd := &cobra.Command{
		Use:           "create",
		Short:         "Creates a new resource",
		Long:          `Creates a new resource in an Edge Service based on a giver servce_id.`,
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

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			trigger, err := cmd.Flags().GetString("trigger")
			if err != nil {
				return err
			}

			content_type, err := cmd.Flags().GetString("content-type")
			if err != nil {
				return err
			}
			if content_type == SHELL_SCRIPT {
				if trigger == "" {
					return utils.ErrorInvalidResourceTrigger
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

			client, err := requests.CreateClient(f, cmd)
			if err != nil {
				return err
			}

			if err := createNewResource(client, ids[0], name, trigger, content_type, stringFile); err != nil {
				return fmt.Errorf("%v. %v", err, utils.GenericUseHelp)
			}

			return nil
		},
	}

	createCmd.Flags().StringP("name", "n", "", "<PATH>/<RESOURCE_NAME>")
	_ = createCmd.MarkFlagRequired("name")
	createCmd.Flags().String("trigger", "", "<Install|Reload|Uninstall>")
	createCmd.Flags().String("content-type", "", "<\"Shell Script\"|\"Text\"")
	_ = createCmd.MarkFlagRequired("content-type")
	createCmd.Flags().String("content-file", "", "Absolute path to where the file with the content is located at")
	_ = createCmd.MarkFlagRequired("content-file")

	return createCmd
}

func createNewResource(client *sdk.APIClient, service_id int64, name string, trigger string, content_type string, file string) error {
	c := context.Background()
	api := client.DefaultApi

	request := sdk.CreateResourceRequest{
		Name:        name,
		Trigger:     trigger,
		ContentType: content_type,
		Content:     file,
	}

	resp, httpResp, err := api.PostResource(c, service_id).CreateResourceRequest(request).Execute()
	if err != nil {
		if httpResp.StatusCode >= 500 {
			return utils.ErrorInternalServerError
		}
		body, err := ioutil.ReadAll(httpResp.Body)
		if err != nil {
			return err
		}

		return errors.New(string(body))
	}

	fmt.Println(resp.Name)

	return nil
}
