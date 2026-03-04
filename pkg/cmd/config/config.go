package config

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/config"
	"github.com/aziontech/azion-cli/pkg/cmd/config/apply"
	configdelete "github.com/aziontech/azion-cli/pkg/cmd/config/delete"
	configinit "github.com/aziontech/azion-cli/pkg/cmd/config/init"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion config --help
		$ azion config init
		$ azion config apply
		$ azion config apply --config-dir ./my-project
		$ azion config delete
		$ azion config delete --force
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(apply.NewCmd(f))
	cmd.AddCommand(configinit.NewCmd(f))
	cmd.AddCommand(configdelete.NewCmd(f))

	return cmd
}
