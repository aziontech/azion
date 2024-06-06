package cachesetting

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/cache_setting"

	api "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	apiEdgeApp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	ApplicationID                  int64
	CacheSettingID                 int64
	Name                           string
	browserCacheSettings           string
	adaptiveDeliveryAction         string
	browserCacheSettingsMaximumTtl int64
	cdnCacheSettings               string
	cdnCacheSettingsMaximumTtl     int64
	cacheByQueryString             string
	queryStringFields              []string
	enableQueryStringSort          string
	cacheByCookies                 string
	cookieNames                    []string
	enableCachingForPost           string
	enableCachingForOptions        string
	l2CachingEnabled               string
	isSliceConfigurationEnabled    string
	isSliceL2CachingEnabled        string
	sliceConfigurationRange        int64
	Path                           string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.UpdateShortDescription,
		Long:          msg.UpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion update cache-setting --application-id 1673635839 --cache-setting-id 123123421 --name "phototypesetting"
        $ azion update cache-setting --application-id 1673635839 --cache-setting-id 123123421 --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			clientEdgeApp := apiEdgeApp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			request := api.UpdateRequest{}

			if !cmd.Flags().Changed("application-id") {
				answers, err := utils.AskInput(msg.CreateAskInputApplicationID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				applicationID, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return utils.ErrorConvertingStringToInt
				}

				fields.ApplicationID = int64(applicationID)
			}

			if !cmd.Flags().Changed("cache-setting-id") {
				answers, err := utils.AskInput(msg.UpdateAskInputCacheSettingID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				cacheSettingID, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return utils.ErrorConvertingStringToInt
				}

				fields.CacheSettingID = int64(cacheSettingID)
			}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			if err := appAccelerationNoEnabled(clientEdgeApp, fields, request); err != nil {
				return err
			}

			response, err := client.Update(context.Background(), &request, fields.ApplicationID, fields.CacheSettingID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateCacheSettings.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.UpdateOutputSuccess, response.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&updateOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)
	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	flags.Int64Var(&fields.CacheSettingID, "cache-setting-id", 0, msg.FlagCacheSettingID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.browserCacheSettings, "browser-cache-settings", "honor", msg.FlagBrowserCacheSettings)
	flags.StringSliceVar(&fields.queryStringFields, "query-string-fields", []string{}, msg.FlagQueryStringFields)
	flags.StringSliceVar(&fields.cookieNames, "cookie-names", []string{}, msg.FlagCookieNames)
	flags.StringVar(&fields.cacheByCookies, "cache-by-cookies", "ignore", msg.FlagCacheByCookiesEnabled)
	flags.StringVar(&fields.cacheByQueryString, "cache-by-query-string", "ignore", msg.FlagCacheByQueryString)
	flags.StringVar(&fields.cdnCacheSettings, "cdn-cache-settings", "honor", msg.FlagCdnCacheSettingsEnabled)
	flags.StringVar(&fields.enableCachingForOptions, "enable-caching-for-options", "false", msg.FlagCachingForOptionsEnabled)
	flags.StringVar(&fields.enableCachingForPost, "enable-caching-for-post", "", msg.FlagCachingForPostEnabled)
	flags.StringVar(&fields.enableQueryStringSort, "enable-caching-string-sort", "", msg.FlagCachingStringSortEnabled)
	flags.StringVar(&fields.isSliceConfigurationEnabled, "slice-configuration-enabled", "", msg.FlagSliceConfigurationEnabled)
	flags.StringVar(&fields.isSliceL2CachingEnabled, "slice-l2-caching-enabled", "", msg.FlagSliceL2CachingEnabled)
	flags.StringVar(&fields.l2CachingEnabled, "l2-caching-enabled", "", msg.FlagL2CachingEnabled)
	flags.Int64Var(&fields.sliceConfigurationRange, "slice-configuration-range", 0, msg.FlagSliceConfigurationRange)
	flags.Int64Var(&fields.cdnCacheSettingsMaximumTtl, "cnd-cache-settings-maximum-ttl", 60, msg.FlagCdnCacheSettingsMaxTtl)
	flags.Int64Var(&fields.browserCacheSettingsMaximumTtl, "browser-cache-settings-maximum-ttl", 0, msg.FlagBrowserCacheSettingsMaxTtl)
	flags.StringVar(&fields.adaptiveDeliveryAction, "adaptive-delivery-action", "ignore", msg.FlagAdaptiveDeliveryAction)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.UpdateFlagHelp)
}

