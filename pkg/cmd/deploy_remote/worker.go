package deploy

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aziontech/azion-cli/pkg/api/s3"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

const (
	maxRetries        = 20
	initialBackoff    = 1 * time.Second
	maxBackoff        = 30 * time.Second
	backoffMultiplier = 1.5
	rateLimitBackoff  = 5 * time.Second
)

// isRateLimitError checks if the error is due to rate limiting (HTTP 429)
func isRateLimitError(err error) bool {
	if err == nil {
		return false
	}
	// Check for various rate limit error patterns
	errStr := err.Error()
	return strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "too many requests") ||
		strings.Contains(errStr, "throttl")
}

// getRetryDelay calculates the delay for retry with exponential backoff
// For rate limit errors, uses a longer initial delay
func getRetryDelay(retryCount int, isRateLimit bool) time.Duration {
	delay := initialBackoff
	if isRateLimit {
		delay = rateLimitBackoff
	}

	for i := 0; i < retryCount; i++ {
		delay = time.Duration(float64(delay) * backoffMultiplier)
		if delay > maxBackoff {
			delay = maxBackoff
			break
		}
	}
	return delay
}

// Worker reads the range of jobs and uploads the file, if there is an error during upload, we return it through the results channel
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

		var retryCount int
		var lastErr error

		for retryCount < maxRetries {
			err := s3.UploadFile(context.Background(), cfg, &job, bucket, prefix)
			if err == nil {
				break
			}

			lastErr = err
			retryCount++

			isRateLimit := isRateLimitError(err)

			if isRateLimit {
				logger.Debug("Rate limit detected, applying backoff for file: "+job.Path,
					zap.Int("retry", retryCount),
					zap.Error(err))
			} else {
				logger.Debug("Error while worker tried to upload file: <"+job.Path+"> to storage api",
					zap.Int("retry", retryCount),
					zap.Error(err))
			}

			if retryCount >= maxRetries {
				logger.Debug("Max retries reached for file: "+job.Path, zap.Int("total_retries", retryCount))
				break
			}

			delay := getRetryDelay(retryCount, isRateLimit)
			logger.Debug("Waiting before retry",
				zap.String("file", job.Path),
				zap.Duration("delay", delay),
				zap.Int("retry", retryCount))

			time.Sleep(delay)

			_, seekErr := job.FileContent.Seek(0, 0)
			if seekErr != nil {
				logger.Debug("An error occurred while seeking fileContent", zap.Error(seekErr))
				results <- seekErr
				return
			}
		}

		if retryCount >= maxRetries && lastErr != nil {
			logger.Debug("Upload failed after max retries",
				zap.String("file", job.Path),
				zap.Int("total_retries", retryCount),
				zap.Error(lastErr))
			results <- lastErr
			return
		}

		atomic.AddInt64(currentFile, 1)
		results <- nil
	}
}
