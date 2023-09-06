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
func worker(jobs <-chan contracts.FileOps, results chan<- error, currentFile *int64, clientUpload *storage.Client) {

	for job := range jobs {
		if err := clientUpload.Upload(context.Background(), &job); err != nil {
			logger.Debug("Error while worker tried to upload file to storage api", zap.Error(err))
			results <- err
			return
		}
		atomic.AddInt64(currentFile, 1)
		results <- nil
	}
}
