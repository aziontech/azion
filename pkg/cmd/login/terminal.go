package login

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	msg "github.com/aziontech/azion-cli/messages/login"
	api "github.com/aziontech/azion-cli/pkg/api/personal_token"
	cmdPersToken "github.com/aziontech/azion-cli/pkg/cmd/create/personal_token"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func terminalLogin(cmd *cobra.Command, f *cmdutil.Factory) error {

	if !cmd.Flags().Changed("username") {
		answers, err := utils.AskInput(msg.AskUsername)
		if err != nil {
			return err
		}

		username = answers
	}

	if !cmd.Flags().Changed("password") {
		answers, err := utils.AskPassword(msg.AskPassword)
		if err != nil {
			return err
		}

		password = answers
	}

	client, err := token.New(&token.Config{Client: f.HttpClient})
	if err != nil {
		return err
	}

	resp, err := client.Create(b64(username, password))
	if err != nil {
		logger.Debug("Error while creating basic token", zap.Error(err))
		return err
	}

	err = validateToken(client)
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

	clientPersonalToken := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	response, err := clientPersonalToken.Create(context.Background(), &request)
	if err != nil {
		return fmt.Errorf(msg.ErrorLogin, err.Error())
	}

	tokenValue = response.GetKey()
	uuid = response.GetUuid()

	err = validateToken(client)
	if err != nil {
		return err
	}

	return nil
}

func b64(username, password string) string {
	str := utils.Concat(username, ":", password)
	b := []byte(str)
	return base64.StdEncoding.EncodeToString(b)
}
