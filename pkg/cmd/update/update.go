package update

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update"
	cacheSetting "github.com/aziontech/azion-cli/pkg/cmd/update/cache_setting"
	domain "github.com/aziontech/azion-cli/pkg/cmd/update/domain"
	edgeApplication "github.com/aziontech/azion-cli/pkg/cmd/update/edge_application"
	edgeFunction "github.com/aziontech/azion-cli/pkg/cmd/update/edge_function"
	edgeStorage "github.com/aziontech/azion-cli/pkg/cmd/update/edge_storage"
	origin "github.com/aziontech/azion-cli/pkg/cmd/update/origin"
	rulesEngine "github.com/aziontech/azion-cli/pkg/cmd/update/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmd/update/variables"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion update --help
		$ azion update edge-application
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplication.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(domain.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(edgeFunction.NewCmd(f))
	cmd.AddCommand(cacheSetting.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(edgeStorage.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
