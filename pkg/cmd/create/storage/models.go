package storage

import (
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type fieldsObjects struct {
	BucketName string `json:"bucket-name"`
	ObjectKey  string `json:"object-key"`
	Source     string `json:"source"`
	fileJSON   string
}

type FieldsBucket struct {
	Name            string
	WorkloadsAccess string
	FileJSON        string
	Factory         *cmdutil.Factory
}

type Fields interface {
	RunE(cmd *cobra.Command, args []string) error
	AddFlags(flags *pflag.FlagSet)
	CreateRequestFromFlags(cmd *cobra.Command, request *api.RequestBucket) error
}
