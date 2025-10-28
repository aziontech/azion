package login

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/login"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/skratchdot/open-golang/open"
)

var (
	username, password, tokenValue, uuid string
	userInfo                             token.UserInfo
	createNewProfile                     bool
	newProfileName                       string
)

var confirmFn = utils.Confirm

type login struct {
	factory     *cmdutil.Factory
	askOne      func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
	run         func(input string) error
	server      Server
	token       token.TokenInterface
	marshalToml func(v interface{}) ([]byte, error)
	askInput    func(msg string) (string, error)
	askPassword func(msg string) (string, error)
}

func New(f *cmdutil.Factory) *cobra.Command {
	return cmd(factory(f))
}

func factory(f *cmdutil.Factory) *login {
	tk := token.New(&token.Config{Client: f.HttpClient})
	return &login{
		factory:     f,
		askOne:      survey.AskOne,
		run:         open.Run,
		server:      &http.Server{Addr: ":8080"},
		token:       tk,
		marshalToml: toml.Marshal,
		askInput:    utils.AskInput,
		askPassword: utils.AskPassword,
	}
}

func cmd(l *login) *cobra.Command {
	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion login --help
		$ azion login --username fulanodasilva@gmail.com --password "senhasecreta"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			createNewProfile = confirmFn(l.factory.GlobalFlagAll, msg.QuestionCreateProfile, true)
			if createNewProfile {
				profileNameInput, err := utils.AskInput(msg.AskProfileName)
				if err != nil {
					return fmt.Errorf(msg.ErrorGetProfileName.Error(), err)
				}
				newProfileName = profileNameInput
			}

			answer, err := l.selectLoginMode()
			if err != nil {
				return err
			}

			switch {
			case strings.Contains(answer, "browser"):
				err := l.browserLogin(l.server)
				if err != nil {
					return err
				}
			case strings.Contains(answer, "terminal"):
				err := l.terminalLogin(cmd)
				if err != nil {
					return err
				}
			default:
				return msg.ErrorInvalidLogin
			}

			err = l.validateToken(tokenValue)
			if err != nil {
				return err
			}

			err = l.saveSettings()
			if err != nil {
				return err
			}

			loginOut := output.GeneralOutput{
				Msg:   fmt.Sprintf(msg.Success),
				Out:   l.factory.IOStreams.Out,
				Flags: l.factory.Flags,
			}
			return output.Print(&loginOut)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&username, "username", "", msg.FlagUsername)
	flags.StringVar(&password, "password", "", msg.FlagPassword)
	flags.BoolP("help", "h", false, msg.FlagHelp)

	return cmd
}

func (l *login) selectLoginMode() (answer string, err error) {
	prompt := &survey.Select{
		Message: "Choose a login method:",
		Options: []string{"Log in via browser", "Log in via terminal"},
	}
	err = l.askOne(prompt, &answer)
	if err != nil {
		return "", err
	}
	return answer, nil
}

func (l *login) saveSettings() error {
	settings := token.Settings{
		UUID:        uuid,
		Token:       tokenValue,
		ClientId:    userInfo.Results.ClientID,
		Email:       userInfo.Results.Email,
		S3AccessKey: "",
		S3SecretKey: "",
		S3Bucket:    "",
	}

	var profileName string
	if createNewProfile {
		profileName = newProfileName

		err := token.WriteSettings(settings, profileName)
		if err != nil {
			return err
		}

		profile := token.Profile{Name: profileName}
		err = token.WriteProfiles(profile)
		if err != nil {
			return fmt.Errorf(msg.ErrorSetActiveProfile.Error(), err)
		}

		fmt.Fprintf(l.factory.IOStreams.Out, msg.ProfileCreated+"\n", profileName)
	} else {
		profileName = l.factory.GetActiveProfile()
		err := token.WriteSettings(settings, profileName)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(l.factory.IOStreams.Out, msg.TokenSavedToProfile+"\n", profileName)
	return nil
}

func (l *login) validateToken(token string) error {
	tokenValid, user, err := l.token.Validate(&token)
	userInfo = user
	if err != nil {
		logger.Debug("Error while validating the token", zap.Error(err))
		return err
	}

	if !tokenValid {
		return msg.ErrorTokenCreateInvalid
	}

	return nil
}
