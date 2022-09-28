package version

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

var BinVersion = "development"

// versionCmd represents the version command
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   msg.VersionUsage,
		Short: msg.VersionShortDescription,
		Long:  msg.VersionLongDescription,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(f.IOStreams.Out, "azioncli version "+BinVersion)
		},
	}

	versionCmd.SetIn(f.IOStreams.In)
	versionCmd.SetOut(f.IOStreams.Out)
	versionCmd.SetErr(f.IOStreams.Err)

	return versionCmd
}
