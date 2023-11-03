package update

import (
	"context"
	"fmt"
	apiEdgeApp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/cache_setting"

	api "github.com/aziontech/azion-cli/pkg/api/cache_setting"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID                  int64
	CacheSettingsID                int64
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
		Use:           msg.CacheSettingsUpdateUsage,
		Short:         msg.CacheSettingsUpdateShortDescription,
		Long:          msg.CacheSettingsUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion cache_settings update -a 1673635839 -c 115247 --name "cachesettingstest"
        $ azion cache_settings update -a 1673635839 -c 115247 --name "cachesettingswithfields" --browser-cache-settings honor --cdn-cache-settings honor --cache-by-query-string ignore 
        $ azion cache_settings update -a 1673635839 --in "update.json"
        $ azion cache_settings update -a 1674767911 -c 115247 --name "updateagain" --browser-cache-settings override --browser-cache-settings-maximum-ttl 60  --cdn-cache-settings honor --cnd-cache-settings-maximum-ttl 60 --cache-by-query-string ignore --cache-by-query-string whitelist --query-string-fields "heyyy,yoooo" --adaptive-delivery-action ignore --cache-by-cookies blacklist --cookie-names "nem,vem" --enable-caching-for-options true --enable-caching-for-post true --enable-caching-string-sort true --slice-configuration-enabled true
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			clientEdgeApp := apiEdgeApp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			request := api.UpdateRequest{}

			if cmd.Flags().Changed("in") {
				if !cmd.Flags().Changed("application-id") { // flags requireds
					return msg.ErrorMandatoryUpdateInFlags
				}
				var (
					file *os.File
					err  error
				)
				if fields.Path == "-" {
					file = os.Stdin
				} else {
					file, err = os.Open(fields.Path)
					if err != nil {
						return fmt.Errorf("%w: %s", utils.ErrorOpeningFile, fields.Path)
					}
				}
				err = cmdutil.UnmarshallJsonFromReader(file, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}

				ctx := context.Background()
				str := strconv.FormatInt(fields.ApplicationID, 10)
				application, err := clientEdgeApp.Get(ctx, str)
				if err != nil {
					return err
				}

				if (request.GetEnableCachingForOptions() || request.GetEnableCachingForPost() || request.GetEnableQueryStringSort()) && !application.GetApplicationAcceleration() {
					return msg.ErrorApplicationAccelerationNotEnabled
				}

			} else {
				if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("cache-settings-id") { // flags requireds
					return msg.ErrorMandatoryUpdateFlags
				}

				request.Id = fields.CacheSettingsID

				ctx := context.Background()
				str := strconv.FormatInt(fields.ApplicationID, 10)
				application, err := clientEdgeApp.Get(ctx, str)
				if err != nil {
					return err
				}

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

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
					if cachingOptions && !application.GetApplicationAcceleration() {
						return msg.ErrorApplicationAccelerationNotEnabled
					}
					request.SetEnableCachingForPost(cachingOptions)
				}

				if cmd.Flags().Changed("enable-caching-for-post") {
					cachingPost, err := strconv.ParseBool(fields.enableCachingForPost)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorCachingForPostFlag, fields.enableCachingForPost)
					}
					if cachingPost && !application.GetApplicationAcceleration() {
						return msg.ErrorApplicationAccelerationNotEnabled
					}
					request.SetEnableCachingForPost(cachingPost)
				}

				if cmd.Flags().Changed("enable-caching-string-sort") {
					stringSort, err := strconv.ParseBool(fields.enableQueryStringSort)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorCachingStringSortFlag, fields.enableQueryStringSort)
					}
					if stringSort && !application.GetApplicationAcceleration() {
						return msg.ErrorApplicationAccelerationNotEnabled
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
			}

			response, err := client.Update(context.Background(), &request, fields.ApplicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateCacheSettings.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.CacheSettingsUpdateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.CreateFlagEdgeApplicationId)
	flags.Int64VarP(&fields.CacheSettingsID, "cache-settings-id", "c", 0, msg.CacheSettingsId)
	flags.StringVar(&fields.Name, "name", "", msg.CreateFlagName)
	flags.StringVar(&fields.browserCacheSettings, "browser-cache-settings", "honor", msg.CreateFlagBrowserCacheSettings)
	flags.StringSliceVar(&fields.queryStringFields, "query-string-fields", []string{}, msg.CreateFlagQueryStringFields)
	flags.StringSliceVar(&fields.cookieNames, "cookie-names", []string{}, msg.CreateFlagCookieNames)
	flags.StringVar(&fields.cacheByCookies, "cache-by-cookies", "ignore", msg.CreateFlagCacheByCookies)
	flags.StringVar(&fields.cacheByQueryString, "cache-by-query-string", "ignore", msg.CreateFlagCacheByQueryString)
	flags.StringVar(&fields.cdnCacheSettings, "cdn-cache-settings", "honor", msg.CreateFlagCdnCacheSettingsEnabled)
	flags.StringVar(&fields.enableCachingForOptions, "enable-caching-for-options", "false", msg.CreateFlagCachingForOptionsEnabled)
	flags.StringVar(&fields.enableCachingForPost, "enable-caching-for-post", "", msg.CreateFlagCachingForPostEnabled)
	flags.StringVar(&fields.enableQueryStringSort, "enable-caching-string-sort", "", msg.CreateFlagCachingStringSortEnabled)
	flags.StringVar(&fields.isSliceConfigurationEnabled, "slice-configuration-enabled", "", msg.CreateFlagSliceConfigurationEnabled)
	flags.StringVar(&fields.isSliceL2CachingEnabled, "slice-l2-caching-enabled", "", msg.CreateFlagSliceL2CachingEnabled)
	flags.StringVar(&fields.l2CachingEnabled, "l2-caching-enabled", "", msg.CreateFlagL2CachingEnabled)
	flags.Int64Var(&fields.sliceConfigurationRange, "slice-configuration-range", 0, msg.CreateFlagSliceConfigurationRange)
	flags.Int64Var(&fields.cdnCacheSettingsMaximumTtl, "cnd-cache-settings-maximum-ttl", 60, msg.CreateFlagCdnCacheSettingsMaxTtl)
	flags.Int64Var(&fields.browserCacheSettingsMaximumTtl, "browser-cache-settings-maximum-ttl", 0, msg.CreateFlagBrowserCacheSettingsMaxTtl)
	flags.StringVar(&fields.adaptiveDeliveryAction, "adaptive-delivery-action", "ignore", msg.CreateFlagAdaptiveDeliveryAction)
	flags.StringVar(&fields.Path, "in", "", msg.CreateFlagIn)
	flags.BoolP("help", "h", false, msg.CreateHelpFlag)
	return cmd
}
