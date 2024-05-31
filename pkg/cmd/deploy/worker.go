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
func worker(jobs <-chan contracts.FileOps, results chan<- error, currentFile *int64, client *storage.Client, conf *contracts.AzionApplicationOptions) {
	var lastError error
	for job := range jobs {
		for Retries <= 5 {
			fileInfo, err := job.FileContent.Stat()
			if err != nil {
				logger.Debug("Error while worker tried to read file stats", zap.Error(err))
				lastError = err
				break
			}

			if fileInfo.Size() == 0 {
				logger.Debug("skipping upload of empty file", zap.String("path", job.Path))
				results <- nil
				atomic.AddInt64(currentFile, 1)
				break
			}

			if err := client.Upload(context.Background(), &job, conf); err != nil {
				logger.Debug("Error while worker tried to upload file: <"+job.Path+"> to storage api", zap.Error(err))
				atomic.AddInt64(&Retries, 1)
				_, err := job.FileContent.Seek(0, 0)
				if err != nil {
					logger.Debug("An error occurred while seeking fileContent", zap.Error(err))
					break
				}

				if Retries < 5 {
					logger.Debug("Retrying to upload the following file: <"+job.Path+"> to storage api", zap.Error(err))
					continue
				} else {
					break
				}
			}

			atomic.AddInt64(currentFile, 1)
			results <- nil
			break
		}

		if Retries > 5 || lastError != nil {
			results <- lastError
			return
		}
	}
}
