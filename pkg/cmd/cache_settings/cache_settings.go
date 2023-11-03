package cache_settings

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings/create"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings/delete"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings/describe"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings/list"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	msg "github.com/aziontech/azion-cli/pkg/messages/cache_settings"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	// cmd represents the cache settings command
	cacheCmd := &cobra.Command{
		Use:   msg.CacheSettingsUsage,
		Short: msg.CacheSettingsShortDescription,
		Long:  msg.CacheSettingsLongDescription,
		Example: heredoc.Doc(`
		$ azion cache_settings --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cacheCmd.AddCommand(create.NewCmd(f))
	cacheCmd.AddCommand(update.NewCmd(f))
	cacheCmd.AddCommand(list.NewCmd(f))
	cacheCmd.AddCommand(describe.NewCmd(f))
	cacheCmd.AddCommand(delete.NewCmd(f))
	cacheCmd.Flags().BoolP("help", "h", false, msg.CacheSettingsFlagHelp)
	return cacheCmd
}
