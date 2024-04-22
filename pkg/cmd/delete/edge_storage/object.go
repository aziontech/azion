package edge_storage

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewObject(f *cmdutil.Factory) *cobra.Command {
	object := object{
		factory: f,
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         msg.SHORT_DESCRIPTION_DELETE_OBJECTS,
		Long:          msg.LONG_DESCRIPTION_DELETE_OBJECTS,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.EXAMPLE_DELETE_OBJECTS),
		RunE:          object.runE,
	}
	object.addFlags(cmd.Flags())
	return cmd
}

func (b *object) runE(cmd *cobra.Command, _ []string) error {
	if !cmd.Flags().Changed("bucket-name") {
		answer, err := utils.AskInput(msg.ASK_NAME_CREATE_BUCKET)
		if err != nil {
			return err
		}
		b.bucketName = answer
	}
	if !cmd.Flags().Changed("object-key") {
		answer, err := utils.AskInput(msg.ASK_OBJECT_DELETE_OBJECT)
		if err != nil {
			return err
		}
		b.objectKey = answer
	}
	client := api.NewClient(
		b.factory.HttpClient,
		b.factory.Config.GetString("storage_url"),
		b.factory.Config.GetString("token"))
	ctx := context.Background()
	err := client.DeleteObject(ctx, b.bucketName, b.objectKey)
	if err != nil {
		return fmt.Errorf(msg.ERROR_DELETE_OBJECT, err.Error())
	}

	deleteOut := output.GeneralOutput{
		Msg: fmt.Sprintf(msg.ERROR_DELETE_OBJECT, b.objectKey),
		Out: b.factory.IOStreams.Out}
	return output.Print(&deleteOut)

}

func (f *object) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&f.bucketName, "bucket-name", "", msg.FLAG_NAME_BUCKET)
	flags.StringVar(&f.objectKey, "object-key", "", msg.FLAG_OBJECT_KEY_OBJECT)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_DELETE_BUCKET)
}
