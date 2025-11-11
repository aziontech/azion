package delete

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/delete"
	application "github.com/aziontech/azion-cli/pkg/cmd/delete/application"
	cache "github.com/aziontech/azion-cli/pkg/cmd/delete/cache_setting"
	connector "github.com/aziontech/azion-cli/pkg/cmd/delete/connector"
	function "github.com/aziontech/azion-cli/pkg/cmd/delete/function"
	functionInstance "github.com/aziontech/azion-cli/pkg/cmd/delete/function_instance"
	origin "github.com/aziontech/azion-cli/pkg/cmd/delete/origin"
	token "github.com/aziontech/azion-cli/pkg/cmd/delete/personal_token"
	profile "github.com/aziontech/azion-cli/pkg/cmd/delete/profile"
	rulesEngine "github.com/aziontech/azion-cli/pkg/cmd/delete/rules_engine"
	storage "github.com/aziontech/azion-cli/pkg/cmd/delete/storage"
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
		$ azion delete application -h
		$ azion delete workload -h
		$ azion delete origin -h
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(application.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(workloads.NewCmd(f))
	cmd.AddCommand(token.NewCmd(f))
	cmd.AddCommand(origin.NewCmd(f))
	cmd.AddCommand(function.NewCmd(f))
	cmd.AddCommand(cache.NewCmd(f))
	cmd.AddCommand(variables.NewCmd(f))
	cmd.AddCommand(storage.NewCmd(f))
	cmd.AddCommand(connector.NewCmd(f))
	cmd.AddCommand(functionInstance.NewCmd(f))
	cmd.AddCommand(profile.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
