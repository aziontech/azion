package object

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type DeleteObjectCmd struct {
	Io           *iostreams.IOStreams
	ReadInput    func(string) (string, error)
	DeleteObject func(context.Context, string, string) error
	AskInput     func(string) (string, error)
	PrintOutput  func(*output.GeneralOutput) error
	Bucket       string
	ObjectKey    string
}

func NewDeleteObjectCmd(f *cmdutil.Factory) *DeleteObjectCmd {
	return &DeleteObjectCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteObject: func(ctx context.Context, bucketName, objectKey string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
			return client.DeleteObject(ctx, bucketName, objectKey)
		},
		AskInput: utils.AskInput,
		PrintOutput: func(out *output.GeneralOutput) error {
			return output.Print(out)
		},
	}
}

func NewObjectCmd(delete *DeleteObjectCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.USAGE_OBJECTS,
		Short:         msg.SHORT_DESCRIPTION_DELETE_OBJECTS,
		Long:          msg.LONG_DESCRIPTION_DELETE_OBJECTS,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc("$ azion delete edge-storage object --bucket-name 'bucket-name' --object-key 'object-key'"),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("bucket-name") {
				answer, err := delete.AskInput(msg.ASK_NAME_CREATE_BUCKET)
				if err != nil {
					return err
				}
				delete.Bucket = answer
			}

			if !cmd.Flags().Changed("object-key") {
				answer, err := delete.AskInput(msg.ASK_OBJECT_DELETE_OBJECT)
				if err != nil {
					return err
				}
				delete.ObjectKey = answer
			}

			ctx := context.Background()

			err = delete.DeleteObject(ctx, delete.Bucket, delete.ObjectKey)
			if err != nil {
				return fmt.Errorf(msg.ERROR_DELETE_OBJECT, err.Error())
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OUTPUT_DELETE_OBJECT, delete.ObjectKey),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return delete.PrintOutput(&deleteOut)
		},
	}

	cobraCmd.Flags().StringVar(&delete.Bucket, "bucket-name", "", msg.FLAG_NAME_BUCKET)
	cobraCmd.Flags().StringVar(&delete.ObjectKey, "object-key", "", msg.FLAG_OBJECT_KEY_OBJECT)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP_DELETE_BUCKET)

	return cobraCmd
}

func NewObject(f *cmdutil.Factory) *cobra.Command {
	return NewObjectCmd(NewDeleteObjectCmd(f), f)
}
