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
	"errors"
	"net/http"
	"time"

	"github.com/aziontech/azion-cli/token"
	"github.com/spf13/cobra"
)

var ctoken string

// configureCmd represents the configure command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure parameters and credentials",
	Long:  `This command configures cli parameters and credentials used for connecting to our services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c := &http.Client{Timeout: 10 * time.Second}
		t := token.NewToken(c)

		if ctoken == "" {
			return errors.New("token not provided, loading the saved one")
		}

		valid, err := t.Validate(&ctoken)
		if err != nil {
			return err
		}

		if !valid {
			return errors.New("invalid token")
		}

		if t.Save() != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)

	configureCmd.Flags().StringVarP(&ctoken, "token", "t", "", "Validate token and save it in $HOME_DIR/.azion/credentials")
}
