package delete

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/cache_settings"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

var (
	applicationID   int64
	cacheSettingsID int64
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           cache_settings.CacheSettingsDeleteUsage,
		Short:         cache_settings.CacheSettingsDeleteShortDescription,
		Long:          cache_settings.CacheSettingsDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion cache_settings delete --application-id 1673635839 --cache-settings-id 107313
        $ azion cache_settings delete -a 1673635839 -c 107313
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("cache-settings-id") {
				return cache_settings.ErrorMissingArguments
			}
			if err := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token")).
				DeleteCacheSettings(context.Background(), applicationID, cacheSettingsID); err != nil {
				return fmt.Errorf(cache_settings.ErrorFailToDelete.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, cache_settings.CacheSettingsDeleteOutputSuccess, cacheSettingsID)
			return nil
		},
	}

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, cache_settings.CacheSettingsDeleteFlagApplicationID)
	cmd.Flags().Int64VarP(&cacheSettingsID, "cache-settings-id", "c", 0, cache_settings.CacheSettingsDeleteFlagCacheSettingsID)
	cmd.Flags().BoolP("help", "h", false, cache_settings.CacheSettingsDeleteHelpFlag)
	return cmd
}
