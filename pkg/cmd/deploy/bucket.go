package deploy

import (
	"context"
	"errors"

	"github.com/aziontech/azionapi-go-sdk/storage"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
	"go.uber.org/zap"
)

func (cmd *DeployCmd) doBucket(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Bucket != "" || (conf.Preset == "javascript" || conf.Preset == "typescript") {
		return nil
	}

	nameBucket := conf.Name

	logger.FInfo(cmd.Io.Out, msg.ProjectNameMessage)
	for {
		err := client.CreateBucket(ctx, api.RequestBucket{
			BucketCreate: storage.BucketCreate{Name: nameBucket, EdgeAccess: storage.READ_WRITE}})
		if err != nil {
			// if the name is already in use, we ask for another one
			if errors.Is(err, utils.ErrorNameInUse) {
				if NoPrompt {
					return err
				}
				logger.FInfo(cmd.Io.Out, msg.BucketInUse)
				if Auto {
					nameBucket = thoth.GenerateName()
				} else {
					nameBucket, err = askForInput(msg.AskInputName, thoth.GenerateName())
					if err != nil {
						return err
					}
				}
				conf.Bucket = nameBucket
				continue
			}
			return err
		}
		break
	}

	conf.Bucket = nameBucket
	err := cmd.WriteAzionJsonContent(conf)
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
