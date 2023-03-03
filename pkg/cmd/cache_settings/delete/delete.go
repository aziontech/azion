package delete

import (
    "context"
    "fmt"

    "github.com/MakeNowJust/heredoc"
    msg "github.com/aziontech/azion-cli/messages/cache_settings"
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
        Use:           msg.CacheSettingsDeleteUsage,
        Short:         msg.CacheSettingsDeleteShortDescription,
        Long:          msg.CacheSettingsDeleteLongDescription,
        SilenceUsage:  true,
        SilenceErrors: true,
        Example: heredoc.Doc(`
        $ azioncli cache_settings delete --application-id 1673635839 --cache-settings-id 107313
        $ azioncli cache_settings delete -a 1673635839 -c 107313
        `),
        RunE: func(cmd *cobra.Command, args []string) error {
            if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("cache-settings-id") {
                return msg.ErrorMissingArguments
            }
            if err := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token")).
            DeleteCacheSettings(context.Background(), applicationID, cacheSettingsID); err != nil {
                return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
            }
            fmt.Fprintf(f.IOStreams.Out, msg.CacheSettingsDeleteOutputSuccess, cacheSettingsID)
            return nil
        },
    }


    cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, msg.CacheSettingsDeleteFlagApplicationID)
    cmd.Flags().Int64VarP(&cacheSettingsID, "cache-settings-id", "c", 0, msg.CacheSettingsDeleteFlagCacheSettingsID)
    cmd.Flags().BoolP("help", "h", false, msg.CacheSettingsDeleteHelpFlag)
    return cmd
}
