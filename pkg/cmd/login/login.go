package login

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/login"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var username, password, tokenValue, uuid string

func NewCmd(f *cmdutil.Factory) *cobra.Command {

	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion login --help
		$ azion login --username fulanodasilva@gmail.com --password "senhasecreta"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			answer, err := selectLoginMode()
			if err != nil {
				return err
			}

			switch {
			case strings.Contains(answer, "browser"):
				browserLogin(f)
			case strings.Contains(answer, "terminal"):
				terminalLogin(cmd, f)
			default:
				return msg.ErrorInvalidLogin
			}

			client, err := token.New(&token.Config{Client: f.HttpClient})
			if err != nil {
				return err
			}

			err = validateToken(client, tokenValue)
			if err != nil {
				return err
			}

			err = saveSettings(client)
			if err != nil {
				return err
			}

			logger.LogSuccess(f.IOStreams.Out, msg.Success)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&username, "username", "", msg.FlagUsername)
	flags.StringVar(&password, "password", "", msg.FlagPassword)
	flags.BoolP("help", "h", false, msg.FlagHelp)

	return cmd
}

func selectLoginMode() (string, error) {
	answer := ""
	prompt := &survey.Select{
		Message: "Choose a mode:",
		Options: []string{"Log in via browser", "Log in via terminal"},
	}
	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func saveSettings(client *token.Token) error {
	settings := token.Settings{
		UUID:  uuid,
		Token: tokenValue,
	}

	byteSettings, err := toml.Marshal(settings)
	if err != nil {
		logger.Debug("Error while marshalling toml file", zap.Error(err))
		return err
	}

	_, err = client.Save(byteSettings)
	if err != nil {
		logger.Debug("Error while saving settings", zap.Error(err))
		return err
	}
	return nil
}

func validateToken(client *token.Token, token string) error {
	tokenValid, err := client.Validate(&token)
	if err != nil {
		logger.Debug("Error while validating the token", zap.Error(err))
		return err
	}

	if !tokenValid {
		return fmt.Errorf(msg.ErrorTokenCreateInvalid)
	}

	return nil
}
