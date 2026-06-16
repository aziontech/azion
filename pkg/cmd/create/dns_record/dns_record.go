package dnsrecord

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
		Short:         msg.DNSRecordCreateShortDescription,
		Long:          msg.DNSRecordCreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create dns-record --zone-id 107313 --name "www" --type "A" --rdata "192.0.2.1" --ttl 3600
        $ azion create dns-record --zone-id 107313 --name "www" --type "A" --rdata "192.0.2.1,192.0.2.2"
        $ azion create dns-record --zone-id 107313 --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := sdk.RecordRequest{}

			if !cmd.Flags().Changed("zone-id") {
				answers, err := utils.AskInput(msg.DNSRecordCreateAskInputZoneID)
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

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.Path, &request)
				if err != nil {
					logger.Debug("Error while parsing <"+fields.Path+"> file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
			} else {
				if !cmd.Flags().Changed("name") {
					answers, err := utils.AskInput(msg.DNSRecordCreateAskInputName)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					fields.Name = answers
				}

				if !cmd.Flags().Changed("type") {
					answers, err := utils.AskInput(msg.DNSRecordCreateAskInputType)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					fields.Type = answers
				}

				if !cmd.Flags().Changed("rdata") {
					answers, err := utils.AskInput(msg.DNSRecordCreateAskInputRdata)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					fields.Rdata = splitRdata(answers)
				}

				request.SetName(fields.Name)
				request.SetType(fields.Type)
				request.SetRdata(fields.Rdata)

				if cmd.Flags().Changed("ttl") {
					request.SetTtl(fields.TTL)
				}

				if cmd.Flags().Changed("policy") {
					request.SetPolicy(fields.Policy)
				}

				// Weight and description are only relevant to the 'weighted' policy.
				if fields.Policy == "weighted" {
					if cmd.Flags().Changed("weight") {
						request.SetWeight(fields.Weight)
					}
					if cmd.Flags().Changed("description") {
						request.SetDescription(fields.Description)
					}
				}
			}

			response, err := client.Create(context.Background(), request, fields.ZoneID)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateDNSRecord.Error(), err)
			}

			creatOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DNSRecordCreateOutputSuccess, response.GetId()),
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
	flags.Int64Var(&fields.ZoneID, "zone-id", 0, msg.DNSRecordFlagZoneID)
	flags.StringVar(&fields.Name, "name", "", msg.DNSRecordCreateFlagName)
	flags.StringVar(&fields.Type, "type", "", msg.DNSRecordCreateFlagType)
	flags.StringSliceVar(&fields.Rdata, "rdata", []string{}, msg.DNSRecordCreateFlagRdata)
	flags.Int64Var(&fields.TTL, "ttl", 0, msg.DNSRecordCreateFlagTTL)
	flags.StringVar(&fields.Policy, "policy", "", msg.DNSRecordCreateFlagPolicy)
	flags.Int64Var(&fields.Weight, "weight", 0, msg.DNSRecordCreateFlagWeight)
	flags.StringVar(&fields.Description, "description", "", msg.DNSRecordCreateFlagDescription)
	flags.StringVar(&fields.Path, "file", "", msg.DNSRecordCreateFlagIn)
	flags.BoolP("help", "h", false, msg.DNSRecordCreateHelpFlag)
}

// splitRdata parses a comma-separated answer into a trimmed list of values.
func splitRdata(answer string) []string {
	parts := strings.Split(answer, ",")
	rdata := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			rdata = append(rdata, trimmed)
		}
	}
	return rdata
}
