package cachesetting

import (
	"context"
	"fmt"
	"strconv"

	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/cache_setting"
	apiEdgeApp "github.com/aziontech/azion-cli/pkg/api/applications"
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
	Name                    string
	browserCacheBehavior    string
	browserCacheMaxAge      int64
	adaptiveDeliveryAction  string
	cacheByQueryString      string
	queryStringFields       []string
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
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create cache-setting --application-id 1673635839 --name "phototypesetting"
        $ azion create cache-setting --application-id 1673635839 --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClientV4(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			clientEdgeApp := apiEdgeApp.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := sdk.CacheSettingRequest{}

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

			if err := appAccelerationNotEnabled(clientEdgeApp, fields, &request); err != nil {
				return err
			}

			response, err := client.Create(context.Background(), request, fields.ApplicationID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateCacheSettings.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.CreateOutputSuccess, response.GetId()),
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
	flags.Int64Var(&fields.ApplicationID, "application-id", 0, msg.FlagApplicationID)
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.browserCacheBehavior, "browser-cache-behavior", "honor", msg.FlagBrowserCacheBehavior)
	flags.Int64Var(&fields.browserCacheMaxAge, "browser-cache-max-age", 0, msg.FlagBrowserCacheMaxAge)
	flags.StringSliceVar(&fields.queryStringFields, "query-string-fields", []string{}, msg.FlagQueryStringFields)
	flags.StringSliceVar(&fields.cookieNames, "cookie-names", []string{}, msg.FlagCookieNames)
	flags.StringVar(&fields.cacheByCookies, "cache-by-cookies", "ignore", msg.FlagCacheByCookiesEnabled)
	flags.StringVar(&fields.cacheByQueryString, "cache-by-query-string", "ignore", msg.FlagCacheByQueryString)
	flags.StringVar(&fields.enableCachingForOptions, "enable-caching-for-options", "false", msg.FlagCachingForOptionsEnabled)
	flags.StringVar(&fields.enableCachingForPost, "enable-caching-for-post", "", msg.FlagCachingForPostEnabled)
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}

func appAccelerationNotEnabled(client *apiEdgeApp.Client, fields *Fields, request *sdk.CacheSettingRequest) error {
	ctx := context.Background()
	application, err := client.Get(ctx, fields.ApplicationID)
	if err != nil {
		return err
	}

	acc := application.GetModules()
	appAcc := acc.GetApplicationAccelerator()

	if request.GetModules().EdgeCache != nil {
		edgeCache := request.GetModules().ApplicationAccelerator
		if len(edgeCache.CacheVaryByMethod) > 0 && !appAcc.GetEnabled() {
			return msg.ErrorApplicationAccelerationNotEnabled
		}
	}

	return nil
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *sdk.CacheSettingRequest) error {
	request.SetName(fields.Name)
	if cmd.Flags().Changed("browser-cache-behavior") {
		if fields.browserCacheBehavior == "override" && !cmd.Flags().Changed("browser-cache-max-age") {
			return msg.ErrorBrowserMaximumTtlNotSent
		}

		req := sdk.BrowserCacheModuleRequest{}
		req.SetBehavior(fields.browserCacheBehavior)
		req.SetMaxAge(fields.browserCacheMaxAge)
		request.SetBrowserCache(req)
	}

	modules := &sdk.CacheSettingsModulesRequest{}
	appAcc := &sdk.CacheSettingsApplicationAcceleratorModuleRequest{}
	setAcc := false

	if cmd.Flags().Changed("query-string-fields") && cmd.Flags().Changed("cache-by-query-string") {
		cacheVary := &sdk.CacheVaryByQuerystringModuleRequest{}
		cacheVary.SetFields(fields.queryStringFields)
		cacheVary.SetBehavior(fields.cacheByQueryString)
		appAcc.SetCacheVaryByQuerystring(*cacheVary)
		setAcc = true
	}

	if cmd.Flags().Changed("cookie-names") && cmd.Flags().Changed("cache-by-cookies") {
		controls := &sdk.CacheVaryByCookiesModuleRequest{}
		controls.SetCookieNames(fields.cookieNames)
		controls.SetBehavior(fields.cacheByCookies)
		appAcc.SetCacheVaryByCookies(*controls)
		setAcc = true
	}

	if cmd.Flags().Changed("enable-caching-for-options") || cmd.Flags().Changed("enable-caching-for-post") {
		cacheByMethod := []string{}
		if fields.enableCachingForOptions != "" {
			cacheByMethod = append(cacheByMethod, "options")
		}

		if fields.enableCachingForPost != "" {
			cacheByMethod = append(cacheByMethod, "post")
		}

		appAcc.SetCacheVaryByMethod(cacheByMethod)
		setAcc = true

	}

	if setAcc {
		modules.SetApplicationAccelerator(*appAcc)
		request.SetModules(*modules)
	}

	return nil
}
