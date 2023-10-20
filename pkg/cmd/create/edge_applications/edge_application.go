package edge_applications

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/edge_application"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/logger"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const example = `
        $ azion create edge-applications --name "naruno"
        $ azion create edge-applications --in create.json
        $ json example "create.json": 
        {
            "name": "New Edge Application",
            "delivery_protocol": "http",
            "origin_type": "single_origin",
            "address": "www.new.api",
            "origin_protocol_policy": "preserve",
            "host_header": "${host}",
            "browser_cache_settings": "override",
            "browser_cache_settings_maximum_ttl": 20,
            "cdn_cache_settings": "honor",
            "cdn_cache_settings_maximum_ttl": 60
        }
        `

type Fields struct {
	Name                           string
	DeliveryProtocol               string
	Http3                          string
	HttpPort                       string
	OriginType                     string
	Address                        string
	OriginProtocolPolicy           string
	HostHeader                     string
	BrowserCacheSettings           string
	CdnCacheSettings               string
	BrowserCacheSettingsMaximumTtl int64
	CdnCacheSettingsMaximumTtl     int64
	DebugRules                     string
	SupportedCiphers               string
	Websocket                      string
	Path                           string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(example),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.CreateRequest{}

			if cmd.Flags().Changed("in") {
				err := utils.FlagINUnmarshalFileJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				err := createRequestFromFlags(fields, &request)
				if err != nil {
					return err
				}
			}

			response, err := api.NewClient(
				f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"),
			).Create(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}

			logger.LogSuccess(f.IOStreams.Out, fmt.Sprintf(msg.OutputSuccess, response.GetId()))

			return nil
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func createRequestFromFlags(fields *Fields, request *api.CreateRequest) error {
	if utils.IsEmpty(fields.Name) {
		answers, err := utils.AskInput("What is the name of the Edge Application?")
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}

		fields.Name = answers
	}

	if utils.IsEmpty(fields.Name) {
		return msg.ErrorMandatoryCreateFlags
	}

	request.SetName(fields.Name)

	if !utils.IsEmpty(fields.DeliveryProtocol) {
		request.SetDeliveryProtocol(fields.DeliveryProtocol)
	}

	if !utils.IsEmpty(fields.Http3) {
		http3, err := strconv.ParseBool(fields.Http3)
		if err != nil {
			logger.Debug("Error while parsing <"+fields.Http3+"> ", zap.Error(err))
			return utils.ErrorConvertingStringToBool
		}

		request.SetHttp3(http3)
	}

	if !utils.IsEmpty(fields.DebugRules) {
		debugRules, err := strconv.ParseBool(fields.DebugRules)
		if err != nil {
			logger.Debug("Error while parsing <"+fields.DebugRules+"> ", zap.Error(err))
			return utils.ErrorConvertingStringToBool
		}

		request.SetDebugRules(debugRules)
	}

	if !utils.IsEmpty(fields.SupportedCiphers) {
		request.SetSupportedCiphers(fields.SupportedCiphers)
	}

	if !utils.IsEmpty(fields.Websocket) {
		websocket, err := strconv.ParseBool(fields.Websocket)
		if err != nil {
			logger.Debug("Error while parsing <"+fields.Websocket+"> ", zap.Error(err))
			return utils.ErrorConvertingStringToBool
		}

		request.SetWebsocket(websocket)
	}

	if !utils.IsEmpty(fields.HttpPort) {
		request.SetHttpPort(fields.HttpPort)
	}

	if !utils.IsEmpty(fields.DeliveryProtocol) {
		request.SetOriginType(fields.OriginType)
	}

	if !utils.IsEmpty(fields.Address) {
		request.SetAddress(fields.Address)
	}

	if !utils.IsEmpty(fields.OriginProtocolPolicy) {
		request.SetOriginProtocolPolicy(fields.OriginProtocolPolicy)
	}

	if !utils.IsEmpty(fields.HostHeader) {
		request.SetHostHeader(fields.HostHeader)
	}

	if !utils.IsEmpty(fields.BrowserCacheSettings) {
		request.SetBrowserCacheSettings(fields.BrowserCacheSettings)
	}

	if !utils.IsEmpty(fields.CdnCacheSettings) {
		request.SetCdnCacheSettings(fields.CdnCacheSettings)
	}

	if fields.BrowserCacheSettingsMaximumTtl <= 0 {
		request.SetBrowserCacheSettingsMaximumTtl(fields.BrowserCacheSettingsMaximumTtl)
	}

	if fields.CdnCacheSettingsMaximumTtl <= 0 {
		request.SetCdnCacheSettingsMaximumTtl(fields.CdnCacheSettingsMaximumTtl)
	}

	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.DeliveryProtocol, "delivery-protocol", "", msg.FlagDeliveryProtocol)
	flags.StringVar(&fields.Http3, "http3", "", msg.FlagHttp3)
	flags.StringVar(&fields.HttpPort, "http-port", "", msg.FlagHttpPort)
	flags.StringVar(&fields.OriginType, "origin-type", "", msg.FlagOriginType)
	flags.StringVar(&fields.Address, "address", "", msg.FlagAddress)
	flags.StringVar(&fields.OriginProtocolPolicy, "origin-protocol-policy", "", msg.FlagOriginProtocolPolicy)
	flags.StringVar(&fields.HostHeader, "host-header", "", msg.FlagHostHeader)
	flags.StringVar(&fields.BrowserCacheSettings, "browser-cache-settings", "", msg.FlagBrowserCacheSettings)
	flags.Int64Var(&fields.BrowserCacheSettingsMaximumTtl, "browser-cache-settings-maximum-ttl", 0, msg.FlagBrowserCacheSettingsMaximumTtl)
	flags.StringVar(&fields.CdnCacheSettings, "cdn-cache-settings", "", msg.FlagCdnCacheSettings)
	flags.StringVar(&fields.DebugRules, "debug-rules", "", msg.FlagDebugRules)
	flags.StringVar(&fields.SupportedCiphers, "supported-ciphers", "", msg.FlagSupportedCiphers)
	flags.StringVar(&fields.Websocket, "websocket", "", msg.FlagWebsocket)
	flags.Int64Var(&fields.CdnCacheSettingsMaximumTtl, "cdn-cache-settings-maximum-ttl", 0, msg.FlagCdnCacheSettingsMaximumTtl)
	flags.StringVar(&fields.Path, "in", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.FlagHelp)
}
