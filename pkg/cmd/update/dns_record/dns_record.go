package dnsrecord

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/dns_record"
	api "github.com/aziontech/azion-cli/pkg/api/dns_record"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	ZoneID      int64
	RecordID    int64
	Name        string
	Type        string
	Rdata       []string
	TTL         int64
	Policy      string
	Weight      int64
	Description string
	Path        string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DNSRecordUsage,
		Short:         msg.DNSRecordUpdateShortDescription,
		Long:          msg.DNSRecordUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion update dns-record --zone-id 107313 --record-id 56789 --rdata "192.0.2.3"
        $ azion update dns-record --zone-id 107313 --record-id 56789 --name "www" --ttl 7200
        $ azion update dns-record --zone-id 107313 --record-id 56789 --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := sdk.PatchedRecordRequest{}

			if !cmd.Flags().Changed("zone-id") {
				answers, err := utils.AskInput(msg.DNSRecordUpdateAskInputZoneID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				zoneID, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return msg.ErrorConvertIdZone
				}

				fields.ZoneID = int64(zoneID)
			}

			if !cmd.Flags().Changed("record-id") {
				answers, err := utils.AskInput(msg.DNSRecordUpdateAskInputRecordID)
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				recordID, err := strconv.Atoi(answers)
				if err != nil {
					logger.Debug("Error while parsing string to integer", zap.Error(err))
					return msg.ErrorConvertIdRecord
				}

				fields.RecordID = int64(recordID)
			}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				// The DNS record update is a full replacement (PUT); name, type
				// and rdata are required. Fetch the current record and preserve
				// any attribute the user didn't provide.
				current, err := client.Get(context.Background(), fields.ZoneID, fields.RecordID)
				if err != nil {
					return fmt.Errorf(msg.ErrorGetDNSRecord.Error(), err)
				}

				if !cmd.Flags().Changed("name") {
					fields.Name = current.GetName()
				}
				if !cmd.Flags().Changed("type") {
					fields.Type = current.GetType()
				}
				if !cmd.Flags().Changed("rdata") {
					fields.Rdata = current.GetRdata()
				}

				request.SetName(fields.Name)
				request.SetType(fields.Type)
				request.SetRdata(fields.Rdata)

				if cmd.Flags().Changed("ttl") {
					request.SetTtl(fields.TTL)
				} else if current.HasTtl() {
					request.SetTtl(current.GetTtl())
				}

				if cmd.Flags().Changed("policy") {
					request.SetPolicy(fields.Policy)
				} else if current.HasPolicy() {
					fields.Policy = current.GetPolicy()
					request.SetPolicy(fields.Policy)
				}

				// Weight and description are only relevant to the 'weighted' policy.
				if fields.Policy == "weighted" {
					if cmd.Flags().Changed("weight") {
						request.SetWeight(fields.Weight)
					} else if current.HasWeight() {
						request.SetWeight(current.GetWeight())
					}

					if cmd.Flags().Changed("description") {
						request.SetDescription(fields.Description)
					} else if current.HasDescription() {
						request.SetDescription(current.GetDescription())
					}
				}
			}

			response, err := client.Update(context.Background(), request, fields.ZoneID, fields.RecordID)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateDNSRecord.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DNSRecordUpdateOutputSuccess, response.GetId()),
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
	flags.Int64Var(&fields.ZoneID, "zone-id", 0, msg.DNSRecordFlagZoneID)
	flags.Int64Var(&fields.RecordID, "record-id", 0, msg.DNSRecordFlagRecordID)
	flags.StringVar(&fields.Name, "name", "", msg.DNSRecordUpdateFlagName)
	flags.StringVar(&fields.Type, "type", "", msg.DNSRecordUpdateFlagType)
	flags.StringSliceVar(&fields.Rdata, "rdata", []string{}, msg.DNSRecordUpdateFlagRdata)
	flags.Int64Var(&fields.TTL, "ttl", 0, msg.DNSRecordUpdateFlagTTL)
	flags.StringVar(&fields.Policy, "policy", "", msg.DNSRecordUpdateFlagPolicy)
	flags.Int64Var(&fields.Weight, "weight", 0, msg.DNSRecordUpdateFlagWeight)
	flags.StringVar(&fields.Description, "description", "", msg.DNSRecordUpdateFlagDescription)
	flags.StringVar(&fields.Path, "file", "", msg.DNSRecordUpdateFlagIn)
	flags.BoolP("help", "h", false, msg.DNSRecordUpdateHelpFlag)
}
