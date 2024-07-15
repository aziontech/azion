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
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	applicationID   int64
	cacheSettingsID int64
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, int64, int64) (api.GetResponse, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, appID, cacheID int64) (api.GetResponse, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			return client.Get(ctx, appID, cacheID)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
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
				answer, err := describe.AskInput(msg.DescibeAskInputApplicationID)
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
				answer, err := describe.AskInput(msg.DescribeAskInputCacheID)
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

			ctx := context.Background()
			resp, err := describe.Get(ctx, applicationID, cacheSettingsID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetCache.Error(), err)
			}

			fields := make(map[string]string, 0)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["BrowserCacheSettings"] = "Browser Cache Settings"
			fields["BrowserCacheSettingsMaximumTtl"] = "Browser Cache Settings maximum TTL"
			fields["CdnCacheSettings"] = "Cdn Cache Settings"
			fields["CdnCacheSettingsMaximumTtl"] = "Cdn Cache Settings maximum TTL"
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
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: resp,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&applicationID, "application-id", 0, msg.DescribeFlagApplicationID)
	cobraCmd.Flags().Int64Var(&cacheSettingsID, "cache-setting-id", 0, msg.DescribeFlagCacheSettingsID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DescribeHelpFlag)
	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
