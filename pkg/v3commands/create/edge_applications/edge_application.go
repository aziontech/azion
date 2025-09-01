package edge_applications

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/edge_application"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/edge_applications"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const example = `
        $ azion create edge-application --name "naruno"
        $ azion create edge-application --file create.json
        $ json example to be used with '--file flag' "create.json": 
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
	Name                           string `json:"name,omitempty"`
	DeliveryProtocol               string `json:"delivery_protocol,omitempty"`
	Http3                          string `json:"http_3,omitempty"`
	HttpPort                       string `json:"http_port,omitempty"`
	HttpsPort                      string `json:"https_port,omitempty"`
	OriginType                     string `json:"origin_type,omitempty"`
	Address                        string `json:"address,omitempty"`
	OriginProtocolPolicy           string `json:"origin_protocol_policy,omitempty"`
	HostHeader                     string `json:"host_header,omitempty"`
	BrowserCacheSettings           string `json:"browser_cache_settings,omitempty"`
	CdnCacheSettings               string `json:"cdn_cache_settings,omitempty"`
	BrowserCacheSettingsMaximumTtl int64  `json:"browser_cache_settings_maximum_ttl,omitempty"`
	CdnCacheSettingsMaximumTtl     int64  `json:"cdn_cache_settings_maximum_ttl,omitempty"`
	DebugRules                     string `json:"debug_rules,omitempty"`
	SupportedCiphers               string `json:"supported_ciphers,omitempty"`
	Websocket                      string `json:"websocket,omitempty"`
	Path                           string
	OutPath                        string
	Format                         string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           "edge-application",
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(example),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.CreateRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
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

			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OutputSuccess, response.GetId()),
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

func createRequestFromFlags(fields *Fields, request *api.CreateRequest) error {

	if utils.IsEmpty(fields.Name) {
		answers, err := utils.AskInput("Enter the new Edge Application's name")
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

	if !utils.IsEmpty(fields.HttpsPort) {
		request.SetHttpsPort(fields.HttpsPort)
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
	flags.StringVar(&fields.HttpsPort, "https-port", "", msg.FlagHttpsPort)
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
	flags.StringVar(&fields.Path, "file", "", msg.FlagFile)
	flags.BoolP("help", "h", false, msg.FlagHelp)
}
