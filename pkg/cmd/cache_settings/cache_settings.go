package cache_settings

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/cache_settings"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings/create"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings/list"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// cmd represents the cache settings command
	cacheCmd := &cobra.Command{
		Use:   msg.CacheSettingsUsage,
		Short: msg.CacheSettingsShortDescription,
		Long:  msg.CacheSettingsLongDescription,
		Example: heredoc.Doc(`
		$ azioncli cache_settings --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cacheCmd.AddCommand(create.NewCmd(f))
	cacheCmd.AddCommand(update.NewCmd(f))
	cacheCmd.AddCommand(list.NewCmd(f))
	cacheCmd.Flags().BoolP("help", "h", false, msg.CacheSettingsFlagHelp)
	return cacheCmd
}
