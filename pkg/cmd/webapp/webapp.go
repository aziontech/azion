package webapp

import (
	buildCmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/build"
	initCmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/init"
	publishCmd "github.com/aziontech/azion-cli/pkg/cmd/webapp/publish"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	webappCmd := &cobra.Command{
		Use:   "webapp",
		Short: "Initialize, Build and Publish your Web Applications",
		Long:  `You can Initialize, Build and Publish your Web Applications`,
		Annotations: map[string]string{
			"Category": "Build",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	webappCmd.AddCommand(initCmd.NewCmd(f))
	webappCmd.AddCommand(buildCmd.NewCmd(f))
	webappCmd.AddCommand(publishCmd.NewCmd(f))

	return webappCmd
}
