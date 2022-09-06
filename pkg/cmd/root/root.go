package cmd

import (
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/pkg/cmd/configure"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmd/webapp"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "azioncli",
		Short: "Azion-CLI",
		Long: `Interact easily with Azion services.

Azion CLI is currently in Beta. We’d love to hear your feedback at https://forms.gle/uBBkyXZCVcrgpvAB8`,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Version: version.BinVersion,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceErrors: true,
	}

	rootCmd.SetIn(f.IOStreams.In)
	rootCmd.SetOut(f.IOStreams.Out)
	rootCmd.SetErr(f.IOStreams.Err)

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		rootHelpFunc(f, cmd, args)
	})

	rootCmd.AddCommand(configure.NewCmd(f))
	rootCmd.AddCommand(version.NewCmd(f))
	rootCmd.AddCommand(edge_services.NewCmd(f))
	rootCmd.AddCommand(edge_functions.NewCmd(f))
	rootCmd.AddCommand(webapp.NewCmd(f))

	return rootCmd
}

func Execute() {
	streams := iostreams.System()
	httpClient := &http.Client{
		Timeout: 10 * time.Second, // TODO: Configure this somewhere
	}

	// TODO: Ignoring errors since the file might not exist, maybe warn the user?
	tok, _ := token.ReadFromDisk()
	viper.SetEnvPrefix("AZIONCLI")
	viper.AutomaticEnv()
	viper.SetDefault("token", tok)
	viper.SetDefault("api_url", constants.ApiURL)

	factory := &cmdutil.Factory{
		HttpClient: httpClient,
		IOStreams:  streams,
		Config:     viper.GetViper(),
	}

	cmd := NewRootCmd(factory)

	cobra.CheckErr(cmd.Execute())
}
