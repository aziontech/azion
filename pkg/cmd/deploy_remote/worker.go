package deploy

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aziontech/azion-cli/pkg/api/s3"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

// worker reads the range of jobs and uploads the file, if there is an error during upload, we returning it through the results channel
func Worker(jobs <-chan contracts.FileOps, results chan<- error, currentFile *int64, cfg aws.Config, bucket, prefix string) {
	for job := range jobs {
		// Once ENG-27343 is completed, we might be able to remove this piece of code
		fileInfo, err := job.FileContent.Stat()
		if err != nil {
			logger.Debug("Error while worker tried to read file stats", zap.Error(err))
			logger.Debug("File that caused the error: " + job.Path)
			results <- err
			return
		}

		// Check if the file size is zero
		if fileInfo.Size() == 0 {
			logger.Debug("\nSkipping upload of empty file: " + job.Path)
			results <- nil
			atomic.AddInt64(currentFile, 1)
			continue
		}

		if err := s3.UploadFile(context.Background(), cfg, &job, bucket, prefix); err != nil {
			logger.Debug("Error while worker tried to upload file: <"+job.Path+"> to storage api", zap.Error(err))

			for Retries < 20 {
				atomic.AddInt64(&Retries, 1)

				time.Sleep(time.Second * time.Duration(Retries))

				_, seekErr := job.FileContent.Seek(0, 0)
				if seekErr != nil {
					logger.Debug("An error occurred while seeking fileContent", zap.Error(seekErr))
					break
				}

				logger.Debug("Retrying to upload the following file: <"+job.Path+"> to storage api", zap.Error(err))
				err = s3.UploadFile(context.Background(), cfg, &job, bucket, prefix)
				if err == nil {
					break
				}
			}

			if Retries >= 20 {
				logger.Debug("There have been 20 retries already, quitting upload")
				results <- err
				return
			}
		}

		atomic.AddInt64(currentFile, 1)
		results <- nil
	}
}
