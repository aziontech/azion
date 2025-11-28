package bucket

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/schedule"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type DeleteBucketCmd struct {
	Io            *iostreams.IOStreams
	ReadInput     func(string) (string, error)
	DeleteBucket  func(context.Context, string) error
	AskInput      func(string) (string, error)
	DeleteAll     func(*api.Client, context.Context, string, string) error
	ConfirmDelete func(bool, string, bool) bool
	PrintOutput   func(*output.GeneralOutput) error
}

var (
	bucketName  string
	forceDelete bool
)

func NewDeleteBucketCmd(f *cmdutil.Factory) *DeleteBucketCmd {
	return &DeleteBucketCmd{
		Io: f.IOStreams,
		ReadInput: func(prompt string) (string, error) {
			return utils.AskInput(prompt)
		},
		DeleteBucket: func(ctx context.Context, bucketName string) error {
			client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
			return client.DeleteBucket(ctx, bucketName)
		},
		AskInput:      utils.AskInput,
		DeleteAll:     deleteAllObjects,
		ConfirmDelete: utils.Confirm,
		PrintOutput: func(out *output.GeneralOutput) error {
			return output.Print(out)
		},
	}
}

func NewBucketCmd(delete *DeleteBucketCmd, f *cmdutil.Factory) *cobra.Command {
	cobraCmd := &cobra.Command{
		Use:           msg.USAGE_BUCKET,
		Short:         msg.SHORT_DESCRIPTION_DELETE_BUCKET,
		Long:          msg.LONG_DESCRIPTION_DELETE_BUCKET,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc("$ azion delete edge-storage bucket --bucket-id 1234"),
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if !cmd.Flags().Changed("name") {
				answer, err := delete.AskInput(msg.ASK_NAME_DELETE_BUCKET)
				if err != nil {
					return err
				}
				bucketName = answer
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err = client.DeleteBucket(ctx, bucketName)
			if err != nil {
				if strings.Contains(err.Error(), msg.ERROR_NO_EMPTY_BUCKET) {
					if !forceDelete {
						if !delete.ConfirmDelete(f.GlobalFlagAll, msg.ASK_NOT_EMPTY_BUCKET, false) {
							return nil
						}
					}
					logger.FInfo(f.IOStreams.Out, "Delete all objects from bucket\n")
					logger.FInfo(f.IOStreams.Out, "Deleting objects...\n")
					if err := delete.DeleteAll(client, ctx, bucketName, ""); err != nil {
						return err
					}
					err := client.DeleteBucket(ctx, bucketName)
					if err != nil {
						if strings.Contains(err.Error(), msg.ERROR_NO_EMPTY_BUCKET) {
							logger.FInfo(f.IOStreams.Out, "Bucket deletion was scheduled successfully\n")
							return schedule.NewSchedule(nil, f, bucketName, schedule.DELETE_BUCKET)
						} else {
							return fmt.Errorf(msg.ERROR_DELETE_BUCKET, err.Error())
						}
					}
					return nil
				}
				return fmt.Errorf(msg.ERROR_DELETE_BUCKET, err.Error())
			}

			deleteOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.OUTPUT_DELETE_BUCKET, bucketName),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return delete.PrintOutput(&deleteOut)
		},
	}

	cobraCmd.Flags().StringVar(&bucketName, "name", "", msg.FLAG_NAME_BUCKET)
	cobraCmd.Flags().BoolVar(&forceDelete, "force", false, msg.FLAG_FORCE)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP_DELETE_BUCKET)

	return cobraCmd
}

func NewBucket(f *cmdutil.Factory) *cobra.Command {
	return NewBucketCmd(NewDeleteBucketCmd(f), f)
}

func deleteAllObjects(client *api.Client, ctx context.Context, name, continuationToken string) error {
	objects, err := client.ListObject(ctx, name, &contracts.ListOptions{ContinuationToken: continuationToken})
	if err != nil {
		return err
	}
	for _, object := range objects.Results {
		err := client.DeleteObject(ctx, name, object.GetKey())
		if err != nil {
			return err
		}
	}
	return nil
}