func appAccelerationNoEnabled(client *apiEdgeApp.Client, fields *Fields, request api.UpdateRequest) error {
	ctx := context.Background()
	str := strconv.FormatInt(fields.ApplicationID, 10)
	application, err := client.Get(ctx, str)
	if err != nil {
		return err
	}

	if (request.GetEnableCachingForOptions() ||
		request.GetEnableCachingForPost() ||
		request.GetEnableQueryStringSort()) &&
		!application.GetApplicationAcceleration() {
		return msg.ErrorApplicationAccelerationNotEnabled
	}
	return nil
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.UpdateRequest) error {
	request.SetName(fields.Name)
	if cmd.Flags().Changed("browser-cache-settings") {
		if fields.browserCacheSettings == "override" && !cmd.Flags().Changed("browser-cache-settings-maximum-ttl") {
			return msg.ErrorBrowserMaximumTtlNotSent
		}
		request.SetBrowserCacheSettings(fields.browserCacheSettings)
	}

	if cmd.Flags().Changed("query-string-fields") {
		request.SetQueryStringFields(fields.queryStringFields)
	}

	if cmd.Flags().Changed("cookie-names") {
		request.SetCookieNames(fields.cookieNames)
	}

	if cmd.Flags().Changed("cache-by-cookies") {
		request.SetCacheByCookies(fields.cacheByCookies)
	}

	if cmd.Flags().Changed("cache-by-query-string") {
		request.SetCacheByQueryString(fields.cacheByQueryString)
	}

	if cmd.Flags().Changed("cdn-cache-settings") {
		request.SetCdnCacheSettings(fields.cdnCacheSettings)
	}

	if cmd.Flags().Changed("slice-configuration-range") {
		request.SetSliceConfigurationRange(fields.sliceConfigurationRange)
	}

	if cmd.Flags().Changed("cnd-cache-settings-maximum-ttl") {
		request.SetCdnCacheSettingsMaximumTtl(fields.cdnCacheSettingsMaximumTtl)
	}

	if cmd.Flags().Changed("browser-cache-settings-maximum-ttl") {
		request.SetBrowserCacheSettingsMaximumTtl(fields.browserCacheSettingsMaximumTtl)
	}

	if cmd.Flags().Changed("adaptive-delivery-action") {
		request.SetAdaptiveDeliveryAction(fields.adaptiveDeliveryAction)
	}

	if cmd.Flags().Changed("browser-cache-settings-maximum-ttl") {
		request.SetBrowserCacheSettingsMaximumTtl(fields.browserCacheSettingsMaximumTtl)
	}

	if cmd.Flags().Changed("enable-caching-for-options") {
		cachingOptions, err := strconv.ParseBool(fields.enableCachingForOptions)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorCachingForOptionsFlag, fields.enableCachingForOptions)
		}
		request.SetEnableCachingForPost(cachingOptions)
	}

	if cmd.Flags().Changed("enable-caching-for-post") {
		cachingPost, err := strconv.ParseBool(fields.enableCachingForPost)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorCachingForPostFlag, fields.enableCachingForPost)
		}
		request.SetEnableCachingForPost(cachingPost)
	}

	if cmd.Flags().Changed("enable-caching-string-sort") {
		stringSort, err := strconv.ParseBool(fields.enableQueryStringSort)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorCachingStringSortFlag, fields.enableQueryStringSort)
		}
		request.SetEnableQueryStringSort(stringSort)
	}

	if cmd.Flags().Changed("slice-configuration-enabled") {
		sliceEnable, err := strconv.ParseBool(fields.isSliceConfigurationEnabled)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorSliceConfigurationFlag, fields.isSliceConfigurationEnabled)
		}
		request.SetIsSliceConfigurationEnabled(sliceEnable)
		request.SetIsSliceEdgeCachingEnabled(true) //Edge Cache is mandatory in this case
	}

	if cmd.Flags().Changed("slice-l2-caching-enabled") {
		sliceEnable, err := strconv.ParseBool(fields.isSliceL2CachingEnabled)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorSliceL2CachingFlag, fields.isSliceL2CachingEnabled)
		}
		request.SetIsSliceL2CachingEnabled(sliceEnable)
	}

	if cmd.Flags().Changed("l2-caching-enabled") {
		lsEnable, err := strconv.ParseBool(fields.l2CachingEnabled)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorL2CachingEnabledFlag, fields.l2CachingEnabled)
		}
		request.SetIsSliceConfigurationEnabled(lsEnable)
	}

	return nil
}
