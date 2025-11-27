package describe

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe"
	edgeApplications "github.com/aziontech/azion-cli/pkg/cmd/describe/applications"
	cache "github.com/aziontech/azion-cli/pkg/cmd/describe/cache_setting"
	edgeConnector "github.com/aziontech/azion-cli/pkg/cmd/describe/connector"
	function "github.com/aziontech/azion-cli/pkg/cmd/describe/function"
	functioninstance "github.com/aziontech/azion-cli/pkg/cmd/describe/function_instance"
	networklist "github.com/aziontech/azion-cli/pkg/cmd/describe/network_list"
	origin "github.com/aziontech/azion-cli/pkg/cmd/describe/origin"
	"github.com/aziontech/azion-cli/pkg/cmd/describe/personal_token"
	ruleEngine "github.com/aziontech/azion-cli/pkg/cmd/describe/rules_engine"
	edgeStorage "github.com/aziontech/azion-cli/pkg/cmd/describe/storage"
	"github.com/aziontech/azion-cli/pkg/cmd/describe/variables"
	"github.com/aziontech/azion-cli/pkg/cmd/describe/workloads"
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
		$ azion describe application -h
		$ azion describe workload -h
		$ azion describe origin -h
		$ azion describe network-list -h
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplications.NewCmd(f))
	cmd.AddCommand(ruleEngine.NewCmd(f))
	cmd.AddCommand(workloads.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(cache.NewCmd(f))
	cmd.AddCommand(function.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(edgeStorage.NewCmd(f))
	cmd.AddCommand(personal_token.NewCmd(f))
	cmd.AddCommand(edgeConnector.NewCmd(f))
	cmd.AddCommand(functioninstance.NewCmd(f))
	cmd.AddCommand(networklist.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
