package edge_storage

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"go.uber.org/zap"

	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

type Fields struct {
	Name       string
	EdgeAccess string
	FileJSON   string
}

func NewBucket(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{}

	cmd := &cobra.Command{
		Use:           msg.USAGE_BUCKET,
		Short:         msg.SHORT_DESCRIPTION_CREATE_BUCKET,
		Long:          msg.LONG_DESCRIPTION_CREATE_BUCKET,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.EXAMPLE_UPDATE_BUCKET),
		RunE:          runE(f, fields),
	}

	flags := cmd.Flags()
	addFlags(flags, fields)
	return cmd
}

func runE(f *cmdutil.Factory, fields *Fields) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		request := api.RequestBucket{}
		if cmd.Flags().Changed("file") {
			err := utils.FlagFileUnmarshalJSON(fields.FileJSON, &request)
			if err != nil {
				return utils.ErrorUnmarshalReader
			}
		} else {
			err := createRequestFromFlags(cmd, fields, &request)
			if err != nil {
				return err
			}
		}

		client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
		err := client.UpdateBucket(context.Background(), request.GetName())
		if err != nil {
			return fmt.Errorf(msg.ERROR_UPDATE_BUCKET, err)
		}

		logger.FInfo(f.IOStreams.Out, msg.OUTPUT_UPDATE_BUCKET)
		return nil
	}
}

func createRequestFromFlags(cmd *cobra.Command, fields *Fields, request *api.RequestBucket) error {
	if !cmd.Flags().Changed("name") {
		answers, err := utils.AskInput(msg.ASK_NAME_UPDATE_BUCKET)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		fields.Name = answers
	}

	request.SetName(fields.Name)
	return nil
}

func addFlags(flags *pflag.FlagSet, fields *Fields) {
	flags.StringVar(&fields.Name, "name", "", msg.FLAG_NAME_BUCKET)
	flags.StringVar(&fields.FileJSON, "file", "", msg.FLAG_FILE_JSON_CREATE_BUCKET)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_CREATE_BUCKET)
}
