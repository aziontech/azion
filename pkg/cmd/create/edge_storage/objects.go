package edge_storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/zRedShift/mimemagic"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func NewObjects(f *cmdutil.Factory) *cobra.Command {
	fields := &FieldsObjects{
		Factory: f,
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         msg.SHORT_DESCRIPTION_CREATE_OBJECTS,
		Long:          msg.LONG_DESCRIPTION_CREATE_OBJECTS,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.EXAMPLE_CREATE_OBJECTS),
		RunE:          fields.RunE,
	}
	fields.AddFlags(cmd.Flags())
	return cmd
}

func (fields *FieldsObjects) RunE(cmd *cobra.Command, args []string) error {
	f := fields.Factory
	if cmd.Flags().Changed("file") {
		err := utils.FlagFileUnmarshalJSON(fields.FileJSON, &fields)
		if err != nil {
			return utils.ErrorUnmarshalReader
		}
	} else {
		err := fields.CreateRequestFromFlags(cmd)
		if err != nil {
			return err
		}
	}
	fileOptions, err := fileOptions(fields)
	if err != nil {
		return err
	}
	client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
	err = client.CreateObject(context.Background(), fileOptions, fields.BucketName, fields.ObjectKey)
	if err != nil {
		return fmt.Errorf(msg.ERROR_CREATE_OBJECT, err)
	}
	logger.FInfo(f.IOStreams.Out, msg.OUTPUT_CREATE_OBJECT)
	return nil
}

func (fields *FieldsObjects) CreateRequestFromFlags(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("bucket-name") {
		answers, err := utils.AskInput(msg.ASK_NAME_CREATE_BUCKET)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.BucketName = answers
	}
	if !cmd.Flags().Changed("object-key") {
		answers, err := utils.AskInput(msg.ASK_OBJECT_KEY_CREATE_OBJECT)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.ObjectKey = answers
	}
	if !cmd.Flags().Changed("source") {
		answers, err := utils.AskInput(msg.ASK_SOURCE_CREATE_OBJECT)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.Source = answers
	}
	return nil
}

func (fields *FieldsObjects) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&fields.BucketName, "bucket-name", "", msg.FLAG_NAME_CREATE_BUCKET)
	flags.StringVar(&fields.ObjectKey, "object-key", "", msg.FLAG_NAME_CREATE_OBJECT)
	flags.StringVar(&fields.Source, "source", "", msg.FLAG_NAME_CREATE_SOURCE)
	flags.StringVar(&fields.FileJSON, "file", "", msg.FLAG_FILE_JSON_CREATE_OBJECTS)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_CREATE_OBJECTS)
}

func fileOptions(fields *FieldsObjects) (*contracts.FileOps, error) {
	pathWorkingDir, err := os.Getwd()
	if err != nil {
		return &contracts.FileOps{}, utils.ErrorInternalServerError
	}
	pathFull := filepath.Join(pathWorkingDir, fields.Source)
	fileContent, err := os.Open(pathFull)
	if err != nil {
		logger.Debug("Error while trying to read file <"+fields.Source+"> about to be created object of the edge storage", zap.Error(err))
		return &contracts.FileOps{}, err
	}
	mimeType, err := mimemagic.MatchFilePath(fields.Source, -1)
	if err != nil {
		logger.Debug("Error while matching file path", zap.Error(err))
		return &contracts.FileOps{}, err
	}
	return &contracts.FileOps{
		Path:        fields.Source,
		MimeType:    mimeType.MediaType(),
		FileContent: fileContent,
	}, nil
}
