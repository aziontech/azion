package deploy

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/api/s3"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/schollz/progressbar/v3"
	"github.com/zRedShift/mimemagic"
	"go.uber.org/zap"
)

var (
	Jobs    chan contracts.FileOps
	Retries int64
)

func ReadAllFiles(pathStatic string, cmd *DeployCmd) ([]contracts.FileOps, error) {
	var listFiles []contracts.FileOps

	if err := cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Handle symlinks and resolve actual path type
		resolvedInfo, statErr := os.Stat(path)
		if statErr != nil {
			logger.Debug("Error resolving path", zap.String("path", path), zap.Error(statErr))
			return statErr
		}

		// Skip directories like node_modules
		if resolvedInfo.IsDir() && (strings.Contains(path, "node_modules")) {
			logger.Debug("Skipping directory: " + path)
			return filepath.SkipDir
		}

		// Process files
		if !resolvedInfo.IsDir() {
			fileContent, err := cmd.Open(path)
			if err != nil {
				logger.Debug("Error while trying to read file <"+path+"> about to be uploaded", zap.Error(err))
				return err
			}

			fileString := strings.TrimPrefix(path, pathStatic)
			mimeType, err := mimemagic.MatchFilePath(path, -1)
			if err != nil {
				logger.Debug("Error while matching file path", zap.Error(err))
				return err
			}

			fileOptions := contracts.FileOps{
				Path:        fileString,
				MimeType:    mimeType.MediaType(),
				FileContent: fileContent,
			}

			listFiles = append(listFiles, fileOptions)
		}

		return nil
	}); err != nil {
		logger.Debug("Error while reading files", zap.Error(err))
		return nil, err
	}

	return listFiles, nil
}

func uploadFiles(f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string, pathStatic, bucket string, cmd *DeployCmd, settings token.Settings) error {
	cfg, err := s3.New(settings.S3AccessKey, settings.S3SecretKey)
	if err != nil {
		return errors.New(msg.ErrorUnableSDKConfig + err.Error())
	}

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.UploadStart, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadStart)

	// create lots zip files
	allFiles, err := ReadAllFiles(pathStatic, cmd)
	if err != nil {
		return err
	}

	err = CreateZipsInBatches(allFiles)
	if err != nil {
		return err
	}

	// Ensure cleanup of temporary zip files after upload completion or on error
	defer func() {
		if cleanupErr := CleanupZipFiles(); cleanupErr != nil {
			logger.Debug("Failed to cleanup temporary zip files", zap.Error(cleanupErr))
		}
	}()

	listZip, err := ReadZip()
	if err != nil {
		return err
	}

	numFiles := len(listZip)

	noOfWorkers := 5
	var currentFile int64
	results := make(chan error, noOfWorkers)
	Jobs := make(chan contracts.FileOps, numFiles)

	// Create worker goroutines
	for i := 1; i <= noOfWorkers; i++ {
		go Worker(Jobs, results, &currentFile, cfg, bucket, conf.Prefix)
	}

	bar := progressbar.NewOptions(
		numFiles,
		progressbar.OptionSetDescription("Uploading files"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWriter(cmd.F.IOStreams.Out),
		progressbar.OptionClearOnFinish(),
	)

	if f.Silent {
		bar = nil
	}

	for _, f := range listZip {
		Jobs <- f
	}

	close(Jobs)

	// Check for errors from workers
	for a := 1; a <= numFiles; a++ {
		result := <-results
		if result != nil {
			return result
		}

		if bar != nil {
			err := bar.Set(int(currentFile))
			if err != nil {
				return err
			}
		}
	}

	// All jobs are processed, no more values will be sent on results:
	close(results)
	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.UploadSuccessful, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadSuccessful)

	return nil
}

// worker reads the range of jobs and uploads the file, if there is an error during upload, we returning it through the results channel
func Worker(jobs <-chan contracts.FileOps, results chan<- error, currentFile *int64, clientUpload aws.Config, bucket, prefix string) {
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

		if err := s3.UploadFile(context.Background(), clientUpload, &job, bucket, prefix); err != nil {
			logger.Debug("Error while worker tried to upload file: <"+job.Path+"> to storage api", zap.Error(err))
			for Retries < 20 {
				atomic.AddInt64(&Retries, 1)

				time.Sleep(time.Second * time.Duration(Retries))

				_, err := job.FileContent.Seek(0, 0)
				if err != nil {
					logger.Debug("An error occurred while seeking fileContent", zap.Error(err))
					break
				}

				logger.Debug("Retrying to upload the following file: <"+job.Path+"> to storage api", zap.Error(err))
				err = s3.UploadFile(context.Background(), clientUpload, &job, bucket, prefix)
				if err != nil {
					continue
				}
				break
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
