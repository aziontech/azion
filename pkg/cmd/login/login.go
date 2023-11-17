package login

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/login"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	cmdPersToken "github.com/aziontech/azion-cli/pkg/cmd/create/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	var username, password string

	cmd := &cobra.Command{
		Use:   msg.Usage,
		Short: msg.ShortDescription,
		Long:  msg.LongDescription,
		Example: heredoc.Doc(`
		$ azion login --help
		$ azion login --username fulanodasilva@gmail.com --password "senhasecreta"
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			if !cmd.Flags().Changed("username") {
				answers, err := utils.AskInput(msg.AskUsername)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				username = answers
			}

			if !cmd.Flags().Changed("password") {
				answers, err := utils.AskPassword(msg.AskPassword)

				if err != nil {
					logger.Debug("Error while parsing answer", zap.Error(err))
					return utils.ErrorParseResponse
				}

				password = answers
			}

			client, err := token.New(&token.Config{Client: f.HttpClient})
			if err != nil {
				return err
			}

			resp, err := client.Create(b64(username, password))
			if err != nil {
				return err
			}

			tokenValid, err := client.Validate(&resp.Token)
			if err != nil {
				return err
			}

			if !tokenValid {
				return fmt.Errorf(msg.ErrorTokenCreateInvalid)
			}

			date, err := cmdPersToken.ParseExpirationDate(time.Now(), "1m")
			if err != nil {
				return err
			}

			request := api.Request{}
			request.SetName(username)
			request.SetExpiresAt(date)

			clientPersonalToken := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
			response, err := clientPersonalToken.Create(context.Background(), &request)
			if err != nil {
				return fmt.Errorf(msg.ErrorLogin, err.Error())
			}

			settings := token.Settings{
				UUID:  response.GetUuid(),
				Token: response.GetKey(),
			}

			byteSettings, err := toml.Marshal(settings)
			if err != nil {
				return err
			}

			_, err = client.Save(byteSettings)
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

func b64(username, password string) string {
	str := utils.Concat(username, ":", password)
	b := []byte(str)
	return base64.StdEncoding.EncodeToString(b)
}
