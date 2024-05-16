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
func worker(jobs <-chan contracts.FileOps, results chan<- error, currentFile *int64, clientUpload *storage.Client, conf *contracts.AzionApplicationOptions) {
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
			return
		}

		// for

		if err := clientUpload.Upload(context.Background(), &job, conf); err != nil {
			logger.Debug("Error while worker tried to upload file: <"+job.Path+"> to storage api", zap.Error(err))
			if Retries < 5 {
				logger.Debug("Retrying to upload the following file: <"+job.Path+"> to storage api", zap.Error(err))
				atomic.AddInt64(&Retries, 1)
				Jobs <- job
				results <- nil
				return
			}
			logger.Debug("There have been 5 retries already, quitting upload")
			results <- err
			return
		}

		atomic.AddInt64(currentFile, 1)
		results <- nil
	}
}
