package edge_storage

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	msg "github.com/aziontech/azion-cli/messages/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "edge-storage",
		Short:         msg.SHORT_DESCRIPTION_DESCRIBE,
		Long:          msg.LONG_DESCRIPTION_DESCRIBE,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion describe edge-storage bucket --bucket-id 1234
		$ azion describe edge-storage object --bucket-id 1234 --object-key 'object-key'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewObject(f))
	cmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP)
	return cmd
}
