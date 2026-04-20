package deploy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/storage-api"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	api "github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
)

func (cmd *DeployCmd) doBucket(
	client *api.Client,
	ctx context.Context,
	conf *contracts.AzionApplicationOptions,
	msgs *[]string,
	manifestStorage []contracts.StorageManifest) error {

	doBucketStart := time.Now()
	if conf.Bucket != "" {
		return nil
	}

	logger.FInfoFlags(cmd.Io.Out, msg.ProjectNameMessage, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msg.ProjectNameMessage)
	nameBucket := utils.ReplaceInvalidCharsBucket(conf.Name)

	bucketAccess := "read_only"
	if WriteBucket {
		bucketAccess = "read_write"
	} else if manifestStorage[0].WorkloadsAccess != "" {
		bucketAccess = manifestStorage[0].WorkloadsAccess
	}
	apiCallStart := time.Now()
	err := client.CreateBucket(ctx, api.RequestBucket{BucketCreateRequest: sdk.BucketCreateRequest{Name: nameBucket, WorkloadsAccess: bucketAccess}})
	if err != nil {
		// If the name is already in use, try 10 times with different names
		for i := 0; i < 10; i++ {
			nameB := fmt.Sprintf("%s-%s", nameBucket, utils.Timestamp())
			msgf := fmt.Sprintf(msg.NameInUseBucket, nameB)
			logger.FInfoFlags(cmd.Io.Out, msgf, cmd.F.Format, cmd.F.Out)
			*msgs = append(*msgs, msgf)
			apiCallStart = time.Now()
			err := client.CreateBucket(ctx, api.RequestBucket{
				BucketCreateRequest: sdk.BucketCreateRequest{Name: nameB, WorkloadsAccess: bucketAccess}})
			if err != nil {
				GlobalTimingSummary.AddAPICallTime("Bucket.Create (retry)", time.Since(apiCallStart))
				if errors.Is(err, utils.ErrorNameInUse) && i < 9 {
					continue
				}
				return err
			}
			GlobalTimingSummary.AddAPICallTime("Bucket.Create", time.Since(apiCallStart))
			conf.Bucket = nameB
			break
		}
	} else {
		GlobalTimingSummary.AddAPICallTime("Bucket.Create", time.Since(apiCallStart))
		conf.Bucket = nameBucket
	}

	GlobalTimingSummary.BucketCreateTime = time.Since(doBucketStart)

	msgf := fmt.Sprintf(msg.BucketSuccessful, conf.Bucket)
	logger.FInfoFlags(cmd.Io.Out, msgf, cmd.F.Format, cmd.F.Out)
	*msgs = append(*msgs, msgf)
	return cmd.WriteAzionJsonContent(conf, ProjectConf)
}

// CreateBucketCredentials creates S3 credentials for a specific bucket and saves them to the credentials file
func CreateBucketCredentials(ctx context.Context, bucketName string, f *cmdutil.Factory, subdir string) (token.S3Credentials, error) {
	logger.Debug("Creating S3 credentials for bucket")

	storageClient := api.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))

	// Get the current time
	now := time.Now()

	// Add one year to the current time
	oneYearLater := now.AddDate(1, 0, 0)

	request := api.RequestCredentials{}
	request.Name = bucketName
	request.Capabilities = []string{"listAllBucketNames", "listBuckets", "listFiles", "readFiles", "writeFiles", "deleteFiles"}
	request.Buckets = []string{bucketName}
	request.ExpirationDate = &oneYearLater

	creds, err := storageClient.CreateCredentials(ctx, request)
	if err != nil {
		return token.S3Credentials{}, fmt.Errorf("failed to create credentials for bucket %s: %w", bucketName, err)
	}

	s3Creds := token.S3Credentials{
		S3AccessKey: creds.Data.GetAccessKey(),
		S3SecretKey: creds.Data.GetSecretKey(),
	}

	return s3Creds, nil
}

// GetOrCreateCredentials retrieves existing credentials for a bucket or creates new ones if they don't exist
func (cmd *DeployCmd) GetOrCreateCredentials(ctx context.Context, bucketName string, profile string) (token.S3Credentials, error) {
	// First, check if credentials already exist for this bucket
	creds, exists, err := cmd.GetCredentialsForBucket(profile, bucketName)
	if err != nil {
		return token.S3Credentials{}, fmt.Errorf("failed to read credentials: %w", err)
	}

	if exists {
		logger.Debug("Found existing credentials for bucket")
		return creds, nil
	}

	// Credentials don't exist, create them
	logger.Debug("Creating new credentials for bucket")
	creds, err = CreateBucketCredentials(ctx, bucketName, cmd.F, profile)
	if err != nil {
		return token.S3Credentials{}, err
	}

	// Save the credentials to the credentials file
	err = cmd.SaveCredentialsForBucket(profile, bucketName, creds)
	if err != nil {
		return token.S3Credentials{}, fmt.Errorf("failed to save credentials: %w", err)
	}

	return creds, nil
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
