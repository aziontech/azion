package main

import (
	cmd "github.com/aziontech/azion-cli/pkg/cmd/root"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"os"
)

var args string = "bash"

func init() {
	if len(os.Args) > 1 {
		args = os.Args[1]
	}
}

func main() {
	cmd := cmd.NewRootCmd(&cmdutil.Factory{IOStreams: iostreams.System()})
	GenerateAutocomplete(cmd, args)
}

func GenerateAutocomplete(rootCmd *cobra.Command, args string) {
	switch args {
	case "bash":
		rootCmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		rootCmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		rootCmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		rootCmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	}
}
