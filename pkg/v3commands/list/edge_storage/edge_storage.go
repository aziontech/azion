package edge_storage

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	msg "github.com/aziontech/azion-cli/messages/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/v3commands/list/edge_storage/bucket"
	"github.com/aziontech/azion-cli/pkg/v3commands/list/edge_storage/object"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "edge-storage",
		Short:         msg.SHORT_DESCRIPTION_LIST,
		Long:          msg.LONG_DESCRIPTION_LIST,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list edge-storage bucket
			$ azion list edge-storage object
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(bucket.NewBucket(f))
	cmd.AddCommand(object.NewObject(f))
	cmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP)
	return cmd
}
