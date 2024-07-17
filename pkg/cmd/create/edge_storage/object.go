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
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
)

func commandObjects(fact *factoryObjects) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         msg.SHORT_DESCRIPTION_CREATE_OBJECTS,
		Long:          msg.LONG_DESCRIPTION_CREATE_OBJECTS,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.EXAMPLE_CREATE_OBJECTS),
		RunE:          fact.RunE,
	}
	fact.AddFlags(cmd.Flags())
	return cmd
}

type factoryObjects struct {
	Factory               *cmdutil.Factory
	FlagFileUnmarshalJSON func(path string, request interface{}) error                      // utils.FlagFileUnmarshalJSON
	Join                  func(elem ...string) string                                       // filepath.Join
	MatchFilePath         func(path string, limAndPref ...int) (mimemagic.MediaType, error) // mimemagic.MatchFilePath
	AskInput              func(msg string) (string, error)                                  // utils.AskInput
	Open                  func(name string) (*os.File, error)                               // os.Open
	Getwd                 func() (dir string, err error)                                    // os.Getwd
	fieldsObjects
}

func NewFactoryObjects(fact *cmdutil.Factory) *factoryObjects {
	return &factoryObjects{
		Factory:               fact,
		FlagFileUnmarshalJSON: utils.FlagFileUnmarshalJSON,
		Join:                  filepath.Join,
		MatchFilePath:         mimemagic.MatchFilePath,
		AskInput:              utils.AskInput,
		Open:                  os.Open,
		Getwd:                 os.Getwd,
	}
}

func (fact *factoryObjects) RunE(cmd *cobra.Command, args []string) error {
	f := fact.Factory
	if cmd.Flags().Changed("file") {
		err := fact.FlagFileUnmarshalJSON(fact.fieldsObjects.fileJSON, &fact.fieldsObjects)
		if err != nil {
			return utils.ErrorUnmarshalReader
		}
	} else {
		err := fact.CreateRequestFromFlags(cmd)
		if err != nil {
			return err
		}
	}
	fileOptions, err := fileOptions(fact)
	if err != nil {
		return err
	}
	client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
	err = client.CreateObject(context.Background(), fileOptions, fact.BucketName, fact.ObjectKey)
	if err != nil {
		return fmt.Errorf(msg.ERROR_CREATE_OBJECT, err)
	}
	creatOut := output.GeneralOutput{
		Msg: msg.OUTPUT_CREATE_OBJECT,
		Out: f.IOStreams.Out,
	}
	return output.Print(&creatOut)
}

func (fact *factoryObjects) CreateRequestFromFlags(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("bucket-name") {
		answers, err := fact.AskInput(msg.ASK_NAME_CREATE_BUCKET)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fact.BucketName = answers
	}
	if !cmd.Flags().Changed("object-key") {
		answers, err := fact.AskInput(msg.ASK_OBJECT_KEY_CREATE_OBJECT)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fact.ObjectKey = answers
	}
	if !cmd.Flags().Changed("source") {
		answers, err := fact.AskInput(msg.ASK_SOURCE_CREATE_OBJECT)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fact.Source = answers
	}
	return nil
}

func (fact *factoryObjects) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&fact.BucketName, "bucket-name", "", msg.FLAG_NAME_BUCKET)
	flags.StringVar(&fact.ObjectKey, "object-key", "", msg.FLAG_NAME_CREATE_OBJECT)
	flags.StringVar(&fact.Source, "source", "", msg.FLAG_NAME_CREATE_SOURCE)
	flags.StringVar(&fact.fileJSON, "file", "", msg.FLAG_FILE_JSON_CREATE_OBJECTS)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_CREATE_OBJECTS)
}

func fileOptions(fact *factoryObjects) (*contracts.FileOps, error) {
	pathWorkingDir, err := fact.Getwd()
	if err != nil {
		return &contracts.FileOps{}, utils.ErrorInternalServerError
	}
	pathFull := fact.Join(pathWorkingDir, fact.Source)
	fileContent, err := fact.Open(pathFull)
	if err != nil {
		logger.Debug("Error while trying to read file <"+fact.Source+"> about to be created object of the edge storage", zap.Error(err))
		return &contracts.FileOps{}, err
	}
	mimeType, err := fact.MatchFilePath(fact.Source, -1)
	if err != nil {
		logger.Debug("Error while matching file path", zap.Error(err))
		return &contracts.FileOps{}, err
	}
	return &contracts.FileOps{
		Path:        fact.Source,
		MimeType:    mimeType.MediaType(),
		FileContent: fileContent,
	}, nil
}
