package webapp

import (
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/webapp"
	buildCmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/build"
	initCmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/init"
	publishCmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/publish"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	webappCmd := &cobra.Command{
		Use:   msg.WebappUsage,
		Short: msg.WebappShortDescription,
		Long:  msg.WebappLongDescription,
		Example: heredoc.Doc(`
		$ azioncli webapp --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	webappCmd.AddCommand(initCmd.NewCmd(f))
	b := buildCmd.NewCmd(f)
	webappCmd.AddCommand(b)
	webappCmd.AddCommand(publishCmd.NewCmd(f, buildCmd.NewBuildCmd(f)))
	webappCmd.Flags().BoolP("help", "h", false, msg.WebappFlagHelp)

	return webappCmd
}
