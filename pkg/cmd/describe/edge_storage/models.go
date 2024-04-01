package edge_storage

import (
	"github.com/aziontech/azion-cli/pkg/cmdutil"
)

type Fields struct {
	BucketName string `json:"bucket-name"`
	ObjectKey  string `json:"object-key"`
	Factory    *cmdutil.Factory
}
