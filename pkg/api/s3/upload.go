package s3

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"

	"go.uber.org/zap"
)

const (
	region   = "us-east"
	endpoint = "https://s3.us-east-005.stageazionstorage.net"
)

type CustomEndpointResolver struct {
	URL           string
	SigningRegion string
}

// ResolveEndpoint is the method that defines the custom endpoint
func (e *CustomEndpointResolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) { // nolint
	return aws.Endpoint{ //nolint
		URL:           e.URL,
		SigningRegion: e.SigningRegion,
	}, nil
}

func New(s3AccessKey, s3SecretKey string) (aws.Config, error) {
	endpointResolver := &CustomEndpointResolver{
		URL:           endpoint,
		SigningRegion: region,
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3AccessKey, s3SecretKey, "")),
		config.WithEndpointResolver(endpointResolver), // nolint
	)
	if err != nil {
		return aws.Config{}, errors.New(msg.ErrorUnableSDKConfig + err.Error())

	}
	return cfg, nil
}

func UploadFile(ctx context.Context, cfg aws.Config, fileOps *contracts.FileOps, bucketName, prefix string) error {
	file := fileOps.Path
	if prefix != "" {
		file = path.Join(prefix, fileOps.Path)
	}

	s3Client := s3.NewFromConfig(cfg)

	logger.Debug("Object_key: " + file)
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(file), // Name of the file in the bucket
		Body:   fileOps.FileContent,
	}

	_, err := s3Client.PutObject(ctx, uploadInput)
	if err != nil {
		logger.Debug("Error while uploading file <"+fileOps.Path+"> to storage api", zap.Error(err))
		return fmt.Errorf(msg.ErrorUploadFileBucket, bucketName, err)
	}

	return nil
}
