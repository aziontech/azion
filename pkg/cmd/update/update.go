package update

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/update"
	application "github.com/aziontech/azion-cli/pkg/cmd/update/application"
	cacheSetting "github.com/aziontech/azion-cli/pkg/cmd/update/cache_setting"
	connector "github.com/aziontech/azion-cli/pkg/cmd/update/connector"
	function "github.com/aziontech/azion-cli/pkg/cmd/update/function"
	functionInstance "github.com/aziontech/azion-cli/pkg/cmd/update/function_instance"
	networkList "github.com/aziontech/azion-cli/pkg/cmd/update/network_list"
	origin "github.com/aziontech/azion-cli/pkg/cmd/update/origin"
	rulesEngine "github.com/aziontech/azion-cli/pkg/cmd/update/rules_engine"
	storage "github.com/aziontech/azion-cli/pkg/cmd/update/storage"
	"github.com/aziontech/azion-cli/pkg/cmd/update/variables"
	"github.com/aziontech/azion-cli/pkg/cmd/update/workloads"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion update --help
		$ azion update application -h
		$ azion update connector -h
		$ azion update workload -h
		$ azion update network-list -h
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(application.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(function.NewCmd(f))
	cmd.AddCommand(cacheSetting.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(storage.NewCmd(f))
	cmd.AddCommand(workloads.NewCmd(f))
	cmd.AddCommand(connector.NewCmd(f))
	cmd.AddCommand(functionInstance.NewCmd(f))
	cmd.AddCommand(networkList.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
