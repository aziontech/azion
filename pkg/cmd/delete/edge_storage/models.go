package edge_storage

import "github.com/aziontech/azion-cli/pkg/cmdutil"

type bucket struct {
	name    string
	force   bool
	factory *cmdutil.Factory
}

type object struct {
	bucketName string
	objectKey  string
	factory    *cmdutil.Factory
}
