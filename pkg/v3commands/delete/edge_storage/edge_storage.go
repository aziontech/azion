package edge_storage

import (
	"github.com/spf13/cobra"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	bucketpkg "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_storage/bucket"
	objectpkg "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_storage/object"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.USAGE,
		Short:         msg.SHORT_DESCRIPTION_DELETE,
		Long:          msg.LONG_DESCRIPTION_DELETE,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       msg.EXAMPLE_DELETE,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(bucketpkg.NewBucket(f))
	cmd.AddCommand(objectpkg.NewObject(f))
	cmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP)
	return cmd
}
