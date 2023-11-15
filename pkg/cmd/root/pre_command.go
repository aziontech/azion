package root

import (
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/root"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

type PreCmd struct {
	token  string
	config string
}

// doPreCommandCheck carry out all pre-cmd checks needed
func doPreCommandCheck(cmd *cobra.Command, f *cmdutil.Factory, pre PreCmd) error {

	if err := setConfigPath(cmd, pre.config); err != nil {
		return err
	}

	if err := checkTokenSent(cmd, f, pre.token); err != nil {
		return err
	}

	return nil
}

func setConfigPath(cmd *cobra.Command, cfg string) error {

	if cmd.Flags().Changed("config") {
		config.SetPath(cfg)
		return nil
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

		strToken := token.Settings{Token: configureToken}
		bStrToken, err := toml.Marshal(strToken)
		if err != nil {
			return err
		}

		filePath, err := t.Save(bStrToken)
		if err != nil {
			return err
		}

		logger.LogSuccess(f.IOStreams.Out, fmt.Sprintf(msg.TokenSavedIn, filePath))
		logger.FInfo(f.IOStreams.Out, msg.TokenUsedIn)
	}

	return nil
}
