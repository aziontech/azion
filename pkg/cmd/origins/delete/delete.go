package delete

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/origins"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var applicationID int64
	var originKey string
	cmd := &cobra.Command{
		Use:           origins.OriginsDeleteUsage,
		Short:         origins.OriginsDeleteShortDescription,
		Long:          origins.OriginsDeleteLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		  $ azion origins delete --application-id 1673635839 --origin-key 03a6e7bf-8e26-49c7-a66e-ab8eaa425086
		  $ azion origins delete -a 1673635839 -o 03a6e7bf-8e26-49c7-a66e-ab8eaa425086
    `),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("application-id") || !cmd.Flags().Changed("origin-key") {
				return origins.ErrorMissingArgumentsDelete
			}
			if err := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token")).
				DeleteOrigins(context.Background(), applicationID, originKey); err != nil {
				return fmt.Errorf(origins.ErrorFailToDelete.Error(), err)
			}
			fmt.Fprintf(f.IOStreams.Out, origins.OriginsDeleteOutputSuccess, originKey)
			return nil
		},
	}

	cmd.Flags().Int64VarP(&applicationID, "application-id", "a", 0, origins.OriginsDeleteFlagApplicationID)
	cmd.Flags().StringVarP(&originKey, "origin-key", "o", "", origins.OriginsDeleteFlagOriginKey)
	cmd.Flags().BoolP("help", "h", false, origins.OriginsDeleteHelpFlag)
	return cmd
}
