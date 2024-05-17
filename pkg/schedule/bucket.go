package schedule

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/cmdutil"

	api "github.com/aziontech/azion-cli/pkg/api/storage"
)

const DELETE_BUCKET = "DeleteBucket" 

func TriggerDeleteBucket(f *cmdutil.Factory, name string) error {
	client := api.NewClient(
		f.HttpClient,
		f.Config.GetString("storage_url"),
		f.Config.GetString("token"))
	ctx := context.Background()
	return client.DeleteBucket(ctx, name)
}
