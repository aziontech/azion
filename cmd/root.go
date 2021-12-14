package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/token"
	"github.com/spf13/cobra"
)

var BinVersion = "development"
var rootToken string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azioncli",
	Short: "Azion-CLI",
	Long:  `This is a placeholder description used while the actual description is still not ready.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Version: BinVersion,
	RunE: func(cmd *cobra.Command, args []string) error {
		if rootToken == "" {
			rootToken, err := tokenLoadFromConf()
			if err != nil {
				return err
			}
			fmt.Println("Using saved token: " + rootToken)
		} else {
			fmt.Println("Using command line token: " + rootToken)
		}
		return nil
	},
}

func tokenLoadFromConf() (string, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	t := token.NewToken(c)
	return t.ReadFromDisk()
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&rootToken, "token", "t", "", "Use provided token")
}
