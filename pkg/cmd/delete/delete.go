package delete

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete"
	cache "github.com/aziontech/azion-cli/pkg/cmd/delete/cache_setting"
	edgeApplication "github.com/aziontech/azion-cli/pkg/cmd/delete/edge_application"
	function "github.com/aziontech/azion-cli/pkg/cmd/delete/edge_function"
	edgeStorage "github.com/aziontech/azion-cli/pkg/cmd/delete/edge_storage"
	origin "github.com/aziontech/azion-cli/pkg/cmd/delete/origin"
	token "github.com/aziontech/azion-cli/pkg/cmd/delete/personal_token"
	rulesEngine "github.com/aziontech/azion-cli/pkg/cmd/delete/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmd/delete/variables"
	"github.com/aziontech/azion-cli/pkg/cmd/delete/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion delete --help
		$ azion delete edge-application -h
		$ azion delete domain -h
		$ azion delete origin -h
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplication.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(workloads.NewCmd(f))
	cmd.AddCommand(token.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(function.NewCmd(f))
	cmd.AddCommand(cache.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(edgeStorage.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
