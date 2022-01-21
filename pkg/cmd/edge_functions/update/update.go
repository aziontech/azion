package update

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "update <edge_function_id> [flags]",
		Short:         "Update an Edge Function",
		Long:          "Update an Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions update 4185 –code ./mycode/function.js –args ./mycode/myargs.json
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("IMPLEMENT ME")
		},
	}

	return cmd
}
