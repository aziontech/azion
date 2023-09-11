package edge_application

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/create/edge_application"
	api "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/logger"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type Fields struct {
	Name                           string `json:"name"`
	ApplicationAcceleration        string `json:"application_acceleration,omitempty"`
	DeliveryProtocol               string `json:"delivery_protocol,omitempty"`
	OriginType                     string `json:"origin_type,omitempty"`
	Address                        string `json:"address,omitempty"`
	OriginProtocolPolicy           string `json:"origin_protocol_policy,omitempty"`
	HostHeader                     string `json:"host_header,omitempty"`
	BrowserCacheSettings           string `json:"browser_cache_settings,omitempty"`
	CdnCacheSettings               string `json:"cdn_cache_settings,omitempty"`
	BrowserCacheSettingsMaximumTtl int64  `json:"browser_cache_settings_maximum_ttl,omitempty"`
	CdnCacheSettingsMaximumTtl     int64  `json:"cdn_cache_settings_maximum_ttl,omitempty"`
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
		Example: heredoc.Doc(`
        $ azion create --name "naruno"
        $ azion create edge_applications --in create.json
        $ json example "create.json": 
        {
            "name": "New Edge Application",
            "delivery_protocol": "http",
            "origin_type": "single_origin",
            "address": "www.new.api",
            "origin_protocol_policy": "preserve",
            "host_header": "www.new.api",
            "browser_cache_settings": "override",
            "browser_cache_settings_maximum_ttl": 20,
            "cdn_cache_settings": "honor",
            "cdn_cache_settings_maximum_ttl": 60
        }
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.CreateRequest{}

			if cmd.Flags().Changed("in") {
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

				err = cmdutil.UnmarshallJsonFromReader(file, &fields)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}

			}

			if utils.IsEmpty(fields.Name) {
				qs := []*survey.Question{
					{
						Name:      "name",
						Prompt:    &survey.Input{Message: "What is the name on Edge Application? required!"},
						Validate:  survey.Required,
						Transform: survey.Title,
					},
				}

				answers := struct{ Name string }{}

				err := survey.Ask(qs, &answers)
				if err != nil {
					return err
				}

				fields.Name = answers.Name
			}

			if utils.IsEmpty(fields.Name) {
				return msg.ErrorMandatoryCreateFlags
			}

			request.SetName(fields.Name)

			if !utils.IsEmpty(fields.ApplicationAcceleration) {
				applicationAcceleration, err := strconv.ParseBool(fields.ApplicationAcceleration)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorParseBoolToString
				}
				request.SetApplicationAcceleration(applicationAcceleration)
			}

			if !utils.IsEmpty(fields.DeliveryProtocol) {
				request.SetDeliveryProtocol(fields.DeliveryProtocol)
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

			response, err := api.NewClient(
				f.HttpClient,
				f.Config.GetString("api_url"),
				f.Config.GetString("token"),
			).Create(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreate.Error(), err)
			}

			logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.OutputSuccess, response.GetId()))

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&fields.Name, "name", "", msg.FlagName)
	flags.StringVar(&fields.ApplicationAcceleration, "ApplicationAcceleration", "", msg.FlagApplicationAcceleration)
	flags.StringVar(&fields.DeliveryProtocol, "DeliveryProtocol", "", msg.FlagDeliveryProtocol)
	flags.StringVar(&fields.OriginType, "OriginType", "", msg.FlagOriginType)
	flags.StringVar(&fields.Address, "address", "", msg.FlagAddress)
	flags.StringVar(&fields.OriginProtocolPolicy, "OriginProtocolPolicy", "", msg.FlagOriginProtocolPolicy)
	flags.StringVar(&fields.BrowserCacheSettings, "BrowserCacheSettings", "", msg.FlagBrowserCacheSettings)
	flags.StringVar(&fields.CdnCacheSettings, "CdnCacheSettings", "", msg.FlagCdnCacheSettings)
	flags.Int64Var(&fields.BrowserCacheSettingsMaximumTtl, "BrowserCacheSettingsMaximumTtl", 0, msg.FlagBrowserCacheSettingsMaximumTtl)
	flags.Int64Var(&fields.CdnCacheSettingsMaximumTtl, "CdnCacheSettingsMaximumTtl", 0, msg.FlagCdnCacheSettingsMaximumTtl)

	flags.StringVar(&fields.Path, "in", "", msg.FlagIn)
	flags.BoolP("help", "h", false, msg.FlagHelp)

	return cmd
}
