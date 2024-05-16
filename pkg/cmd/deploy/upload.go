package deploy

import (
	"os"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/schollz/progressbar/v3"
	"github.com/zRedShift/mimemagic"
	"go.uber.org/zap"
)

var (
	PathStatic = ".edge/storage"
	Jobs       chan contracts.FileOps
	Retries    int64
)

func (cmd *DeployCmd) uploadFiles(f *cmdutil.Factory, conf *contracts.AzionApplicationOptions) error {
	// Get total amount of files to display progress
	totalFiles := 0
	if err := cmd.FilepathWalk(PathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Debug("Error while reading files to be uploaded", zap.Error(err))
			logger.Debug("File that caused the error: " + PathStatic)
			return err
		}
		if !info.IsDir() {
			totalFiles++
		}
		return nil
	}); err != nil {
		logger.Debug("Error while reading files to be uploaded", zap.Error(err))
		return err
	}

	clientUpload := storage.NewClient(cmd.F.HttpClient, cmd.F.Config.GetString("storage_url"), cmd.F.Config.GetString("token"))

	logger.FInfo(cmd.F.IOStreams.Out, msg.UploadStart)

	noOfWorkers := 5
	var currentFile int64
	Jobs := make(chan contracts.FileOps, totalFiles)
	results := make(chan error, noOfWorkers)

	// Create worker goroutines
	for i := 1; i <= noOfWorkers; i++ {
		go worker(Jobs, results, &currentFile, clientUpload, conf)
	}

	bar := progressbar.NewOptions(
		totalFiles,
		progressbar.OptionSetDescription("Uploading files"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWriter(cmd.F.IOStreams.Out),
		progressbar.OptionClearOnFinish(),
	)

	if f.Silent {
		bar = nil
	}

	if err := cmd.FilepathWalk(PathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileContent, err := cmd.Open(path)
			if err != nil {
				logger.Debug("Error while trying to read file <"+path+"> about to be uploaded", zap.Error(err))
				return err
			}

			fileString := strings.TrimPrefix(path, PathStatic)
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

			Jobs <- fileOptions
		}
		return nil
	}); err != nil {
		logger.Debug("Error while reading files to be uploaded", zap.Error(err))
		return err
	}
	close(Jobs)

	// Check for errors from workers
	for a := 1; a <= totalFiles; a++ {
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
	logger.FInfo(cmd.F.IOStreams.Out, msg.UploadSuccessful)

	return nil
}
