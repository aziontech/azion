package completion

import (
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
	"os"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 "Generate completion script",
		Long:                  "To load completions",
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		Annotations: map[string]string{
			"Category": "skip",
		},
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				err := cmd.Root().GenBashCompletion(os.Stdout)
				if err != nil {
					return
				}
			case "zsh":
				err := cmd.Root().GenZshCompletion(os.Stdout)
				if err != nil {
					return
				}
			case "fish":
				err := cmd.Root().GenFishCompletion(os.Stdout, true)
				if err != nil {
					return
				}
			case "powershell":
				err := cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
				if err != nil {
					return
				}
			}
		},
	}
}
