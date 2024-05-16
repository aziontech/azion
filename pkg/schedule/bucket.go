package schedule

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/cmdutil"

	api "github.com/aziontech/azion-cli/pkg/api/storage"
)

type DeleteBucket struct {
	Factory *cmdutil.Factory
	Name    string
}

func (b DeleteBucket) TriggerEvent() error {
	client := api.NewClient(
		b.Factory.HttpClient,
		b.Factory.Config.GetString("storage_url"),
		b.Factory.Config.GetString("token"))
	ctx := context.Background()
	return client.DeleteBucket(ctx, b.Name)
}
