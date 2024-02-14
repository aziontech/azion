package edgeapplication

import (
	"context"
	"errors"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/delete/edge_application"
	app "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	fun "github.com/aziontech/azion-cli/pkg/api/edge_function"
)

func (del *DeleteCmd) Cascade(ctx context.Context) error {
	azionJson, err := del.GetAzion()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return msg.ErrorMissingAzionJson
		} else {
			return err //message with error
		}
	}
	if azionJson.Application.ID == 0 {
		return msg.ErrorMissingApplicationIdJson
	}
	clientapp := app.NewClient(del.f.HttpClient, del.f.Config.GetString("api_url"), del.f.Config.GetString("token"))
	clientfunc := fun.NewClient(del.f.HttpClient, del.f.Config.GetString("api_url"), del.f.Config.GetString("token"))

	err = clientapp.Delete(ctx, azionJson.Application.ID)
	if err != nil {
		return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
	}

	if azionJson.Function.ID == 0 {
		fmt.Fprintf(del.f.IOStreams.Out, msg.MissingFunction)
	} else {
		err = clientfunc.Delete(ctx, azionJson.Function.ID)
		if err != nil {
			return fmt.Errorf(msg.ErrorFailToDeleteApplication.Error(), err)
		}
	}

	fmt.Fprintf(del.f.IOStreams.Out, "%s\n", msg.CascadeSuccess)

	err = del.UpdateJson(del)
	if err != nil {
		return err
	}

	return nil
}
