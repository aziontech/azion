package edgeapplication

import (
	"context"
	"errors"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/delete/edge_application"
	app "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	fun "github.com/aziontech/azion-cli/pkg/api/edge_function"
	store "github.com/aziontech/azion-cli/pkg/cmd/delete/edge_storage/bucket"
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
	clientapp := app.NewClient(del.f.HttpClient, del.f.Config.GetString("api_url"), del.f.Config.GetString("token"))
	clientfunc := fun.NewClient(del.f.HttpClient, del.f.Config.GetString("api_url"), del.f.Config.GetString("token"))
	storagecmd := store.NewBucket(del.f)

	err = clientapp.Delete(ctx, azionJson.Application.ID)
	if err != nil {
		return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
	}

	if azionJson.Function.ID == 0 {
		_, err := fmt.Fprint(del.f.IOStreams.Out, msg.MissingFunction)
		if err != nil {
			return err
		}
	} else {
		err = clientfunc.Delete(ctx, azionJson.Function.ID)
		if err != nil {
			return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
		}
	}

	if azionJson.Bucket != "" {
		storagecmd.SetArgs([]string{"", "--name", azionJson.Bucket, "--force"})
		err := storagecmd.Execute()
		if err != nil {
			return msg.ErrorFailCascadeStorage
		}
	}

	err = del.UpdateJson(del)
	if err != nil {
		return err
	}

	deleteOut := output.GeneralOutput{
		Msg: msg.CascadeSuccess,
		Out: del.f.IOStreams.Out}
	return output.Print(&deleteOut)

}
