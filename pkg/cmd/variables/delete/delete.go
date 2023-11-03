package delete

import (
	"context"
	"fmt"
	"github.com/aziontech/azion-cli/pkg/messages/variables"

	"github.com/MakeNowJust/heredoc"
	api "github.com/aziontech/azion-cli/pkg/api/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var variable_id string

	deleteCmd := &cobra.Command{
		Use:           variables.DeleteUsage,
		Short:         variables.DeleteShortDescription,
		Long:          variables.DeleteLongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
		$ azion variables delete --variable-id 7a187044-4a00-4a4a-93ed-d230900421f3
		$ azion variables delete -v 7a187044-4a00-4a4a-93ed-d230900421f3
		`),

		RunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Changed("variable-id") {
				return variables.ErrorMissingVariableIdArgumentDelete
			}

			client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

			ctx := context.Background()

			err := client.Delete(ctx, variable_id)

			if err != nil {
				return fmt.Errorf(variables.ErrorFailToDeleteVariable.Error(), err)
			}
			out := f.IOStreams.Out
			fmt.Fprintf(out, variables.DeleteOutputSuccess, variable_id)

			return nil
		},
	}

	deleteCmd.Flags().StringVarP(&variable_id, "variable-id", "v", "", variables.FlagVariableID)
	deleteCmd.Flags().BoolP("help", "h", false, variables.DeleteHelpFlag)

	return deleteCmd
}
