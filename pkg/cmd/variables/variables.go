package variables

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/variables"
	"github.com/aziontech/azion-cli/pkg/cmd/variables/list"
	"github.com/aziontech/azion-cli/pkg/cmd/variables/delete"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)


func NewCmd(f *cmdutil.Factory) *cobra.Command {
	variablesCmd := &cobra.Command{
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

	variablesCmd.AddCommand(list.NewCmd(f)
	variablesCmd.AddCommand(delete.NewCmd(f))
	variablesCmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return variablesCmd
}

