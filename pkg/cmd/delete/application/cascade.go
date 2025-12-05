package application

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/delete/application"
	app "github.com/aziontech/azion-cli/pkg/api/applications"
	fun "github.com/aziontech/azion-cli/pkg/api/function"
	workload "github.com/aziontech/azion-cli/pkg/api/workloads"
	store "github.com/aziontech/azion-cli/pkg/cmd/delete/storage/bucket"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
)

func CascadeDelete(ctx context.Context, del *DeleteCmd) error {
	logger.FInfo(del.f.IOStreams.Out, "Cascade deleting resources\n")
	azionJson, err := del.GetAzion(ProjectConf)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return msg.ErrorMissingAzionJson
		} else {
			return err
		}
	}
	if azionJson.Application.ID == 0 {
		return msg.ErrorMissingApplicationIdJson
	}

	// Initialize clients
	clientapp := app.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientfunc := fun.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	clientworkload := workload.NewClient(del.f.HttpClient, del.f.Config.GetString("api_v4_url"), del.f.Config.GetString("token"))
	storagecmd := store.NewBucket(del.f)

	// Collect all errors
	var errs []string

	// Delete workload first if it exists
	if azionJson.Workloads.Id != 0 {
		logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Deleting workload with ID %d\n", azionJson.Workloads.Id))
		err = clientworkload.Delete(ctx, azionJson.Workloads.Id)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failed to delete workload: %v", err))
			logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Failed to delete workload: %v\n", err))
		}
	}

	// Delete edge application
	logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Deleting edge application with ID %d\n", azionJson.Application.ID))
	err = clientapp.Delete(ctx, azionJson.Application.ID)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Failed to delete application: %v", err))
		logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Failed to delete application: %v\n", err))
	}

	// Delete functions
	for _, funcJson := range azionJson.Function {
		if funcJson.ID == 0 {
			_, err := fmt.Fprint(del.f.IOStreams.Out, msg.MissingFunction)
			if err != nil {
				errs = append(errs, fmt.Sprintf("Failed to print missing function message: %v", err))
			}
		} else {
			logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Deleting function with ID %d\n", funcJson.ID))
			err = clientfunc.Delete(ctx, funcJson.ID)
			if err != nil {
				errs = append(errs, fmt.Sprintf("Failed to delete function: %v", err))
				logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Failed to delete function: %v\n", err))
			}
		}
	}

	// Delete bucket if it exists
	if azionJson.Bucket != "" {
		logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Deleting storage bucket %s\n", azionJson.Bucket))
		storagecmd.SetArgs([]string{"", "--name", azionJson.Bucket, "--force"})
		err := storagecmd.Execute()
		if err != nil {
			errs = append(errs, fmt.Sprintf("Failed to delete storage bucket: %v", err))
			logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Failed to delete storage bucket: %v\n", err))
		}
	}

	// Update JSON file
	err = del.UpdateJson(del)
	if err != nil {
		errs = append(errs, fmt.Sprintf("Failed to update JSON file: %v", err))
		logger.FInfo(del.f.IOStreams.Out, fmt.Sprintf("Failed to update JSON file: %v\n", err))
	}

	// Check if there were any errors
	if len(errs) > 0 {
		return fmt.Errorf("Cascade delete completed with errors:\n%s", strings.Join(errs, "\n"))
	}

	deleteOut := output.GeneralOutput{
		Msg: msg.CascadeSuccess,
		Out: del.f.IOStreams.Out}
	return output.Print(&deleteOut)
}
