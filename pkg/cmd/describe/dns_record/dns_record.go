package dnsrecord

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/dns_record"
	api "github.com/aziontech/azion-cli/pkg/api/dns_record"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
)

var (
	zoneID   int64
	recordID int64
)

type DescribeCmd struct {
	Io       *iostreams.IOStreams
	AskInput func(string) (string, error)
	Get      func(context.Context, int64, int64) (sdk.Record, error)
}

func NewDescribeCmd(f *cmdutil.Factory) *DescribeCmd {
	return &DescribeCmd{
		Io: f.IOStreams,
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		Get: func(ctx context.Context, zoneID, recordID int64) (sdk.Record, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Get(ctx, zoneID, recordID)
		},
	}
}

func NewCobraCmd(describe *DescribeCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.DescribeOptions{}
	cobraCmd := &cobra.Command{
		Use:           msg.DNSRecordUsage,
		Short:         msg.DNSRecordDescribeShortDescription,
		Long:          msg.DNSRecordDescribeLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion describe dns-record --zone-id 107313 --record-id 56789
        $ azion describe dns-record --zone-id 107313 --record-id 56789 --format json
        $ azion describe dns-record --zone-id 107313 --record-id 56789 --out "./tmp/test.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("zone-id") {
				answer, err := describe.AskInput(msg.DNSRecordDescribeAskInputZoneID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdZone
				}

				zoneID = num
			}

			if !cmd.Flags().Changed("record-id") {
				answer, err := describe.AskInput(msg.DNSRecordDescribeAskInputRecordID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdRecord
				}

				recordID = num
			}

			ctx := context.Background()
			resp, err := describe.Get(ctx, zoneID, recordID)
			if err != nil {
				return fmt.Errorf(msg.ErrorGetDNSRecord.Error(), err)
			}

			fields := make(map[string]string, 0)
			fields["Id"] = "ID"
			fields["Name"] = "Name"
			fields["Type"] = "Type"
			fields["Rdata"] = "Rdata"
			fields["Ttl"] = "TTL"
			fields["Policy"] = "Policy"
			fields["Weight"] = "Weight"
			fields["Description"] = "Description"

			describeOut := output.DescribeOutput{
				GeneralOutput: output.GeneralOutput{
					Out:   f.IOStreams.Out,
					Msg:   filepath.Clean(opts.OutPath),
					Flags: f.Flags,
				},
				Fields: fields,
				Values: &resp,
			}
			return output.Print(&describeOut)
		},
	}

	cobraCmd.Flags().Int64Var(&zoneID, "zone-id", 0, msg.DNSRecordFlagZoneID)
	cobraCmd.Flags().Int64Var(&recordID, "record-id", 0, msg.DNSRecordFlagRecordID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DNSRecordDescribeHelpFlag)
	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDescribeCmd(f), f)
}
