package create

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "create [flags]",
		Short:         "Create a new Edge Function",
		Long:          "Create a new Edge Function",
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_functions create -–name myfunc -–language javascript –-code ./mycode/function.js  -–state active --initiator-type edge-application
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("IMPLEMENT ME")
		},
	}

	return cmd
}
