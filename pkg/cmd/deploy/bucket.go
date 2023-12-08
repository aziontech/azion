package deploy

import (
	"context"
	"errors"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
)

func (cmd *DeployCmd) doBucket(client *api.ClientStorage, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	var err error
	name := conf.Name

	if conf.BucketName != "" {
		return nil
	}

	logger.FInfo(cmd.Io.Out, msg.ProjectNameMessage)
	for {
		err = client.CreateBucket(ctx, name)
		// if the bucket name is already in use, we ask for another one
		if errors.Is(err, utils.ErrorBucketInUse) {
			logger.FInfo(cmd.Io.Out, msg.BucketInUse)
			projName, err := askForInput(msg.AskInputName, thoth.GenerateName())
			if err != nil {
				return err
			}
			name = projName
			continue
		}
		break
	}

	if err != nil {
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
