package cachesetting

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/cache_setting"

	api "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	ApplicationID           int64
	CacheSettingID          int64
	Name                    string
	browserCacheSettings    string
	browserCacheMaxAge      int64
	cacheByQueryString      string
	queryStringFields       []string
	enableQueryStringSort   string
	cacheByCookies          string
	cookieNames             []string
	enableCachingForPost    string
	enableCachingForOptions string
	Path                    string
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
			client := api.NewClientV4(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := api.RequestUpdate{}

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

			if request.Name == nil || *request.Name == "" {
				resp, err := client.Get(context.Background(), fields.ApplicationID, fields.CacheSettingID)
				if err != nil {
					return fmt.Errorf(msg.ErrorGetCache.Error(), err)
				}
				name := resp.GetName()
				request.Name = &name
			}

			response, err := client.Update(context.Background(), &request, fields.ApplicationID, fields.CacheSettingID)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateCacheSettings.Error(), err)
			}

			data := response.GetData()
			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.UpdateOutputSuccess, data.GetId()),
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
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.FlagApplicationID)
	flags.Int64Var(&fields.CacheSettingID, "cache-setting-id", 0, msg.FlagCacheSettingID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.browserCacheSettings, "browser-cache-settings", "honor", msg.FlagBrowserCacheBehavior)
	flags.StringSliceVar(&fields.queryStringFields, "query-string-fields", []string{}, msg.FlagQueryStringFields)
	flags.StringSliceVar(&fields.cookieNames, "cookie-names", []string{}, msg.FlagCookieNames)
	flags.StringVar(&fields.cacheByCookies, "cache-by-cookies", "ignore", msg.FlagCacheByCookiesEnabled)
	flags.StringVar(&fields.cacheByQueryString, "cache-by-query-string", "ignore", msg.FlagCacheByQueryString)
	flags.StringVar(&fields.enableCachingForOptions, "enable-caching-for-options", "false", msg.FlagCachingForOptionsEnabled)
	flags.StringVar(&fields.enableCachingForPost, "enable-caching-for-post", "", msg.FlagCachingForPostEnabled)
	flags.StringVar(&fields.enableQueryStringSort, "enable-caching-string-sort", "", msg.FlagCachingStringSortEnabled)
	flags.Int64Var(&fields.browserCacheMaxAge, "browser-cache-max-age", 0, msg.FlagBrowserCacheMaxAge)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.UpdateFlagHelp)
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.RequestUpdate) error {
	request.SetName(fields.Name)
	if cmd.Flags().Changed("browser-cache-settings") {
		if fields.browserCacheSettings == "override" && !cmd.Flags().Changed("browser-cache-max-age") {
			return msg.ErrorBrowserMaximumTtlNotSent
		}

		req := sdk.BrowserCacheModuleRequest{}
		req.SetBehavior(fields.browserCacheSettings)
		req.SetMaxAge(fields.browserCacheMaxAge)
		request.SetBrowserCache(req)
	}

	if !anyFlagChanged(cmd, acceleratorFlags...) {
		return nil
	}

	// The sub-objects are pointers in the SDK, so they must be allocated before
	// being written to, and the modules struct must be set back on the request.
	mods := request.GetModules()
	if mods.ApplicationAccelerator == nil {
		mods.ApplicationAccelerator = &sdk.CacheSettingsApplicationAcceleratorModuleRequest{}
	}
	appAcc := mods.ApplicationAccelerator

	// Allocate each sub-object only when a flag actually writes to it, so a
	// PATCH never carries an empty object for a module the user did not touch.
	if anyFlagChanged(cmd, "query-string-fields", "cache-by-query-string", "enable-caching-string-sort") {
		if appAcc.CacheVaryByQuerystring == nil {
			appAcc.CacheVaryByQuerystring = &sdk.CacheVaryByQuerystringModuleRequest{}
		}

		if cmd.Flags().Changed("query-string-fields") {
			appAcc.CacheVaryByQuerystring.SetFields(fields.queryStringFields)
		}

		if cmd.Flags().Changed("cache-by-query-string") {
			appAcc.CacheVaryByQuerystring.SetBehavior(fields.cacheByQueryString)
		}

		if cmd.Flags().Changed("enable-caching-string-sort") {
			stringSort, err := strconv.ParseBool(fields.enableQueryStringSort)
			if err != nil {
				return fmt.Errorf("%w: %q", msg.ErrorCachingStringSortFlag, fields.enableQueryStringSort)
			}
			appAcc.CacheVaryByQuerystring.SetSortEnabled(stringSort)
		}
	}

	if anyFlagChanged(cmd, "cookie-names", "cache-by-cookies") {
		if appAcc.CacheVaryByCookies == nil {
			appAcc.CacheVaryByCookies = &sdk.CacheVaryByCookiesModuleRequest{}
		}

		if cmd.Flags().Changed("cookie-names") {
			appAcc.CacheVaryByCookies.SetCookieNames(fields.cookieNames)
		}

		if cmd.Flags().Changed("cache-by-cookies") {
			appAcc.CacheVaryByCookies.SetBehavior(fields.cacheByCookies)
		}
	}

	if anyFlagChanged(cmd, "enable-caching-for-options", "enable-caching-for-post") {
		cacheByMethod := appAcc.GetCacheVaryByMethod()

		if cmd.Flags().Changed("enable-caching-for-options") {
			cachingOptions, err := strconv.ParseBool(fields.enableCachingForOptions)
			if err != nil {
				return fmt.Errorf("%w: %q", msg.ErrorCachingForOptionsFlag, fields.enableCachingForOptions)
			}
			cacheByMethod = setCacheVaryByMethod(cacheByMethod, "options", cachingOptions)
		}

		if cmd.Flags().Changed("enable-caching-for-post") {
			cachingPost, err := strconv.ParseBool(fields.enableCachingForPost)
			if err != nil {
				return fmt.Errorf("%w: %q", msg.ErrorCachingForPostFlag, fields.enableCachingForPost)
			}
			cacheByMethod = setCacheVaryByMethod(cacheByMethod, "post", cachingPost)
		}

		appAcc.SetCacheVaryByMethod(cacheByMethod)
	}

	request.SetModules(mods)
	return nil
}

// acceleratorFlags are the flags handled by the Application Accelerator module.
var acceleratorFlags = []string{
	"query-string-fields",
	"cache-by-query-string",
	"enable-caching-string-sort",
	"cookie-names",
	"cache-by-cookies",
	"enable-caching-for-options",
	"enable-caching-for-post",
}

func anyFlagChanged(cmd *cobra.Command, names ...string) bool {
	for _, name := range names {
		if cmd.Flags().Changed(name) {
			return true
		}
	}
	return false
}

// setCacheVaryByMethod adds method to the cache_vary_by_method list when enabled
// is true and removes it when false, keeping the other methods untouched. The
// returned slice is always non-nil so that clearing the list serializes as [].
func setCacheVaryByMethod(methods []string, method string, enabled bool) []string {
	out := make([]string, 0, len(methods)+1)
	for _, m := range methods {
		if m != method {
			out = append(out, m)
		}
	}
	if enabled {
		out = append(out, method)
	}
	return out
}
