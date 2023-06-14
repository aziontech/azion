package delete

import (
	"context"
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/variables"
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)


func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var variables_id string

	deleteCmd := &cobra.Command{
		Use: msg.VariableDeleteUsage,
		Short: msg.VariableDeleteShortDescription,
		Long: msg.CacheSettingsDeleteLongDescription,
		SilenceErrors: true,
		SilenceUsage: true,
		Example: heredoc.Doc(`
		$ azioncli variables delete --variable-id 1234
		$ azioncli variables delete -v 1234
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("variable-id") {
				return msg.ErrorMissingVariableIdArgumentDelete
			}

			client :=  api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			
			ctx := context.Background()

			err := client.Delete(ctx, variables_id)

			if err != nil {
				return fmt.Errorf(msg.ErrorFailToDeleteVariable.Error(), err)
			}
			out := f.IOStreams.Out
			fmt.Fprintf(out, msg.VariableDeleteOutputSuccess, variables_id)

			return nil
		},
	}

	deleteCmd.Flags().StringVarP(&variables_id, "variable-id", "v", "", msg.VariableFlagId)
	deleteCmd.Flags().BoolP("help", "h", false, msg.ValiableDeleteHelpFlag)

	return deleteCmd
}