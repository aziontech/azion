package variables

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/variables"
	"github.com/aziontech/azion-cli/pkg/cmd/variables/describe"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	deviceGroupsCmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azioncli variables --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	deviceGroupsCmd.AddCommand(describe.NewCmd(f))

	deviceGroupsCmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return deviceGroupsCmd
}
