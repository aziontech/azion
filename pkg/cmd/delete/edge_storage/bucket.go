package edge_storage

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/schedule"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func NewBucket(f *cmdutil.Factory) *cobra.Command {
	bucket := bucket{
		factory: f,
	}
	cmd := &cobra.Command{
		Use:           msg.USAGE_BUCKET,
		Short:         msg.SHORT_DESCRIPTION_DELETE_BUCKET,
		Long:          msg.LONG_DESCRIPTION_DELETE_BUCKET,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.EXAMPLE_DELETE_BUCKET),
		RunE:          bucket.runE,
	}
	bucket.addFlags(cmd.Flags())
	return cmd
}

func (b *bucket) runE(cmd *cobra.Command, _ []string) error {
	if !cmd.Flags().Changed("name") {
		answer, err := utils.AskInput(msg.ASK_NAME_DELETE_BUCKET)
		if err != nil {
			return err
		}
		b.name = answer
	}
	client := api.NewClient(
		b.factory.HttpClient,
		b.factory.Config.GetString("storage_url"),
		b.factory.Config.GetString("token"))
	ctx := context.Background()
	err := client.DeleteBucket(ctx, b.name)
	if err != nil {
		if msg.ERROR_NO_EMPTY_BUCKET == err.Error() {
			if !b.force {
				if !utils.Confirm(b.factory.GlobalFlagAll, msg.ASK_NOT_EMPTY_BUCKET, false) {
					return nil
				}
			}
			logger.FInfo(b.factory.IOStreams.Out, "Delete all objects bucket\n")
			if err := deleteAllObjects(client, ctx, b.name, ""); err != nil {
				return err
			}
			err := client.DeleteBucket(ctx, b.name)
			if err != nil {
				if msg.ERROR_NO_EMPTY_BUCKET == err.Error() {
					logger.FInfo(b.factory.IOStreams.Out, fmt.Sprintf("schedules a delete for the bucket: %s", b.name))
					return schedule.NewSchedule(b.name, schedule.DELETE_BUCKET)
				} else {
					return fmt.Errorf(msg.ERROR_DELETE_BUCKET, err.Error())
				}
			}
			return nil
		}
		return fmt.Errorf(msg.ERROR_DELETE_BUCKET, err.Error())
	}
	deleteOut := output.GeneralOutput{
		Msg:         fmt.Sprintf(msg.OUTPUT_DELETE_BUCKET, b.name),
		Out:         b.factory.IOStreams.Out,
		FlagOutPath: b.factory.Out,
		FlagFormat:  b.factory.Format,
	}
	return output.Print(&deleteOut)
}

func (f *bucket) addFlags(flags *pflag.FlagSet) {
	flags.StringVar(&f.name, "name", "", msg.FLAG_NAME_BUCKET)
	flags.BoolVar(&f.force, "force", false, msg.FLAG_FORCE)
	flags.BoolP("help", "h", false, msg.FLAG_HELP_DELETE_BUCKET)
}

func deleteAllObjects(client *api.Client, ctx context.Context, name, continuationToken string) error {
	objects, err := client.ListObject(ctx, name, &contracts.ListOptions{ContinuationToken: continuationToken})
	if err != nil {
		return err
	}
	if len(objects.GetResults()) > 0 {
		for _, o := range objects.GetResults() {
			err := client.DeleteObject(ctx, name, o.GetKey())
			if err != nil {
				return err
			}
		}
		return deleteAllObjects(client, ctx, name, continuationToken)
	}
	return nil
}
