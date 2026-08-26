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
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var (
	zoneID   int64
	recordID int64
)

type DeleteCmd struct {
	Io              *iostreams.IOStreams
	DeleteDNSRecord func(context.Context, int64, int64) error
	AskInput        func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		DeleteDNSRecord: func(ctx context.Context, zoneID, recordID int64) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Delete(ctx, zoneID, recordID)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(del *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.DNSRecordUsage,
		Short:         msg.DNSRecordDeleteShortDescription,
		Long:          msg.DNSRecordDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete dns-record --zone-id 107313 --record-id 56789
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("zone-id") {
				answer, err := del.AskInput(msg.DNSRecordDeleteAskInputZoneID)
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
				answer, err := del.AskInput(msg.DNSRecordDeleteAskInputRecordID)
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

			err := del.DeleteDNSRecord(ctx, zoneID, recordID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DNSRecordDeleteOutputSuccess, recordID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&zoneID, "zone-id", 0, msg.DNSRecordFlagZoneID)
	cobraCmd.Flags().Int64Var(&recordID, "record-id", 0, msg.DNSRecordFlagRecordID)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DNSRecordDeleteHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
