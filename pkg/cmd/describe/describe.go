package describe

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/describe"
	cache "github.com/aziontech/azion-cli/pkg/cmd/describe/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmd/describe/domains"
	edgeApplications "github.com/aziontech/azion-cli/pkg/cmd/describe/edge_applications"
	function "github.com/aziontech/azion-cli/pkg/cmd/describe/edge_function"
	origin "github.com/aziontech/azion-cli/pkg/cmd/describe/origin"
	ruleEngine "github.com/aziontech/azion-cli/pkg/cmd/describe/rules_engine"
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
		$ azion describe edge-application
		$ azion describe domain 
		$ azion describe origin 
		$ azion describe rule-engine 
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

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
