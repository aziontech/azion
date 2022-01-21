package list

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "list [flags]",
		Short:         "List the Edge Functions of your account",
		Long:          "List the Edge Functions of your account",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions list [--details]
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("IMPLEMENT ME")
		},
	}

	return cmd
}
