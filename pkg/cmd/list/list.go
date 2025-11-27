package list

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/list"
	applications "github.com/aziontech/azion-cli/pkg/cmd/list/applications"
	cache "github.com/aziontech/azion-cli/pkg/cmd/list/cache_setting"
	connector "github.com/aziontech/azion-cli/pkg/cmd/list/connector"
	function "github.com/aziontech/azion-cli/pkg/cmd/list/function"
	functioninstance "github.com/aziontech/azion-cli/pkg/cmd/list/function_instance"
	networklist "github.com/aziontech/azion-cli/pkg/cmd/list/network_list"
	origin "github.com/aziontech/azion-cli/pkg/cmd/list/origin"
	token "github.com/aziontech/azion-cli/pkg/cmd/list/personal_token"
	rule "github.com/aziontech/azion-cli/pkg/cmd/list/rule_engine"
	storage "github.com/aziontech/azion-cli/pkg/cmd/list/storage"
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
		$ azion list network-list -h
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(applications.NewCmd(f))
	cmd.AddCommand(rule.NewCmd(f))
	cmd.AddCommand(workloads.NewCmd(f))
	cmd.AddCommand(wdeployments.NewCmd(f))
	cmd.AddCommand(token.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(cache.NewCmd(f))
	cmd.AddCommand(function.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(storage.NewCmd(f))
	cmd.AddCommand(connector.NewCmd(f))
	cmd.AddCommand(functioninstance.NewCmd(f))
	cmd.AddCommand(networklist.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
