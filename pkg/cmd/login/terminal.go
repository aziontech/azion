package login

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	msg "github.com/aziontech/azion-cli/messages/login"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	cmdPersToken "github.com/aziontech/azion-cli/pkg/cmd/create/personal_token"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func (l *login) terminalLogin(cmd *cobra.Command) error {
	if !cmd.Flags().Changed("username") {
		answer, err := l.askInput(msg.AskUsername)
		if err != nil {
			return err
		}

		username = answer
	}

	if !cmd.Flags().Changed("password") {
		answer, err := l.askPassword(msg.AskPassword)
		if err != nil {
			return err
		}

		password = answer
	}

	resp, err := l.token.Create(b64(username, password))
	if err != nil {
		logger.Debug("Error while creating basic token", zap.Error(err))
		return err
	}

	err = l.validateToken(resp.Token)
	if err != nil {
		return err
	}

	viper.SetDefault("token", resp.Token)

	date, err := cmdPersToken.ParseExpirationDate(time.Now(), "1m")
	if err != nil {
		logger.Debug("Error while formatting expiration date", zap.Error(err))
		return err
	}

	request := api.Request{}
	request.SetName(username)
	request.SetExpiresAt(date)
	clientPersonalToken := api.NewClient(l.factory.HttpClient, l.factory.Config.GetString("api_url"), l.factory.Config.GetString("token"))
	response, err := clientPersonalToken.Create(context.Background(), &request)
	if err != nil {
		return fmt.Errorf(msg.ErrorLogin.Error(), err.Error())
	}

	tokenValue = response.GetKey()
	uuid = response.GetUuid()

	return nil
}

func b64(username, password string) string {
	str := utils.Concat(username, ":", password)
	b := []byte(str)
	return base64.StdEncoding.EncodeToString(b)
}
