package cachesetting

import (
	"context"
	"fmt"
	sdk "github.com/aziontech/azionapi-v4-go-sdk/edge"
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
	ApplicationID               int64
	Name                        string
	browserCacheBehavior        string
	browserCacheMaxAge          int64
	adaptiveDeliveryAction      string
	cacheByQueryString          string
	queryStringFields           []string
	enableQueryStringSort       string
	cacheByCookies              string
	cookieNames                 []string
	enableCachingForPost        string
	enableCachingForOptions     string
	isSliceConfigurationEnabled string
	sliceConfigurationRange     int64
	Path                        string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create cache-setting --application-id 1673635839 --name "phototypesetting"
        $ azion create cache-setting --application-id 1673635839 --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClientV4(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			clientEdgeApp := apiEdgeApp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			request := api.Request{}

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

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}

			} else {
				if !cmd.Flags().Changed("name") {
					answers, err := utils.AskInput(msg.CreateAskInputName)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					fields.Name = answers
				}

				err := createRequestFromFlags(cmd, fields, &request)
				if err != nil {
					return err
				}
			}

			if err := appAccelerationNoEnabled(clientEdgeApp, fields, request); err != nil {
				return err
			}

			response, err := client.Create(context.Background(), &request, fields.ApplicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateCacheSettings.Error(), err)
			}

			data := response.GetData()
			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.CreateOutputSuccess, data.GetId()),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&creatOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)
	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.FlagEdgeApplicationID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.browserCacheBehavior, "browser-cache-behavior", "honor", msg.FlagBrowserCacheBehavior)
	flags.Int64Var(&fields.browserCacheMaxAge, "browser-cache-max-age", 0, msg.FlagBrowserCacheMaxAge)
	flags.StringSliceVar(&fields.queryStringFields, "query-string-fields", []string{}, msg.FlagQueryStringFields)
	flags.StringSliceVar(&fields.cookieNames, "cookie-names", []string{}, msg.FlagCookieNames)
	flags.StringVar(&fields.cacheByCookies, "cache-by-cookies", "ignore", msg.FlagCacheByCookiesEnabled)
	flags.StringVar(&fields.cacheByQueryString, "cache-by-query-string", "ignore", msg.FlagCacheByQueryString)
	flags.StringVar(&fields.enableCachingForOptions, "enable-caching-for-options", "false", msg.FlagCachingForOptionsEnabled)
	flags.StringVar(&fields.enableCachingForPost, "enable-caching-for-post", "", msg.FlagCachingForPostEnabled)
	flags.StringVar(&fields.enableQueryStringSort, "enable-caching-string-sort", "", msg.FlagCachingStringSortEnabled)
	flags.StringVar(&fields.isSliceConfigurationEnabled, "slice-configuration-enabled", "", msg.FlagSliceConfigurationEnabled)
	flags.Int64Var(&fields.sliceConfigurationRange, "slice-configuration-range", 0, msg.FlagSliceConfigurationRange)
	flags.StringVar(&fields.adaptiveDeliveryAction, "adaptive-delivery-action", "ignore", msg.FlagAdaptiveDeliveryAction)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}

func appAccelerationNoEnabled(client *apiEdgeApp.Client, fields *Fields, request api.Request) error {
	ctx := context.Background()
	str := strconv.FormatInt(fields.ApplicationID, 10)
	application, err := client.Get(ctx, str)
	if err != nil {
		return err
	}

	acc := application.GetModules()

	enabled := request.GetEdgeCache()

	if (enabled.GetCachingForOptionsEnabled() ||
		enabled.GetCachingForPostEnabled()) &&
		!acc.GetApplicationAcceleratorEnabled() {
		return msg.ErrorApplicationAccelerationNotEnabled
	}
	return nil
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.Request) error {
	request.SetName(fields.Name)
	if cmd.Flags().Changed("browser-cache-behavior") {
		if fields.browserCacheBehavior == "override" && !cmd.Flags().Changed("browser-cache-settings-maximum-ttl") {
			return msg.ErrorBrowserMaximumTtlNotSent
		}

		req := sdk.BrowserCacheModuleRequest{
			Behavior: fields.browserCacheBehavior,
			MaxAge:   fields.browserCacheMaxAge,
		}
		request.SetBrowserCache(req)
	}

	if cmd.Flags().Changed("query-string-fields") {
		controls := request.GetApplicationControls()
		controls.SetQueryStringFields(fields.queryStringFields)
	}

	if cmd.Flags().Changed("cookie-names") {
		controls := request.GetApplicationControls()
		controls.SetCookieNames(fields.cookieNames)
	}

	if cmd.Flags().Changed("cache-by-cookies") {
		controls := request.GetApplicationControls()
		controls.SetCacheByCookies(fields.cacheByCookies)
	}

	if cmd.Flags().Changed("cache-by-query-string") {
		controls := request.GetApplicationControls()
		controls.SetCacheByQueryString(fields.cacheByQueryString)
	}

	if cmd.Flags().Changed("slice-configuration-range") {
		controls := request.GetSliceControls()
		controls.SetSliceConfigurationRange(fields.sliceConfigurationRange)
	}

	if cmd.Flags().Changed("adaptive-delivery-action") {
		controls := request.GetApplicationControls()
		controls.SetAdaptiveDeliveryAction(fields.adaptiveDeliveryAction)
	}

	if cmd.Flags().Changed("enable-caching-for-options") {
		cachingOptions, err := strconv.ParseBool(fields.enableCachingForOptions)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorCachingForOptionsFlag, fields.enableCachingForOptions)
		}

		edgeCache := request.GetEdgeCache()
		edgeCache.SetCachingForOptionsEnabled(cachingOptions)
	}

	if cmd.Flags().Changed("enable-caching-for-post") {
		cachingPost, err := strconv.ParseBool(fields.enableCachingForPost)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorCachingForPostFlag, fields.enableCachingForPost)
		}

		edgeCache := request.GetEdgeCache()
		edgeCache.SetCachingForPostEnabled(cachingPost)
	}

	if cmd.Flags().Changed("enable-caching-string-sort") {
		stringSort, err := strconv.ParseBool(fields.enableQueryStringSort)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorCachingStringSortFlag, fields.enableQueryStringSort)
		}

		controls := request.GetApplicationControls()
		controls.SetQueryStringSortEnabled(stringSort)
	}

	if cmd.Flags().Changed("slice-configuration-enabled") {
		sliceEnable, err := strconv.ParseBool(fields.isSliceConfigurationEnabled)
		if err != nil {
			return fmt.Errorf("%w: %q", msg.ErrorSliceConfigurationFlag, fields.isSliceConfigurationEnabled)
		}

		controls := request.GetSliceControls()
		controls.SetSliceConfigurationEnabled(sliceEnable)
		controls.SetSliceEdgeCachingEnabled(true) //Edge Cache is mandatory in this case
	}

	return nil
}
