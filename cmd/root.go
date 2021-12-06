package cmd

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/token"
	"github.com/spf13/cobra"
)

var BinVersion = "development"
var rtoken string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azion",
	Short: "Azion-CLI",
	Long:  `This is a placeholder description used while the actual description is still not ready.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	Version: BinVersion,
	RunE: func(cmd *cobra.Command, args []string) error {
		if rtoken == "" {
			rtoken = tokenLoadFromConf()
			fmt.Println("Using saved token: " + rtoken)
		} else {
			fmt.Println("Using command line token: " + rtoken)
		}
		return nil
	},
}

func tokenLoadFromConf() string {
	c := &http.Client{Timeout: 10 * time.Second}
	t := token.NewToken(c)
	dToken, _ := t.ReadFromDisk()
	isTokenValid, _ := t.Validate(&dToken)
	if isTokenValid {
		return dToken
	}

	return ""
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVarP(&rtoken, "token", "t", "", "Use provided token")
}
