package create

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/create"
	domains "github.com/aziontech/azion-cli/pkg/cmd/create/domains"
	edgeApplications "github.com/aziontech/azion-cli/pkg/cmd/create/edge_applications"
	rulesEngine "github.com/aziontech/azion-cli/pkg/cmd/create/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
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
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(edgeApplications.NewCmd(f))
	cmd.AddCommand(rulesEngine.NewCmd(f))
	cmd.AddCommand(domains.NewCmd(f))

	cmd.Flags().BoolP("help", "h", false, msg.FlagHelp)
	return cmd
}
