package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/cmd/configure"
	"github.com/aziontech/azion-cli/cmd/edge_services"
	"github.com/aziontech/azion-cli/cmd/version"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/token"
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
			if cmd.Flags().Changed("token") {
				fmt.Fprintln(f.IOStreams.Out, "Using command line token: "+rootToken)
				return nil
			}

			client, err := f.HttpClient()
			if err != nil {
				return err
			}

			t := token.NewToken(client)
			if err != nil {
				return err
			}

			tok, err := t.ReadFromDisk()
			if err != nil {
				return err
			}

			fmt.Fprintln(f.IOStreams.Out, "Using saved token: "+tok)

			return nil
		},
	}

	rootCmd.SetIn(f.IOStreams.In)
	rootCmd.SetOut(f.IOStreams.Out)
	rootCmd.SetErr(f.IOStreams.Err)

	rootCmd.PersistentFlags().StringVarP(&rootToken, "token", "t", "", "Use provided token")

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
