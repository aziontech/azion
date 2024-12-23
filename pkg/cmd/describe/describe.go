package describe

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe"
	cache "github.com/aziontech/azion-cli/pkg/cmd/describe/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmd/describe/domains"
	edgeApplications "github.com/aziontech/azion-cli/pkg/cmd/describe/edge_applications"
	function "github.com/aziontech/azion-cli/pkg/cmd/describe/edge_function"
	edgeStorage "github.com/aziontech/azion-cli/pkg/cmd/describe/edge_storage"
	origin "github.com/aziontech/azion-cli/pkg/cmd/describe/origin"
	"github.com/aziontech/azion-cli/pkg/cmd/describe/personal_token"
	ruleEngine "github.com/aziontech/azion-cli/pkg/cmd/describe/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmd/describe/variables"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion describe --help
		$ azion describe edge-application -h
		$ azion describe domain -h
		$ azion describe origin -h
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplications.NewCmd(f))
	cmd.AddCommand(ruleEngine.NewCmd(f))
	cmd.AddCommand(domains.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(cache.NewCmd(f))
	cmd.AddCommand(function.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(edgeStorage.NewCmd(f))
	cmd.AddCommand(personal_token.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
