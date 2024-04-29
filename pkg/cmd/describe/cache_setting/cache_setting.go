package cachesetting

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/cache_setting"

	api "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
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
		Use:           msg.Usage,
		Short:         msg.DescribeShortDescription,
		Long:          msg.DescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion describe cache-setting --application-id 1673635839 --cache-setting-id 107313
        $ azion describe cache-setting --application-id 1673635839 --cache-setting-id 107313 --format json
        $ azion describe cache-setting --application-id 1673635839 --cache-setting-id 107313 --out "./tmp/test.json" 
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("application-id") {
				answer, err := utils.AskInput(msg.DescibeAskInputApplicationID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				applicationID = num
			}

			if !cmd.Flags().Changed("cache-setting-id") {
				answer, err := utils.AskInput(msg.DescribeAskInputCacheID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdApplication
				}

				cacheSettingsID = num
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			ctx := context.Background()
			resp, err := client.Get(ctx, applicationID, cacheSettingsID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetCache.Error(), err)
			}

			fields := make(map[string]string, 0)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["BrowserCacheSettings"] = "Browser cache settings"
			fields["BrowserCacheSettingsMaximumTtl"] = "Browser cache settings maximum TTL"
			fields["CdnCacheSettings"] = "Cdn cache settings"
			fields["CdnCacheSettingsMaximumTtl"] = "Cdn cache settings maximum TTL"
			fields["CacheByQueryString"] = "Cache by query string"
			fields["QueryStringFields"] = "Query string fiedlds"
			fields["EnableQueryStringSort"] = "Enable query string sort"
			fields["CacheByCookies"] = "Cache by cookies"
			fields["CookieNames"] = "Cookie Names"
			fields["AdaptiveDeliveryAction"] = "Adaptive delivery action"
			fields["DeviceGroup"] = "Device group"
			fields["EnableCachingForPost"] = "EnableCachingForPost"
			fields["L2CachingEnabled"] = "L2 caching enabled"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:         f.IOStreams.Out,
					Msg:         filepath.Clean(opts.OutPath),
					FlagOutPath: opts.OutPath,
					FlagFormat:  opts.Format,
				},
				Fields: fields,
				Values: resp,
			}
			return output.Print(&describeOut)
		},
	}

	cmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.DescribeFlagApplicationID)
	cmd.Flags().Int64Var(&cacheSettingsID, "cache-setting-id", 0, msg.DescribeFlagCacheSettingsID)
	cmd.Flags().StringVar(&opts.OutPath, "out", "", msg.DescribeFlagOut)
	cmd.Flags().StringVar(&opts.Format, "format", "", msg.DescribeFlagFormat)
	cmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)
	return cmd
}
