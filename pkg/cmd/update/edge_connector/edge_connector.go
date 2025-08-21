package edgeconnector

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_connector"
	api "github.com/aziontech/azion-cli/pkg/api/edge_connector"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	ID     string
	InPath string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.UpdateShortDescription,
		Long:          msg.UpdateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion update edge-connector --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			// either connector-id or in path should be passed
			if !cmd.Flags().Changed("connector-id") {
				answer, err := utils.AskInput(msg.UpdateAskEdgeConnectorID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				fields.ID = answer
			}

			request := api.UpdateRequest{}

			if cmd.Flags().Changed("file") {
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &request)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request, fields.ID)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateConnector.Error(), err)
			}

			var id int64
			if response.EdgeConnectorHTTP != nil {
				id = response.EdgeConnectorHTTP.GetId()
			} else if response.EdgeConnectorLiveIngest != nil {
				id = response.EdgeConnectorLiveIngest.GetId()
			} else if response.EdgeConnectorStorage != nil {
				id = response.EdgeConnectorStorage.GetId()
			}

			updateOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.UpdateOutputSuccess, id),
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
	flags.StringVar(&fields.ID, "connector-id", "", msg.FlagID)
	flags.StringVar(&fields.InPath, "file", "", msg.UpdateFlagFile)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}
