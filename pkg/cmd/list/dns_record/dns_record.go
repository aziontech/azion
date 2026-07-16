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
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
)

type ListCmd struct {
	Io            *iostreams.IOStreams
	ListDNSRecord func(context.Context, *contracts.ListOptions, int64) (*sdk.PaginatedRecordList, error)
	AskInput      func(string) (string, error)
	ZoneID        int64
}

func NewListCmd(f *cmdutil.Factory) *ListCmd {
	return &ListCmd{
		Io: f.IOStreams,
		ListDNSRecord: func(ctx context.Context, opts *contracts.ListOptions, zoneID int64) (*sdk.PaginatedRecordList, error) {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.List(ctx, opts, zoneID)
		},
		AskInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
	}
}

func NewCobraCmd(list *ListCmd, f *cmdutil.Factory) *cobra.Command {
	opts := &contracts.ListOptions{}
	cmd := &cobra.Command{
		Use:           msg.DNSRecordUsage,
		Short:         msg.DNSRecordListShortDescription,
		Long:          msg.DNSRecordListLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
			$ azion list dns-record --zone-id 107313
			$ azion list dns-record --zone-id 107313 --details
			$ azion list dns-record --zone-id 107313 --order-by "name"
			$ azion list dns-record --zone-id 107313 --order-by "-name"
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("zone-id") {
				answer, err := list.AskInput(msg.DNSRecordListAskInputZoneID)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ErrorConvertIdZone
				}

				list.ZoneID = num
			}

			if err := PrintTable(cmd, f, opts, list); err != nil {
				return fmt.Errorf(msg.ErrorListDNSRecord.Error(), err)
			}
			return nil
		},
	}

	cmdutil.AddAzionApiFlags(cmd, opts)
	cmd.Flags().Int64Var(&list.ZoneID, "zone-id", 0, msg.DNSRecordFlagZoneID)
	cmd.Flags().BoolP("help", "h", false, msg.DNSRecordListHelpFlag)

	return cmd
}

func PrintTable(cmd *cobra.Command, f *cmdutil.Factory, opts *contracts.ListOptions, list *ListCmd) error {
	ctx := context.Background()

	response, err := list.ListDNSRecord(ctx, opts, list.ZoneID)
	if err != nil {
		return err
	}

	listOut := output.ListOutput{}
	listOut.Columns = []string{"ID", "NAME", "TYPE", "TTL", "RDATA"}
	listOut.Out = f.IOStreams.Out
	listOut.Flags = f.Flags

	for _, v := range response.GetResults() {
		ln := []string{
			fmt.Sprintf("%d", v.GetId()),
			v.GetName(),
			v.GetType(),
			fmt.Sprintf("%d", v.GetTtl()),
			strings.Join(v.GetRdata(), ", "),
		}
		listOut.Lines = append(listOut.Lines, ln)
	}

	return output.Print(&listOut)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewListCmd(f), f)
}
