package connector

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/connector"
	api "github.com/aziontech/azion-cli/pkg/api/connector"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	ID     string
	InPath string
	Type   string
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
		$ azion update connector --in "update.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			// either connector-id or in path should be passed
			if !cmd.Flags().Changed("connector-id") {
				answer, err := utils.AskInput(msg.UpdateAskConnectorID)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				fields.ID = answer
			}

			request := api.UpdateRequest{}

			if !cmd.Flags().Changed("type") {
				answer, err := utils.AskInput(msg.UpdateAskConnectorType)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				fields.Type = answer
			}

			if !cmd.Flags().Changed("file") {
				answer, err := utils.AskInput(msg.UpdateAskConnectorFile)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				fields.InPath = answer
			}

			switch fields.Type {
			case "http":
				httpStruct := sdk.PatchedConnectorHTTPRequest{}
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &httpStruct)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
				request.PatchedConnectorHTTPRequest = &httpStruct
			case "storage":
				storageStruct := sdk.PatchedConnectorStorageRequest{}
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &storageStruct)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
				request.PatchedConnectorStorageRequest = &storageStruct
			case "live_ingest":
				liveIngestStruct := sdk.PatchedConnectorLiveIngestRequest{}
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &liveIngestStruct)
				if err != nil {
					return utils.ErrorUnmarshalReader
				}
				request.PatchedConnectorLiveIngestRequest = &liveIngestStruct
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Update(ctx, &request, fields.ID)

			if err != nil {
				return fmt.Errorf(msg.ErrorUpdateConnector.Error(), err)
			}

			var id int64
			if response.ConnectorHTTP != nil {
				id = response.ConnectorHTTP.GetId()
			} else if response.ConnectorLiveIngest != nil {
				id = response.ConnectorLiveIngest.GetId()
			} else if response.ConnectorStorage != nil {
				id = response.ConnectorStorage.GetId()
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
	flags.StringVar(&fields.Type, "type", "", msg.FlagType)
	flags.BoolP("help", "h", false, msg.UpdateHelpFlag)
}
