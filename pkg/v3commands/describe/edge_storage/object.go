package edge_storage

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	"github.com/aziontech/azion-cli/utils"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	api "github.com/aziontech/azion-cli/pkg/v3api/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewObject(f *cmdutil.Factory) *cobra.Command {
	fields := &Fields{
		Factory: f,
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         "",
		Long:          "",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc("$ azion describe edge-storage object --bucket-id 1234 --object-key 'object-key'"),
		RunE:          fields.RunE,
	}
	fields.AddFlags(cmd.Flags())
	return cmd
}

func (f *Fields) RunE(cmd *cobra.Command, args []string) error {
	if !cmd.Flags().Changed("bucket-name") {
		answers, err := utils.AskInput(msg.ASK_NAME_CREATE_BUCKET)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		f.BucketName = answers
	}
	if !cmd.Flags().Changed("object-key") {
		answers, err := utils.AskInput(msg.ASK_OBJECT_KEY)
		if err != nil {
			logger.Debug("Error while parsing answer", zap.Error(err))
			return utils.ErrorParseResponse
		}
		f.ObjectKey = answers
	}

	client := api.NewClient(f.Factory.HttpClient, f.Factory.Config.GetString("storage_url"), f.Factory.Config.GetString("token"))
	ctx := context.Background()
	bFile, err := client.GetObject(ctx, f.BucketName, f.ObjectKey)
	if err != nil {
		return fmt.Errorf(msg.ERROR_DESCRIBE_OBJECT, err)
	}

	describeOut := output.GeneralOutput{
		Msg:   string(bFile),
		Out:   f.Factory.IOStreams.Out,
		Flags: f.Factory.Flags,
	}
	return output.Print(&describeOut)
}

func (f *Fields) AddFlags(flags *pflag.FlagSet) {
	flags.StringVar(&f.BucketName, "bucket-name", "", "")
	flags.StringVar(&f.ObjectKey, "object-key", "", "")
	flags.BoolP("help", "h", false, "")
}
