package describe

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "describe <edge_function_id> [flags]",
		Short:         "Describe a given Edge Function",
		Long:          "Describeb a given Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions describe 1337 [--with-code]
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("IMPLEMENT ME")
		},
	}

	return cmd
}
