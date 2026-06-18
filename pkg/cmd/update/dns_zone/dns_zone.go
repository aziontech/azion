package dnszone

import (
	"context"
	"fmt"
	"strconv"

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
	ZoneID int64
	Name   string
	Active bool
	Path   string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DNSZoneUsage,
		Short:         msg.DNSZoneUpdateShortDescription,
		Long:          msg.DNSZoneUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion update dns-zone --zone-id 12312 --name "my zone"
        $ azion update dns-zone --zone-id 12312 --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := sdk.PatchedUpdateZoneRequest{}

			if !cmd.Flags().Changed("zone-id") {
				answers, err := utils.AskInput(msg.DNSZoneUpdateAskInputZoneID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				id, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return msg.ErrorConvertIdZone
				}

				fields.ZoneID = int64(id)
			}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				// The DNS zone update is a full replacement (PUT), and both
				// name and active are required. Fetch the current values for
				// any attribute the user didn't provide so they are preserved.
				if !cmd.Flags().Changed("name") || !cmd.Flags().Changed("active") {
					current, err := client.Get(context.Background(), fields.ZoneID)
					if err != nil {
						return fmt.Errorf(msg.ErrorGetDNSZone.Error(), err)
					}

					if !cmd.Flags().Changed("name") {
						fields.Name = current.GetName()
					}

					if !cmd.Flags().Changed("active") {
						fields.Active = current.GetActive()
					}
				}

				request.SetName(fields.Name)
				request.SetActive(fields.Active)
			}

			response, err := client.Update(context.Background(), request, fields.ZoneID)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateDNSZone.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DNSZoneUpdateOutputSuccess, response.GetId()),
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
	flags.Int64Var(&fields.ZoneID, "zone-id", 0, msg.DNSZoneFlagId)
	flags.StringVar(&fields.Name, "name", "", msg.DNSZoneUpdateFlagName)
	flags.BoolVar(&fields.Active, "active", true, msg.DNSZoneUpdateFlagActive)
	flags.StringVar(&fields.Path, "file", "", msg.DNSZoneUpdateFlagIn)
	flags.BoolP("help", "h", false, msg.DNSZoneUpdateHelpFlag)
}
