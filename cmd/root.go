package cmd

import (
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/cmd/configure"
	"github.com/aziontech/azion-cli/cmd/edge_services"
	"github.com/aziontech/azion-cli/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/spf13/cobra"
)

var rootToken string

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "azioncli",
		Short: "Azion-CLI",
		Long:  `This is a placeholder description used while the actual description is still not ready.`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Version: version.BinVersion,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.SetIn(f.IOStreams.In)
	rootCmd.SetOut(f.IOStreams.Out)
	rootCmd.SetErr(f.IOStreams.Err)

	rootCmd.PersistentFlags().StringVarP(&rootToken, "token", "t", "", "Use provided token")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Makes azioncli verbose during the operation")

	rootCmd.AddCommand(configure.NewCmd(f))
	rootCmd.AddCommand(version.NewCmd(f))
	rootCmd.AddCommand(edge_services.NewCmdEdgeServices(f))

	return rootCmd
}

func Execute() {
	factory := &cmdutil.Factory{
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Timeout: 10 * time.Second, // TODO: Configure this somewhere
			}, nil
		},
		IOStreams: iostreams.System(),
	}
	cmd := NewRootCmd(factory)
	cobra.CheckErr(cmd.Execute())
}
