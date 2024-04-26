package edge_storage

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	sdk "github.com/aziontech/azionapi-go-sdk/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
)

func NewBucket(f *cmdutil.Factory) *cobra.Command {
	bucket := &bucket{
		factory: f,
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_BUCKET,
		Short:         msg.SHORT_DESCRIPTION_CREATE_BUCKET,
		Long:          msg.LONG_DESCRIPTION_CREATE_BUCKET,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.EXAMPLE_UPDATE_BUCKET),
		RunE:          bucket.runE,
	}
	bucket.addFlags(cmd.Flags())
	return cmd
}

func (b *bucket) runE(cmd *cobra.Command, args []string) error {
	request := api.RequestBucket{}
	if cmd.Flags().Changed("file") {
		err := utils.FlagFileUnmarshalJSON(b.fileJSON, &request)
		if err != nil {
			return utils.ErrorUnmarshalReader
		}
	} else {
		err := b.createRequestFromFlags(cmd, &request)
		if err != nil {
			return err
		}
	}
	client := api.NewClient(
		b.factory.HttpClient,
		b.factory.Config.GetString("storage_url"),
		b.factory.Config.GetString("token"))
	err := client.UpdateBucket(context.Background(), request.GetName(), request.GetEdgeAccess())
	if err != nil {
		return fmt.Errorf(msg.ERROR_UPDATE_BUCKET, err)
	}

	updateOut := output.GeneralOutput{
		Msg: msg.OUTPUT_UPDATE_BUCKET,
		Out: b.factory.IOStreams.Out,
	}
	return output.Print(&updateOut)
}

func (b *bucket) createRequestFromFlags(cmd *cobra.Command, request *api.RequestBucket) error {
	if !cmd.Flags().Changed("name") {
		answers, err := utils.AskInput(msg.ASK_NAME_UPDATE_BUCKET)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		b.name = answers
	}
	if !cmd.Flags().Changed("edge-access") {
		answers, err := utils.Select(
			msg.ASK_EDGE_ACCESSS_CREATE_BUCKET,
			[]string{string(sdk.READ_ONLY), string(sdk.READ_WRITE), string(sdk.RESTRICTED)})
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		b.edgeAccess = answers
	}
	request.SetName(b.name)
	request.SetEdgeAccess(sdk.EdgeAccessEnum(b.edgeAccess))
	return nil
}

func (b *bucket) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&b.name, "name", "", msg.FLAG_NAME_BUCKET)
	flags.StringVar(&b.edgeAccess, "edge-access", "", msg.FLAG_EDGE_ACCESS_CREATE_BUCKET)
	flags.StringVar(&b.fileJSON, "file", "", msg.FLAG_FILE_JSON_CREATE_BUCKET)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_CREATE_BUCKET)
}
