package edge_storage

import (
	"context"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zRedShift/mimemagic"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
)

func NewObject(f *cmdutil.Factory) *cobra.Command {
	object := &object{
		factory: f,
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         msg.SHORT_DESCRIPTION_CREATE_BUCKET,
		Long:          msg.LONG_DESCRIPTION_CREATE_BUCKET,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.EXAMPLE_UPDATE_OBJECT),
		RunE:          object.runE,
	}
	object.addFlags(cmd.Flags())
	return cmd
}

func (o *object) runE(cmd *cobra.Command, args []string) error {
	if cmd.Flags().Changed("file") {
		err := utils.FlagFileUnmarshalJSON(o.fileJSON, o)
		if err != nil {
			return utils.ErrorUnmarshalReader
		}
	} else {
		err := o.createRequestFromFlags(cmd)
		if err != nil {
			return err
		}
	}
	file, err := os.Open(o.Source)
	if err != nil {
		logger.Debug("Error while trying to read file <"+o.Source+"> about to be uploaded", zap.Error(err))
		return err
	}
	mimeType, err := mimemagic.MatchFilePath(o.Source, -1)
	if err != nil {
		logger.Debug("Error while matching file path", zap.Error(err))
		return err
	}
	client := api.NewClient(
		o.factory.HttpClient, o.factory.Config.GetString("storage_url"), o.factory.Config.GetString("token"))
	err = client.UpdateObject(
		context.Background(), o.BucketName, o.ObjectKey, mimeType.MediaType(), file)
	if err != nil {
		return fmt.Errorf(msg.ERROR_UPDATE_BUCKET, err)
	}

	updateOut := output.GeneralOutput{
		Msg: msg.OUTPUT_UPDATE_OBJECT,
		Out: o.factory.IOStreams.Out,
	}
	return output.Print(&updateOut)
}

func (o *object) createRequestFromFlags(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("bucket-name") {
		answers, err := utils.AskInput(msg.ASK_NAME_CREATE_BUCKET)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		o.BucketName = answers
	}
	if !cmd.Flags().Changed("object-key") {
		answers, err := utils.AskInput(msg.ASK_NAME_UPDATE_OBJECT)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		o.ObjectKey = answers
	}
	if !cmd.Flags().Changed("source") {
		answers, err := utils.AskInput(msg.ASK_SOURCE_CREATE_OBJECT)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		o.Source = answers
	}
	return nil
}

func (o *object) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&o.BucketName, "bucket-name", "", msg.FLAG_NAME_BUCKET)
	flags.StringVar(&o.ObjectKey, "object-key", "", msg.FLAG_NAME_OBJECT)
	flags.StringVar(&o.Source, "source", "", msg.FLAG_SOURCE_UPDATE_OBJECT)
	flags.StringVar(&o.fileJSON, "file", "", msg.FLAG_FILE_JSON_CREATE_BUCKET)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_CREATE_BUCKET)
}
