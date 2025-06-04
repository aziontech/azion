package version

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/output"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			versionOut := output.GeneralOutput{
				Msg:   fmt.Sprintf("Azion CLI %s\n", BinVersion),
				Out:   f.IOStreams.Out,
				Flags: f.Flags,
			}
			return output.Print(&versionOut)
		},
	}

	versionCmd.SetIn(f.IOStreams.In)
	versionCmd.SetOut(f.IOStreams.Out)
	versionCmd.SetErr(f.IOStreams.Err)

	versionCmd.Flags().BoolP("help", "h", false, msg.VersionHelpFlag)

	return versionCmd
}
