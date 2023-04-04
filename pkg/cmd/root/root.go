package cmd

import (
	"net/http"
	"time"

	completion "github.com/aziontech/azion-cli/pkg/cmd/completion"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions_instances"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmd/cache_settings"
	"github.com/aziontech/azion-cli/pkg/cmd/configure"
	"github.com/aziontech/azion-cli/pkg/cmd/domains"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_functions"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_services"
	"github.com/aziontech/azion-cli/pkg/cmd/origins"
	"github.com/aziontech/azion-cli/pkg/cmd/rules_engine"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/constants"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   msg.RootUsage,
		Short: msg.RootShortDescription,
		Long:  msg.RootLongDescription,
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
	rootCmd.AddCommand(completion.NewCmd(f))
	rootCmd.AddCommand(version.NewCmd(f))
	rootCmd.AddCommand(edge_services.NewCmd(f))
	rootCmd.AddCommand(edge_functions.NewCmd(f))
	rootCmd.AddCommand(edge_applications.NewCmd(f))
	rootCmd.AddCommand(domains.NewCmd(f))
	rootCmd.AddCommand(origins.NewCmd(f))
	rootCmd.AddCommand(rules_engine.NewCmd(f))
	rootCmd.AddCommand(cache_settings.NewCmd(f))
	rootCmd.AddCommand(edge_functions_instances.NewCmd(f))
	rootCmd.Flags().BoolP("help", "h", false, msg.RootHelpFlag)
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
	viper.SetDefault("storage_url", constants.StorageApiURL)

	factory := &cmdutil.Factory{
		HttpClient: httpClient,
		IOStreams:  streams,
		Config:     viper.GetViper(),
	}

	cmd := NewRootCmd(factory)

	cobra.CheckErr(cmd.Execute())
}
