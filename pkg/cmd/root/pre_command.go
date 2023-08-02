package cmd

import (
	"fmt"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

// doPreCommandCheck carries out all pre-cmd checks needed
func doPreCommandCheck(cmd *cobra.Command, f *cmdutil.Factory, configureToken string) error {

	err := checkTokenSent(cmd, f, configureToken)
	if err != nil {
		return err
	}

	return nil
}

func checkTokenSent(cmd *cobra.Command, f *cmdutil.Factory, configureToken string) error {

	// if global --token flag was sent, verify it and save it locally
	if cmd.Flags().Changed("token") {
		t, err := token.New(&token.Config{
			Client: f.HttpClient,
			Out:    f.IOStreams.Out,
		})
		if err != nil {
			return fmt.Errorf("%s: %w", utils.ErrorTokenManager, err)
		}

		if configureToken == "" {
			return utils.ErrorTokenNotProvided
		}

		valid, err := t.Validate(&configureToken)
		if err != nil {
			return err
		}

		if !valid {
			return utils.ErrorInvalidToken
		}

		if err := t.Save(); err != nil {
			return err
		}
	}

	return nil
}