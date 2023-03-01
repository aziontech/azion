package update

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/MakeNowJust/heredoc"

	msg "github.com/aziontech/azion-cli/messages/cache_settings"
	msgapp "github.com/aziontech/azion-cli/messages/edge_applications"

	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	ApplicationID                      int64
	CacheSettingsID                    int64
	Name                               string
	browser_cache_settings             string
	adaptive_delivery_action           string
	browser_cache_settings_maximum_ttl int64
	cdn_cache_settings                 string
	cdn_cache_settings_maximum_ttl     int64
	cache_by_query_string              string
	query_string_fields                []string
	enable_query_string_sort           string
	cache_by_cookies                   string
	cookie_names                       []string
	enable_caching_for_post            string
	enable_caching_for_options         string
	l2_caching_enabled                 string
	is_slice_configuration_enabled     string
	is_slice_l2_caching_enabled        string
	slice_configuration_range          int64
	Path                               string
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
        $ azioncli cache_settings update -a 1673635839 -c 115247 --name "cachesettingstest"
        $ azioncli cache_settings update -a 1673635839 -c 115247 --name "cachesettingswithfields" --browser-cache-settings honor --cdn-cache-settings honor --cache-by-query-string ignore 
        $ azioncli cache_settings update -a 1673635839 --in "update.json"
		$ azioncli cache_settings update -a 1674767911 -c 115247 --name "updateagain" --browser-cache-settings override --browser-cache-settings-maximum-ttl 60  --cdn-cache-settings honor --cnd-cache-settings-maximum-ttl 60 --cache-by-query-string ignore --cache-by-query-string whitelist --query-string-fields "heyyy,yoooo" --adaptive-delivery-action ignore --cache-by-cookies blacklist --cookie-names "nem,vem" --enable-caching-for-options true --enable-caching-for-post true --enable-caching-string-sort true --slice-configuration-enabled true
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			request := api.UpdateCacheSettingsRequest{}

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
				application, err := client.Get(ctx, str)
				if err != nil {
					return fmt.Errorf(msgapp.ErrorGetApplication.Error(), err)
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
				application, err := client.Get(ctx, str)
				if err != nil {
					return fmt.Errorf(msgapp.ErrorGetApplication.Error(), err)
				}

				if cmd.Flags().Changed("name") {
					request.SetName(fields.Name)
				}

				if cmd.Flags().Changed("browser-cache-settings") {
					if fields.browser_cache_settings == "override" && !cmd.Flags().Changed("browser-cache-settings-maximum-ttl") {
						return msg.ErrorBrowserMaximumTtlNotSent
					}
					request.SetBrowserCacheSettings(fields.browser_cache_settings)
				}

				if cmd.Flags().Changed("query-string-fields") {
					request.SetQueryStringFields(fields.query_string_fields)
				}

				if cmd.Flags().Changed("cookie-names") {
					request.SetCookieNames(fields.cookie_names)
				}

				if cmd.Flags().Changed("cache-by-cookies") {
					request.SetCacheByCookies(fields.cache_by_cookies)
				}

				if cmd.Flags().Changed("cache-by-query-string") {
					request.SetCacheByQueryString(fields.cache_by_query_string)
				}

				if cmd.Flags().Changed("cdn-cache-settings") {
					request.SetCdnCacheSettings(fields.cdn_cache_settings)
				}

				if cmd.Flags().Changed("slice-configuration-range") {
					request.SetSliceConfigurationRange(fields.slice_configuration_range)
				}

				if cmd.Flags().Changed("cnd-cache-settings-maximum-ttl") {
					request.SetCdnCacheSettingsMaximumTtl(fields.cdn_cache_settings_maximum_ttl)
				}

				if cmd.Flags().Changed("browser-cache-settings-maximum-ttl") {
					request.SetBrowserCacheSettingsMaximumTtl(fields.browser_cache_settings_maximum_ttl)
				}

				if cmd.Flags().Changed("adaptive-delivery-action") {
					request.SetAdaptiveDeliveryAction(fields.adaptive_delivery_action)
				}

				if cmd.Flags().Changed("browser-cache-settings-maximum-ttl") {
					request.SetBrowserCacheSettingsMaximumTtl(fields.browser_cache_settings_maximum_ttl)
				}

				if cmd.Flags().Changed("enable-caching-for-options") {
					cachingOptions, err := strconv.ParseBool(fields.enable_caching_for_options)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorCachingForOptionsFlag, fields.enable_caching_for_options)
					}
					if cachingOptions && !application.GetApplicationAcceleration() {
						return msg.ErrorApplicationAccelerationNotEnabled
					}
					request.SetEnableCachingForPost(cachingOptions)
				}

				if cmd.Flags().Changed("enable-caching-for-post") {
					cachingPost, err := strconv.ParseBool(fields.enable_caching_for_post)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorCachingForPostFlag, fields.enable_caching_for_post)
					}
					if cachingPost && !application.GetApplicationAcceleration() {
						return msg.ErrorApplicationAccelerationNotEnabled
					}
					request.SetEnableCachingForPost(cachingPost)
				}

				if cmd.Flags().Changed("enable-caching-string-sort") {
					stringSort, err := strconv.ParseBool(fields.enable_query_string_sort)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorCachingStringSortFlag, fields.enable_query_string_sort)
					}
					if stringSort && !application.GetApplicationAcceleration() {
						return msg.ErrorApplicationAccelerationNotEnabled
					}
					request.SetEnableQueryStringSort(stringSort)
				}

				if cmd.Flags().Changed("slice-configuration-enabled") {
					sliceEnable, err := strconv.ParseBool(fields.is_slice_configuration_enabled)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorSliceConfigurationFlag, fields.is_slice_configuration_enabled)
					}
					request.SetIsSliceConfigurationEnabled(sliceEnable)
					request.SetIsSliceEdgeCachingEnabled(true) //Edge Cache is mandatory in this case
				}

				if cmd.Flags().Changed("slice-l2-caching-enabled") {
					sliceEnable, err := strconv.ParseBool(fields.is_slice_l2_caching_enabled)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorSliceL2CachingFlag, fields.is_slice_l2_caching_enabled)
					}
					request.SetIsSliceL2CachingEnabled(sliceEnable)
				}

				if cmd.Flags().Changed("l2-caching-enabled") {
					lsEnable, err := strconv.ParseBool(fields.l2_caching_enabled)
					if err != nil {
						return fmt.Errorf("%w: %q", msg.ErrorL2CachingEnabledFlag, fields.l2_caching_enabled)
					}
					request.SetIsSliceConfigurationEnabled(lsEnable)
				}
			}

			response, err := client.UpdateCacheSettings(context.Background(), &request, fields.ApplicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateCacheSettings.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, msg.CacheSettingsUpdateOutputSuccess, response.GetId())
			return nil
		},
	}

	flags := cmd.Flags()
	flags.Int64VarP(&fields.ApplicationID, "application-id", "a", 0, msg.CacheSettingsCreateFlagEdgeApplicationId)
	flags.Int64VarP(&fields.CacheSettingsID, "cache-settings-id", "c", 0, msg.CacheSettingsId)
	flags.StringVar(&fields.Name, "name", "", msg.CacheSettingsCreateFlagName)
	flags.StringVar(&fields.browser_cache_settings, "browser-cache-settings", "honor", msg.CacheSettingsCreateFlagBrowserCacheSettings)
	flags.StringSliceVar(&fields.query_string_fields, "query-string-fields", []string{}, msg.CacheSettingsCreateFlagQueryStringFields)
	flags.StringSliceVar(&fields.cookie_names, "cookie-names", []string{}, msg.CacheSettingsCreateFlagCookieNames)
	flags.StringVar(&fields.cache_by_cookies, "cache-by-cookies", "ignore", msg.CacheSettingsCreateFlagCacheByCookies)
	flags.StringVar(&fields.cache_by_query_string, "cache-by-query-string", "ignore", msg.CacheSettingsCreateFlagCacheByQueryString)
	flags.StringVar(&fields.cdn_cache_settings, "cdn-cache-settings", "honor", msg.CacheSettingsCreateFlagCdnCacheSettings)
	flags.StringVar(&fields.enable_caching_for_options, "enable-caching-for-options", "false", msg.CacheSettingsCreateFlagCachingForOptions)
	flags.StringVar(&fields.enable_caching_for_post, "enable-caching-for-post", "", msg.CacheSettingsCreateFlagCachingForPost)
	flags.StringVar(&fields.enable_query_string_sort, "enable-caching-string-sort", "", msg.CacheSettingsCreateFlagCachingStringSort)
	flags.StringVar(&fields.is_slice_configuration_enabled, "slice-configuration-enabled", "", msg.CacheSettingsCreateFlagSliceConfigurationEnabled)
	flags.StringVar(&fields.is_slice_l2_caching_enabled, "slice-l2-caching-enabled", "", msg.CacheSettingsCreateFlagSliceL2CachingEnabled)
	flags.StringVar(&fields.l2_caching_enabled, "l2-caching-enabled", "", msg.CacheSettingsCreateFlagL2CachingEnabled)
	flags.Int64Var(&fields.slice_configuration_range, "slice-configuration-range", 0, msg.CacheSettingsCreateFlagSliceConfigurationRange)
	flags.Int64Var(&fields.cdn_cache_settings_maximum_ttl, "cnd-cache-settings-maximum-ttl", 60, msg.CacheSettingsCreateFlagCdnCacheSettingsMaxTtl)
	flags.Int64Var(&fields.browser_cache_settings_maximum_ttl, "browser-cache-settings-maximum-ttl", 0, msg.CacheSettingsCreateFlagBrowserCacheSettingsMaxTtl)
	flags.StringVar(&fields.adaptive_delivery_action, "adaptive-delivery-action", "ignore", msg.CacheSettingsCreateFlagAdaptiveDeliveryAction)
	flags.StringVar(&fields.Path, "in", "", msg.CacheSettingsCreateFlagIn)
	flags.BoolP("help", "h", false, msg.CacheSettingsCreateHelpFlag)
	return cmd
}
