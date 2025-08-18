package deploy

import (
	"context"
	"sync/atomic"

	"github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

// worker reads the range of jobs and uploads the file, if there is an error during upload, we returning it through the results channel
func Worker(jobs <-chan contracts.FileOps, results chan<- error, currentFile *int64, clientUpload *storage.Client, conf *contracts.AzionApplicationOptions, bucket string) {
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
			return
		}

		if err := clientUpload.Upload(context.Background(), &job, conf, bucket); err != nil {
			logger.Debug("Error while worker tried to upload file: <"+job.Path+"> to storage api", zap.Error(err))
			for Retries < 5 {
				atomic.AddInt64(&Retries, 1)
				_, err := job.FileContent.Seek(0, 0)
				if err != nil {
					logger.Debug("An error occurred while seeking fileContent", zap.Error(err))
					break
				}

				logger.Debug("Retrying to upload the following file: <"+job.Path+"> to storage api", zap.Error(err))
				err = clientUpload.Upload(context.Background(), &job, conf, bucket)
				if err != nil {
					continue
				}
				break
			}

			if Retries >= 5 {
				logger.Debug("There have been 5 retries already, quitting upload")
				results <- err
				return
			}
		}

		atomic.AddInt64(currentFile, 1)
		results <- nil
	}
}
