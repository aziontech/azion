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
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

type Fields struct {
	Name   string
	Type   string
	Path   string
	InPath string
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.CreateShortDescription,
		Long:          msg.CreateLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion create connector --file "create.json"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := api.NewCreateRequest()

			if !cmd.Flags().Changed("type") {
				answer, err := utils.AskInput("Enter the type of your Connector")
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}
				fields.Type = answer
			}

			if !cmd.Flags().Changed("file") {
				answer, err := utils.AskInput("Enter the path of the json to create the Connector:")
				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}
				fields.InPath = answer
			}

			switch fields.Type {
			case "http":
				httpStruct := sdk.ConnectorHTTPRequest{}
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &httpStruct)
				if err != nil {
					logger.Debug("Failed to unmarshal file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
				request.ConnectorHTTPRequest = &httpStruct
			case "storage":
				storageStruct := sdk.ConnectorStorageRequest{}
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &storageStruct)
				if err != nil {
					logger.Debug("Failed to unmarshal file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
				request.ConnectorStorageRequest = &storageStruct
			case "live_ingest":
				liveIngestStruct := sdk.ConnectorLiveIngestRequest{}
				err := utils.FlagFileUnmarshalJSON(fields.InPath, &liveIngestStruct)
				if err != nil {
					logger.Debug("Failed to unmarshal file", zap.Error(err))
					return utils.ErrorUnmarshalReader
				}
				request.ConnectorLiveIngestRequest = &liveIngestStruct
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))

			ctx := context.Background()
			response, err := client.Create(ctx, request)
			if err != nil {
				return fmt.Errorf(msg.ErrorCreateConnector.Error(), err)
			}

			var id int64
			switch fields.Type {
			case "http":
				id = response.ConnectorHTTP.GetId()
			case "storage":
				id = response.ConnectorStorage.GetId()
			case "live_ingest":
				id = response.ConnectorLiveIngest.GetId()
			}

			creatOut := output.GeneralOutput{
				Msg: fmt.Sprintf(msg.CreateOutputSuccess, id),
				Out: f.IOStreams.Out,
			}
			return output.Print(&creatOut)
		},
	}

	flags := cmd.Flags()
	addFlags(flags, fields)

	return cmd
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.InPath, "file", "", msg.FlagIn)
	flags.StringVar(&fields.Type, "type", "", msg.FlagType)
	flags.BoolP("help", "h", false, msg.CreateFlagHelp)
}
