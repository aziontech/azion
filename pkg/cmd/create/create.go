package create

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create"
	edgeApplications "github.com/aziontech/azion-cli/pkg/cmd/create/applications"
	cacheSetting "github.com/aziontech/azion-cli/pkg/cmd/create/cache_setting"
	edgeConnector "github.com/aziontech/azion-cli/pkg/cmd/create/connector"
	edgeFunction "github.com/aziontech/azion-cli/pkg/cmd/create/function"
	functionInstance "github.com/aziontech/azion-cli/pkg/cmd/create/function_instance"
	networkList "github.com/aziontech/azion-cli/pkg/cmd/create/network_list"
	origin "github.com/aziontech/azion-cli/pkg/cmd/create/origin"
	token "github.com/aziontech/azion-cli/pkg/cmd/create/personal_token"
	profile "github.com/aziontech/azion-cli/pkg/cmd/create/profile"
	rulesEngine "github.com/aziontech/azion-cli/pkg/cmd/create/rules_engine"
	edgeStorage "github.com/aziontech/azion-cli/pkg/cmd/create/storage"
	"github.com/aziontech/azion-cli/pkg/cmd/create/variables"
	workloaddeployment "github.com/aziontech/azion-cli/pkg/cmd/create/workload_deployment"
	"github.com/aziontech/azion-cli/pkg/cmd/create/workloads"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription, Example: heredoc.Doc(`
		$ azion create --help
		$ azion create application -h
		$ azion create connector -h
		$ azion create workload -h
		$ azion create network-list -h
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplications.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(token.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(cacheSetting.NewCmd(f))
	cmd.AddCommand(edgeFunction.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(edgeStorage.NewCmd(f))
	cmd.AddCommand(workloads.NewCmd(f))
	cmd.AddCommand(workloaddeployment.NewCmd(f))
	cmd.AddCommand(edgeConnector.NewCmd(f))
	cmd.AddCommand(functionInstance.NewCmd(f))
	cmd.AddCommand(profile.NewCmd(f))
	cmd.AddCommand(networkList.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
