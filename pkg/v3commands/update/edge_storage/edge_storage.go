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
		Short:         msg.SHORT_DESCRIPTION_UPDATE,
		Long:          msg.LONG_DESCRIPTION_UPDATE,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion update edge-storage bucket --bucket-id 1234 --name 'Hello'
			$ azion update edge-storage object --object-id 1234 --name 'Hello'
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(NewBucket(f))
	cmd.AddCommand(NewObject(f))
	cmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP)
	return cmd
}
