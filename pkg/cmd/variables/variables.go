package variables

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmd/variables/create"
	"github.com/aziontech/azion-cli/pkg/cmd/variables/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/variables/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/variables/list"
	"github.com/aziontech/azion-cli/pkg/cmd/variables/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	msg "github.com/aziontech/azion-cli/pkg/messages/variables"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	variablesCmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion variables --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	variablesCmd.AddCommand(describe.NewCmd(f))
	variablesCmd.AddCommand(list.NewCmd(f))
	variablesCmd.AddCommand(delete.NewCmd(f))
	variablesCmd.AddCommand(create.NewCmd(f))
	variablesCmd.AddCommand(update.NewCmd(f))
	variablesCmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return variablesCmd
}
