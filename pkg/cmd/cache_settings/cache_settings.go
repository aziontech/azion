package cache_settings 

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/cache_settings"
	"github.com/aziontech/azion-cli/pkg/cmd/origins/list"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	originsCmd := &cobra.Command{
    Use:   msg.OriginsUsage,
		Short: msg.CacheSettingsShortDescription,
		Long:  msg.CacheSettingsLongDescription, Example: heredoc.Doc(`
		$ azioncli origins --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	originsCmd.AddCommand(list.NewCmd(f))
	originsCmd.Flags().BoolP("help", "h", false, msg.OriginsFlagHelp)
	return originsCmd
}
