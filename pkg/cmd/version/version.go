package version

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	msg "github.com/aziontech/azion-cli/pkg/messages/version"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var BinVersion = "development"

// versionCmd represents the version command
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   msg.VersionUsage,
		Short: msg.VersionShortDescription,
		Long:  msg.VersionLongDescription,
		Example: heredoc.Doc(`
		$ azion version
        `),
		Run: func(cmd *cobra.Command, args []string) {
			color.New(color.Bold).Fprintln(f.IOStreams.Out, "Azion CLI "+BinVersion+"\n")
		},
	}

	versionCmd.SetIn(f.IOStreams.In)
	versionCmd.SetOut(f.IOStreams.Out)
	versionCmd.SetErr(f.IOStreams.Err)

	versionCmd.Flags().BoolP("help", "h", false, msg.VersionHelpFlag)

	return versionCmd
}
