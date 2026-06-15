package dnszone

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/dns_zone"
	api "github.com/aziontech/azion-cli/pkg/api/dns_zone"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	Name   string
	Domain string
	Active bool
	Path   string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DNSZoneUsage,
		Short:         msg.DNSZoneCreateShortDescription,
		Long:          msg.DNSZoneCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create dns-zone --name "my zone" --domain "example.com" --active true
        $ azion create dns-zone --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := sdk.ZoneRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("name") {
					answers, err := utils.AskInput(msg.DNSZoneCreateAskInputName)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					fields.Name = answers
				}

				if !cmd.Flags().Changed("domain") {
					answers, err := utils.AskInput(msg.DNSZoneCreateAskInputDomain)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					fields.Domain = answers
				}

				request.SetName(fields.Name)
				request.SetDomain(fields.Domain)
				request.SetActive(fields.Active)
			}

			response, err := client.Create(context.Background(), request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateDNSZone.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DNSZoneCreateOutputSuccess, response.GetId()),
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
	flags.StringVar(&fields.Name, "name", "", msg.DNSZoneCreateFlagName)
	flags.StringVar(&fields.Domain, "domain", "", msg.DNSZoneCreateFlagDomain)
	flags.BoolVar(&fields.Active, "active", true, msg.DNSZoneCreateFlagActive)
	flags.StringVar(&fields.Path, "file", "", msg.DNSZoneCreateFlagIn)
	flags.BoolP("help", "h", false, msg.DNSZoneCreateHelpFlag)
}
