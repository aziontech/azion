package edgeapplication

import (
	"context"
	"errors"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/delete/edge_application"
	"github.com/aziontech/azion-cli/pkg/output"
	app "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"
	fun "github.com/aziontech/azion-cli/pkg/v3api/edge_function"
)

func CascadeDelete(ctx context.Context, del *DeleteCmd) error {
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

	err = clientapp.Delete(ctx, azionJson.Application.ID)
	if err != nil {
		return fmt.Errorf("%v: %w", msg.ErrorFailToDeleteApplication, err)
	}

	if azionJson.Function.ID == 0 {
		fmt.Fprint(del.f.IOStreams.Out, msg.MissingFunction)
	} else {
		err = clientfunc.Delete(ctx, azionJson.Function.ID)
		if err != nil {
			return fmt.Errorf("%v: %w", msg.ErrorFailToDeleteApplication, err)
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
