/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/aziontech/azion-cli/mocks/configure"
	"github.com/aziontech/azion-cli/token"
	"github.com/spf13/cobra"
)

var stoken string

var ttoken token.Token

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure parameters and credentials",
	Long:  `This command configures cli parameters and credentials used for connecting to our services.`,
	Run: func(cmd *cobra.Command, args []string) {
		ttoken = token.NewToken()
		ttoken.Client = configure.MockClient{}
		if stoken != "" {
			if ttoken.Validation(stoken) {
				token.Save(stoken)
			}
		} else {
			fmt.Println("Token not provided, loading the saved")
		}

	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVarP(&stoken, "token", "t", "", "Validate token and save in $HOME_DIR/.azion/credentials")
}
