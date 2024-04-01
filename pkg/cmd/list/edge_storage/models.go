package edge_storage

import (
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
)

type Objects struct {
	BucketName string
	Options    *contracts.ListOptions
	Factory    *cmdutil.Factory
}

type Bucket struct {
	Options *contracts.ListOptions
	Factory *cmdutil.Factory
}
