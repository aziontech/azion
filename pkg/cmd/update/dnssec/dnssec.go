package dnssec

import (
	"context"
	"fmt"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/dnssec"
	api "github.com/aziontech/azion-cli/pkg/api/dnssec"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Fields struct {
	ZoneID  int64
	Enabled bool
	Path    string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.DNSSECUsage,
		Short:         msg.DNSSECUpdateShortDescription,
		Long:          msg.DNSSECUpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion update dnssec --zone-id 12312 --enabled true
        $ azion update dnssec --zone-id 12312 --enabled false
        $ azion update dnssec --zone-id 12312 --file "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			request := sdk.PatchedDNSSECRequest{}

			if !cmd.Flags().Changed("zone-id") {
				answers, err := utils.AskInput(msg.DNSSECUpdateAskInputZoneID)
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
				if !cmd.Flags().Changed("enabled") {
					answer, err := utils.AskInput(msg.DNSSECUpdateAskInputEnabled)
					if err != nil {
						logger.Debug("Error while parsing answer", zap.Error(err))
						return utils.ErrorParseResponse
					}

					enabled, err := strconv.ParseBool(answer)
					if err != nil {
						logger.Debug("Error while parsing string to boolean", zap.Error(err))
						return msg.ErrorConvertEnabled
					}

					fields.Enabled = enabled
				}

				request.SetEnabled(fields.Enabled)
			}

			_, err := client.Update(context.Background(), request, fields.ZoneID)
			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateDNSSEC.Error(), err)
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.DNSSECUpdateOutputSuccess, fields.ZoneID),
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
	flags.Int64Var(&fields.ZoneID, "zone-id", 0, msg.DNSSECFlagZoneID)
	flags.BoolVar(&fields.Enabled, "enabled", false, msg.DNSSECUpdateFlagEnabled)
	flags.StringVar(&fields.Path, "file", "", msg.DNSSECUpdateFlagIn)
	flags.BoolP("help", "h", false, msg.DNSSECUpdateHelpFlag)
}
