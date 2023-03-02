package describe

import (
    "context"
    "encoding/json"
    "fmt"

    "github.com/fatih/color"

    "github.com/MakeNowJust/heredoc"
    "github.com/MaxwelMazur/tablecli"
    msg "github.com/aziontech/azion-cli/messages/cache_settings"

    api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
    "github.com/aziontech/azion-cli/pkg/cmdutil"
    "github.com/aziontech/azion-cli/pkg/contracts"
    "github.com/aziontech/azion-cli/utils"
    "github.com/spf13/cobra"
)

var (
    applicationID   int64
    cacheSettingsID int64
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
    opts := &contracts.DescribeOptions{}
    cmd := &cobra.Command{
        Use:           msg.CacheSettingsDescribeUsage,
        Short:         msg.CacheSettingsDescribeShortDescription,
        Long:          msg.CacheSettingsDescribeLongDescription,
        SilenceUsage:  true,
        SilenceErrors: true,
        Example: heredoc.Doc(`
        $ azioncli cache_settings describe --application-id 1673635839 --cache-settings-id 107313
        $ azioncli cache_settings describe --application-id 1673635839 --cache-settings-id 107313 --format json
        $ azioncli cache_settings describe --application-id 1673635839 --cache-settings-id 107313 --out "./tmp/test.json" --format json
        `),
        RunE: func(cmd *cobra.Command, args []string) error {
            if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("cache-settings-id") {
                return msg.ErrorMissingArguments
            }

            client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
            ctx := context.Background()
            resp, err := client.GetCacheSettings(ctx, applicationID, cacheSettingsID)
            if err != nil {
                return fmt.Errorf(msg.ErrorGetCache.Error(), err)
            }

            out := f.IOStreams.Out
            formattedFuction, err := format(cmd, resp)
            if err != nil {
                return utils.ErrorFormatOut
            }

            if cmd.Flags().Changed("out") {
                err := cmdutil.WriteDetailsToFile(formattedFuction, opts.OutPath, out)
                if err != nil {
                    return fmt.Errorf("%s: %w", utils.ErrorWriteFile, err)
                }
                fmt.Fprintf(out, msg.CacheSettingsFileWritten, opts.OutPath)
            } else {
                _, err := out.Write(formattedFuction[:])
                if err != nil {
                    return err
                }
            }

            return nil
        },
    }

    cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, msg.CacheSettingsDescribeFlagApplicationID)
    cmd.Flags().Int64VarP(&cacheSettingsID, "cache-settings-id", "c", 0, msg.CacheSettingsDescribeFlagCacheSettingsID)
    cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.CacheSettingsDescribeFlagOut)
    cmd.Flags().StringVar(&opts.Format, "format", "", msg.CacheSettingsDescribeFlagFormat)
    cmd.Flags().BoolP("help", "h", false, msg.CacheSettingsDescribeHelpFlag)
    return cmd
}

func format(cmd *cobra.Command, strResp api.CacheSettingsResponse) ([]byte, error) {
    format, err := cmd.Flags().GetString("format")
    if err != nil {
        return nil, err
    }

    if format == "json" || cmd.Flags().Changed("out") {
        return json.MarshalIndent(strResp, "", " ")
    }

    tbl := tablecli.New("", "")
    tbl.WithFirstColumnFormatter(color.New(color.FgGreen).SprintfFunc())
    tbl.AddRow("Id: ", strResp.GetId())
    tbl.AddRow("Name: ", strResp.GetName())
    tbl.AddRow("Browser cache settings: ", strResp.GetBrowserCacheSettings())
    tbl.AddRow("Browser cache settings maximum TTL: ", strResp.GetBrowserCacheSettingsMaximumTtl())
    tbl.AddRow("Cdn cache settings: ", strResp.GetCdnCacheSettings())
    tbl.AddRow("Cdn cache settings maximum TTL: ", strResp.GetCdnCacheSettingsMaximumTtl())
    tbl.AddRow("Cache by query string: ", strResp.GetCacheByQueryString())
    tbl.AddRow("Query string fields: ", strResp.GetQueryStringFields())
    tbl.AddRow("Enable query string sort: ", strResp.GetEnableCachingForPost())
    tbl.AddRow("Cache by cookies: ", strResp.GetCacheByCookies())
    tbl.AddRow("Cookie Names: ", strResp.GetCookieNames())
    tbl.AddRow("Adaptive delivery action: ", strResp.GetAdaptiveDeliveryAction())
    tbl.AddRow("Device group: ", strResp.GetDeviceGroup())
    tbl.AddRow("EnableCachingForPost: ", strResp.GetEnableCachingForPost())
    tbl.AddRow("L2 caching enabled: ", strResp.GetL2CachingEnabled())
    return tbl.GetByteFormat(), nil
}
