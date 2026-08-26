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
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var zoneID int64

type DeleteCmd struct {
	Io            *iostreams.IOStreams
	DeleteDNSZone func(context.Context, int64) error
	AskInput      func(string) (string, error)
}

func NewDeleteCmd(f *cmdutil.Factory) *DeleteCmd {
	return &DeleteCmd{
		Io: f.IOStreams,
		DeleteDNSZone: func(ctx context.Context, id int64) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			return client.Delete(ctx, id)
		},
		AskInput: utils.AskInput,
	}
}

func NewCobraCmd(del *DeleteCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.DNSZoneUsage,
		Short:         msg.DNSZoneDeleteShortDescription,
		Long:          msg.DNSZoneDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion delete dns-zone --zone-id 107313
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("zone-id") {
				answer, err := del.AskInput(msg.DNSZoneDeleteAskInputZoneID)
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

			ctx := context.Background()

			err := del.DeleteDNSZone(ctx, zoneID)
			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDelete.Error(), err)
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DNSZoneDeleteOutputSuccess, zoneID),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&deleteOut)
		},
	}

	cobraCmd.Flags().Int64Var(&zoneID, "zone-id", 0, msg.DNSZoneFlagId)
	cobraCmd.Flags().BoolP("help", "h", false, msg.DNSZoneDeleteHelpFlag)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}
