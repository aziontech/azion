package edge_storage

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	msg "github.com/aziontech/azion-cli/messages/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	bucketpkg "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_storage/bucket"
	objectpkg "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_storage/object"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "edge-storage",
		Short:         msg.SHORT_DESCRIPTION_DELETE,
		Long:          msg.LONG_DESCRIPTION_DELETE,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete edge-storage bucket --bucket-id 1234
		$ azion delete edge-storage object --bucket-id 1234 --object-key 'object-key'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(bucketpkg.NewBucket(f))
	cmd.AddCommand(objectpkg.NewObject(f))
	cmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP)
	return cmd
}
