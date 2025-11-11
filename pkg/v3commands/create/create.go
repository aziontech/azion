package create

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create"
	profile "github.com/aziontech/azion-cli/pkg/cmd/create/profile"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	cacheSetting "github.com/aziontech/azion-cli/pkg/v3commands/create/cache_setting"
	domain "github.com/aziontech/azion-cli/pkg/v3commands/create/domain"
	edgeApplications "github.com/aziontech/azion-cli/pkg/v3commands/create/edge_applications"
	edgeFunction "github.com/aziontech/azion-cli/pkg/v3commands/create/edge_function"
	edgeStorage "github.com/aziontech/azion-cli/pkg/v3commands/create/edge_storage"
	origin "github.com/aziontech/azion-cli/pkg/v3commands/create/origin"
	token "github.com/aziontech/azion-cli/pkg/v3commands/create/personal_token"
	rulesEngine "github.com/aziontech/azion-cli/pkg/v3commands/create/rules_engine"
	"github.com/aziontech/azion-cli/pkg/v3commands/create/variables"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion create --help
		$ azion create edge-application -h
		$ azion create rules-engine -h
		$ azion create domain -h
		$ azion create personal-token -h
		$ azion create origin -h
		$ azion create cache-setting -h
		$ azion create edge-function -h
		$ azion create variables -h
		$ azion create edge-storage -h
		$ azion create profile -h
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplications.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(domain.NewCmd(f))
	cmd.AddCommand(token.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(cacheSetting.NewCmd(f))
	cmd.AddCommand(edgeFunction.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(edgeStorage.NewCmd(f))
	cmd.AddCommand(profile.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
