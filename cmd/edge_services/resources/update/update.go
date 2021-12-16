package update

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
	// updateCmd represents the update command
	updateCmd := &cobra.Command{
		Use:           "update",
		Short:         "Updates a resource based on a resource_id",
		Long:          `Updates fields of a resource based on a service_id and resource_id`,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) < 2 {
				return utils.ErrorMissingResourceIdArgument
			}

			ids, err := utils.ConvertIdsToInt(args[0], args[1])
			if err != nil {
				return utils.ErrorConvertingIdArgumentToInt
			}

			updateRequest := sdk.UpdateResourceRequest{}
			valueHasChanged := false

			if cmd.Flags().Changed("name") {
				name, err := cmd.Flags().GetString("name")
				if err != nil {
					return err
				}
				updateRequest.SetName(name)
				valueHasChanged = true
			}

			if cmd.Flags().Changed("trigger") {
				trigger, err := cmd.Flags().GetString("trigger")
				if err != nil {
					return err
				}
				updateRequest.SetTrigger(trigger)
				updateRequest.SetContentType(SHELL_SCRIPT)
				valueHasChanged = true
			}

			if cmd.Flags().Changed("content-type") {
				contentType, err := cmd.Flags().GetString("content-type")
				if err != nil {
					return err
				}
				updateRequest.SetContentType(contentType)
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

				updateRequest.SetContent(stringFile)
				valueHasChanged = true
			}

			if !valueHasChanged {
				return utils.ErrorUpdateNoFlagsSent
			}

			client, err := requests.CreateClient(f, cmd)
			if err != nil {
				return err
			}

			if err := updateResource(client, ids[0], ids[1], updateRequest); err != nil {
				return fmt.Errorf("%v. %v", err, utils.GenericUseHelp)
			}

			return nil
		},
	}

	updateCmd.Flags().StringP("name", "n", "", "<PATH>/<RESOURCE_NAME>")
	updateCmd.Flags().String("trigger", "", "<Install|Reload|Uninstall>")
	updateCmd.Flags().String("content-type", "", "<\"Shell Script\"|\"Text\"")
	updateCmd.Flags().String("content-file", "", "Absolute path to where the file with the content is located at")

	return updateCmd
}

func updateResource(client *sdk.APIClient, service_id int64, resource_id int64, update sdk.UpdateResourceRequest) error {
	c := context.Background()
	api := client.DefaultApi

	resp, httpResp, err := api.PatchServiceResource(c, service_id, resource_id).UpdateResourceRequest(update).Execute()
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
