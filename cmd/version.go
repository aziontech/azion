package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Returns bin version",
	Long:  `Returns the version of the binary.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Azion version: " + BinVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
