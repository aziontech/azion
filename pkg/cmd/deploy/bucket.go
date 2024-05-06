package deploy

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/aziontech/azionapi-go-sdk/storage"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
)

func (cmd *DeployCmd) doBucket(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	if conf.Bucket != "" || (conf.Preset == "javascript" || conf.Preset == "typescript") {
		return nil
	}

	logger.FInfo(cmd.Io.Out, msg.ProjectNameMessage)
	nameBucket := replaceInvalidChars(conf.Name)
	err := client.CreateBucket(ctx, api.RequestBucket{
		BucketCreate: storage.BucketCreate{Name: nameBucket, EdgeAccess: storage.READ_WRITE}})
	if err != nil {
		// If the name is already in use, try 10 times with different names
		for i := 0; i < 10; i++ {
			nameB := fmt.Sprintf("%s-%s", nameBucket, utils.Timestamp())
			err := client.CreateBucket(ctx, api.RequestBucket{
				BucketCreate: storage.BucketCreate{Name: nameB, EdgeAccess: storage.READ_WRITE}})
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

	return cmd.WriteAzionJsonContent(conf)
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

// replaceInvalidChars Regular expression to find disallowed characters: "[^a-zA-Z0-9]+" replace invalid characters with -
func replaceInvalidChars(str string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\-]`)
	return re.ReplaceAllString(str, "")
}
