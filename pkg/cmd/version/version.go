package version

import (
	"fmt"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

var BinVersion = "development"

// versionCmd represents the version command
func NewCmd(f *cmdutil.Factory) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Returns bin version",
		Long:  `Returns the version of the binary.`,
		Annotations: map[string]string{
			"IsAditional": "true",
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(f.IOStreams.Out, "Azion version: "+BinVersion)
		},
	}

	versionCmd.SetIn(f.IOStreams.In)
	versionCmd.SetOut(f.IOStreams.Out)
	versionCmd.SetErr(f.IOStreams.Err)

	return versionCmd
}
