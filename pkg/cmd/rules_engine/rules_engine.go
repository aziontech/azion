package rules_engine

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmd/rules_engine/create"
	"github.com/aziontech/azion-cli/pkg/cmd/rules_engine/template"
	"github.com/aziontech/azion-cli/pkg/cmd/rules_engine/update"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	rulesEngineCmd := &cobra.Command{
		Use:   msg.RulesEngineUsage,
		Short: msg.RulesEngineShortDescription,
		Long:  msg.RulesEngineLongDescription,
		Example: heredoc.Doc(`
		$ azion rules_engine --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rulesEngineCmd.AddCommand(update.NewCmd(f))
	rulesEngineCmd.AddCommand(create.NewCmd(f))
	rulesEngineCmd.AddCommand(template.NewCmd(f))

	rulesEngineCmd.Flags().BoolP("help", "h", false, msg.RulesEngineFlagHelp)

	return rulesEngineCmd
}
