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
	for job := range jobs {
		var attempt int
		var lastError error

		defer job.FileContent.Close()

		for attempt = 1; attempt <= 5; attempt++ {
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
				_, err := job.FileContent.Seek(0, 0)
				if err != nil {
					logger.Debug("An error occurred while seeking fileContent", zap.Error(err))
					break
				}

				logger.Debug("Erro ao tentar fazer o upload do arquivo", zap.Error(err))
				lastError = err
				if attempt < 5 {
					continue
				} else {
					break
				}
			}

			atomic.AddInt64(currentFile, 1)
			results <- nil
			break
		}

		if attempt > 5 || lastError != nil {
			logger.Debug("There have been 5 retries already, quitting upload")
			results <- lastError
			return
		}
	}
}
