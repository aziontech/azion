package edge_storage

import "github.com/aziontech/azion-cli/pkg/cmdutil"

type bucket struct {
	name            string
	workloadsAccess string
	fileJSON        string
	factory         *cmdutil.Factory
}

type object struct {
	BucketName string `json:"bucket_name,omitempty"`
	ObjectKey  string `json:"object_key,omitempty"`
	Source     string `json:"source,omitempty"`
	fileJSON   string
	factory    *cmdutil.Factory
}
