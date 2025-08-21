package logs

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/pkg/api/graphql"
	"github.com/aziontech/azion-cli/pkg/cmd/logs/cells"
	"github.com/aziontech/azion-cli/pkg/cmd/logs/http"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion logs cells
		$ azion logs http
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(cells.NewCmd(f))
	cmd.AddCommand(http.NewCmd(f))
	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
