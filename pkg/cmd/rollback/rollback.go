package rollback

import (
	"context"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/rollback"
	apiConnector "github.com/aziontech/azion-cli/pkg/api/connector"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/edge-api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	connectorID int64
	projectPath string
)

type RollbackCmd struct {
	AskInput              func(string) (string, error)
	GetAzionJsonContent   func(pathConf string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
}

func NewDeleteCmd(f *cmdutil.Factory) *RollbackCmd {
	return &RollbackCmd{
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		AskInput:              utils.AskInput,
	}
}

func NewCobraCmd(rollback *RollbackCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.USAGE,
		Short:         msg.SHORTDESCRIPTION,
		Long:          msg.LONGDESCRIPTION,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion rollback --connector-id aaaa-bbbb-cccc-dddd
		`),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("connector-id") {
				answer, err := rollback.AskInput(msg.ASKCONNECTOR)
				if err != nil {
					return err
				}

				num, err := strconv.ParseInt(answer, 10, 64)
				if err != nil {
					logger.Debug("Error while converting answer to int64", zap.Error(err))
					return msg.ERRORCONVERTCONNECTORID
				}

				connectorID = num
			}

			conf, err := rollback.GetAzionJsonContent(projectPath)
			if err != nil {
				logger.Debug("Error while reading azion.json file", zap.Error(err))
				return msg.ERRORAZION
			}

			if conf.Bucket == "" || conf.Prefix == "" {
				return msg.ERRORNEEDSDEPLOY
			}

			timestamp, err := checkForNewTimestamp(f, conf.Prefix, conf.Bucket)
			if err != nil {
				return msg.ERRORROLLBACK
			}

			if timestamp == "" {
				logger.Debug("No previous timestamp found for rollback")
				return msg.ERRORNOPREVIOUS
			}

			logger.Debug("Rolling back to previous timestamp", zap.String("from", conf.Prefix), zap.String("to", timestamp))

			clientConnector := apiConnector.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
			request := apiConnector.UpdateRequest{}

			attributes := sdk.ConnectorStorageAttributesRequest{}
			attributes.SetBucket(conf.Bucket)
			attributes.SetPrefix(timestamp)

			storageRequest := sdk.PatchedConnectorStorageRequest{}
			storageRequest.SetAttributes(attributes)
			request.PatchedConnectorStorageRequest = &storageRequest

			_, err = clientConnector.Update(context.Background(), &request, connectorID)
			if err != nil {
				return msg.ERRORROLLBACK
			}

			conf.Prefix = timestamp
			err = rollback.WriteAzionJsonContent(conf, projectPath)
			if err != nil {
				return msg.ERRORROLLBACK
			}

			rollbackOut := output.GeneralOutput{
				Msg:   msg.SUCCESS,
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&rollbackOut)
		},
	}

	cobraCmd.Flags().Int64Var(&connectorID, "connector-id", 0, msg.FLAGCONNECTORID)
	cobraCmd.Flags().StringVar(&projectPath, "config-dir", "azion", msg.CONFFLAG)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FLAGHELP)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeleteCmd(f), f)
}

func checkForNewTimestamp(f *cmdutil.Factory, referenceTimestamp, bucketName string) (string, error) {
	logger.Debug("Checking if there are previous static files for the following bucket", zap.Any("Bucket name", bucketName))
	client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
	c := context.Background()

	var prevTimestamp string
	var continuationToken string

	for {
		options := &contracts.ListOptions{
			ContinuationToken: continuationToken,
		}

		resp, err := client.ListObject(c, bucketName, options)
		if err != nil {
			return "", err
		}

		for _, object := range resp.Results {
			parts := strings.Split(object.Key, "/")
			if len(parts) > 1 {
				timestamp := parts[0]
				logger.Debug("Found timestamp in bucket", zap.String("timestamp", timestamp), zap.String("key", object.Key))
				if timestamp == referenceTimestamp {
					logger.Debug("Found current timestamp, returning previous", zap.String("current", referenceTimestamp), zap.String("previous", prevTimestamp))
					return prevTimestamp, nil
				}
				prevTimestamp = timestamp
			}
		}

		logger.Debug("continuing to next page", zap.Any("continuation-token", resp.GetContinuationToken()))
		if contToken, ok := resp.GetContinuationTokenOk(); contToken == nil || !ok {
			break
		}
		continuationToken = resp.GetContinuationToken()
	}

	return referenceTimestamp, nil
}
