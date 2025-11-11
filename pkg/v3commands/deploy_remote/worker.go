package deploy

import (
	"context"
	"sync/atomic"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/v3api/storage"
	"go.uber.org/zap"
)

// worker reads the range of jobs and uploads the file, if there is an error during upload, we returning it through the results channel
func Worker(jobs <-chan contracts.FileOps, results chan<- error, currentFile *int64, clientUpload *storage.Client, conf *contracts.AzionApplicationOptionsV3, bucket string) {
	for job := range jobs {
		// Once ENG-27343 is completed, we might be able to remove this piece of code
		fileInfo, err := job.FileContent.Stat()
		if err != nil {
			logger.Debug("Error while worker tried to read file stats", zap.Error(err))
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

		if err := clientUpload.Upload(context.Background(), &job, conf, bucket); err != nil {
			logger.Debug("Error while worker tried to upload file: <"+job.Path+"> to storage api", zap.Error(err))

			fileRetries := 0
			maxRetries := 5

			for fileRetries < maxRetries {
				fileRetries++
				atomic.AddInt64(&Retries, 1)

				_, seekErr := job.FileContent.Seek(0, 0)
				if seekErr != nil {
					logger.Debug("An error occurred while seeking fileContent", zap.Error(seekErr))
					break
				}

				logger.Debug("Retrying to upload file", zap.Int("attempt", fileRetries), zap.Int("maxRetries", maxRetries), zap.String("path", job.Path))
				err = clientUpload.Upload(context.Background(), &job, conf, bucket)
				if err == nil {
					break
				}
			}

			if fileRetries >= maxRetries {
				logger.Debug("Failed to upload file after retries", zap.Int("maxRetries", maxRetries), zap.String("path", job.Path))
				results <- err
				return
			}
		}

		atomic.AddInt64(currentFile, 1)
		results <- nil
	}
}
