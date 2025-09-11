package list

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list"
	cache "github.com/aziontech/azion-cli/pkg/cmd/list/cache_setting"
	edgeApplications "github.com/aziontech/azion-cli/pkg/cmd/list/edge_applications"
	edgeConnector "github.com/aziontech/azion-cli/pkg/cmd/list/edge_connector"
	function "github.com/aziontech/azion-cli/pkg/cmd/list/edge_function"
	edgeStorage "github.com/aziontech/azion-cli/pkg/cmd/list/edge_storage"
	functioninstance "github.com/aziontech/azion-cli/pkg/cmd/list/function_instance"
	origin "github.com/aziontech/azion-cli/pkg/cmd/list/origin"
	token "github.com/aziontech/azion-cli/pkg/cmd/list/personal_token"
	rule "github.com/aziontech/azion-cli/pkg/cmd/list/rule_engine"
	"github.com/aziontech/azion-cli/pkg/cmd/list/variables"
	wdeployments "github.com/aziontech/azion-cli/pkg/cmd/list/workload_deployment"
	"github.com/aziontech/azion-cli/pkg/cmd/list/workloads"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion list --help
		$ azion list application -h
		$ azion list workload -h
		$ azion list origin -h
		$ azion list function-instance -h
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplications.NewCmd(f))
	cmd.AddCommand(rule.NewCmd(f))
	cmd.AddCommand(workloads.NewCmd(f))
	cmd.AddCommand(wdeployments.NewCmd(f))
	cmd.AddCommand(token.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(cache.NewCmd(f))
	cmd.AddCommand(function.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(edgeStorage.NewCmd(f))
	cmd.AddCommand(edgeConnector.NewCmd(f))
	cmd.AddCommand(functioninstance.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
