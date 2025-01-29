package deploy

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aziontech/azionapi-go-sdk/storage"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func (cmd *DeployCmd) doBucket(
	client *api.Client,
	ctx context.Context,
	conf *contracts.AzionApplicationOptions,
	msgs *[]string) error {
	if conf.Bucket != "" || (conf.Preset == "javascript" || conf.Preset == "typescript") {
		return nil
	}

	logger.FInfoFlags(cmd.Io.Out, msg.ProjectNameMessage, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msg.ProjectNameMessage)
	nameBucket := utils.ReplaceInvalidCharsBucket(conf.Name)

	err := client.CreateBucket(ctx, api.RequestBucket{
		BucketCreate: storage.BucketCreate{Name: nameBucket, EdgeAccess: storage.READ_ONLY}})
	if err != nil {
		// If the name is already in use, try 10 times with different names
		for i := 0; i < 10; i++ {
			nameB := fmt.Sprintf("%s-%s", nameBucket, utils.Timestamp())
			msgf := fmt.Sprintf(msg.NameInUseBucket, nameB)
			logger.FInfoFlags(cmd.Io.Out, msgf, cmd.F.Format, cmd.F.Out)
			*msgs = append(*msgs, msgf)
			err := client.CreateBucket(ctx, api.RequestBucket{
				BucketCreate: storage.BucketCreate{Name: nameB, EdgeAccess: storage.READ_ONLY}})
			if err != nil {
				if errors.Is(err, utils.ErrorNameInUse) && i < 9 {
					continue
				}
				return err
			}
			conf.Bucket = nameB
			break
		}
	} else {
		conf.Bucket = nameBucket
	}

	msgf := fmt.Sprintf(msg.BucketSuccessful, conf.Bucket)
	logger.FInfoFlags(cmd.Io.Out, msgf, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msgf)
	return cmd.WriteAzionJsonContent(conf, ProjectConf)
}

func askForInput(msg string, defaultIn string) (string, error) {
	var userInput string
	prompt := &survey.Input{
		Message: msg,
		Default: defaultIn,
	}

	// Prompt the user for input
	err := survey.AskOne(prompt, &userInput, survey.WithKeepFilter(true), survey.WithStdio(os.Stdin, os.Stderr, os.Stdout))
	if err != nil {
		return "", err
	}
	return userInput, nil
}
