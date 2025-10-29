package delete

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete"
	profile "github.com/aziontech/azion-cli/pkg/cmd/delete/profile"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	domain "github.com/aziontech/azion-cli/pkg/v3commands/delete/domain"
	edgeApplication "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_application"
	function "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_function"
	edgeStorage "github.com/aziontech/azion-cli/pkg/v3commands/delete/edge_storage"
	origin "github.com/aziontech/azion-cli/pkg/v3commands/delete/origin"
	token "github.com/aziontech/azion-cli/pkg/v3commands/delete/personal_token"
	rulesEngine "github.com/aziontech/azion-cli/pkg/v3commands/delete/rules_engine"
	"github.com/aziontech/azion-cli/pkg/v3commands/delete/variables"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion delete --help
		$ azion delete profile -h
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplication.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(domain.NewCmd(f))
	cmd.AddCommand(token.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(function.NewCmd(f))
	// cmd.AddCommand(cache.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(edgeStorage.NewCmd(f))
	cmd.AddCommand(profile.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
