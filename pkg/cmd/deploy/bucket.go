package deploy

import (
	"context"
	"errors"

	sdk "github.com/aziontech/azionapi-go-sdk/storage"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/aziontech/azionapi-go-sdk/storage"
	thoth "github.com/aziontech/go-thoth"
	"go.uber.org/zap"
)

func (cmd *DeployCmd) doBucket(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Bucket != "" || (conf.Template == "javascript" || conf.Template == "typescript") {
		return nil
	}

	var err error
	name := conf.Name

	logger.FInfo(cmd.Io.Out, msg.ProjectNameMessage)
	for {
		err = client.CreateBucket(ctx, api.RequestBucket{
			BucketCreate: sdk.BucketCreate{Name: name, EdgeAccess: storage.READ_WRITE}})
		if err != nil {
			// if the name is already in use, we ask for another one
			if errors.Is(err, utils.ErrorNameInUse) {
				logger.FInfo(cmd.Io.Out, msg.BucketInUse)
				if Auto {
					name = thoth.GenerateName()
				} else {
					name, err = askForInput(msg.AskInputName, thoth.GenerateName())
					if err != nil {
						return err
					}
				}
				conf.Application.Name = name
				continue
			}
			return err
		}
		break
	}

	conf.Bucket = name
	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	return nil
}

func askForInput(msg string, defaultIn string) (string, error) {
	var userInput string
	prompt := &survey.Input{
		Message: msg,
		Default: defaultIn,
	}

	// Prompt the user for input
	err := survey.AskOne(prompt, &userInput, survey.WithKeepFilter(true))
	if err != nil {
		return "", err
	}
	return userInput, nil
}
