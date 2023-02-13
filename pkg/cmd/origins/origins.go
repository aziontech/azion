package origins

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/origins"
	"github.com/aziontech/azion-cli/pkg/cmd/origins/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/origins/list"
	"github.com/aziontech/azion-cli/pkg/cmd/origins/create"
	"github.com/aziontech/azion-cli/pkg/cmd/origins/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	originsCmd := &cobra.Command{
    Use:   msg.OriginsUsage,
		Short: msg.OriginsShortDescription,
		Long:  msg.OriginsLongDescription, Example: heredoc.Doc(`
		$ azioncli origins --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	originsCmd.AddCommand(list.NewCmd(f))
	originsCmd.AddCommand(describe.NewCmd(f))
	originsCmd.AddCommand(create.NewCmd(f))
	originsCmd.AddCommand(update.NewCmd(f))
	originsCmd.Flags().BoolP("help", "h", false, msg.OriginsFlagHelp)

	return originsCmd
}
