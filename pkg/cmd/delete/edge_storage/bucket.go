package edge_storage

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_storage"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var name string

func NewBucket(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           msg.USAGE,
		Short:         msg.SHORT_DESCRIPTION_DELETE_BUCKET,
		Long:          msg.LONG_DESCRIPTION_DELETE_BUCKET,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(msg.EXAMPLE_DELETE_BUCKET),
		RunE: runE(f),
	}
	cmd.Flags().StringVar(&name, "name", "", msg.FLAG_NAME_BUCKET)
	cmd.Flags().BoolP("help", "h", false, msg.FLAG_HELP_DELETE_BUCKET)	
	return cmd
}

func runE(f *cmdutil.Factory) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("name") {	
			answer, err := utils.AskInput(msg.ASK_NAME_DELETE_BUCKET)
			if err != nil {
				return err
			}
			name = answer 
		}
		client := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
		ctx := context.Background()
		err := client.DeleteBucket(ctx, name)
		if err != nil {
			return fmt.Errorf(msg.ERROR_DELETE_BUCKET, err.Error())
		}
		logger.FInfo(f.IOStreams.Out, fmt.Sprintf(msg.OUTPUT_DELETE_BUCKET, name))
		return nil
	}
}
