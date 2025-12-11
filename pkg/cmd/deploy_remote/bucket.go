package deploy

import (
	"context"
	"errors"
	"fmt"
	"os"

	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/storage-api"

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
	msgs *[]string,
	manifestStorage []contracts.StorageManifest) error {
	if conf.Bucket != "" {
		return nil
	}

	logger.FInfoFlags(cmd.Io.Out, msg.ProjectNameMessage, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msg.ProjectNameMessage)
	nameBucket := utils.ReplaceInvalidCharsBucket(conf.Name)

	bucketAccess := "read_only"
	if WriteBucket {
		bucketAccess = "read_write"
	} else if manifestStorage[0].EdgeAccess != "" {
		bucketAccess = manifestStorage[0].EdgeAccess
	}
	err := client.CreateBucket(ctx, api.RequestBucket{BucketCreateRequest: sdk.BucketCreateRequest{Name: nameBucket, EdgeAccess: bucketAccess}})
	if err != nil {
		// If the name is already in use, try 10 times with different names
		for i := 0; i < 10; i++ {
			nameB := fmt.Sprintf("%s-%s", nameBucket, utils.Timestamp())
			msgf := fmt.Sprintf(msg.NameInUseBucket, nameB)
			logger.FInfoFlags(cmd.Io.Out, msgf, cmd.F.Format, cmd.F.Out)
			*msgs = append(*msgs, msgf)
			err := client.CreateBucket(ctx, api.RequestBucket{
				BucketCreateRequest: sdk.BucketCreateRequest{Name: nameB, EdgeAccess: bucketAccess}})
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
