package deploy

import (
	"os"
	"path/filepath"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/api/storage"
	remote "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/schollz/progressbar/v3"
	"github.com/zRedShift/mimemagic"
	"go.uber.org/zap"
)

var (
	Jobs    chan contracts.FileOps
	Retries int64
)

func uploadFiles(f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string, pathStatic, bucket string, cmd *DeployCmd) error {
	// Get total amount of files to display progress
	totalFiles := 0
	if err := cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Debug("Error while reading files to be uploaded", zap.Error(err))
			logger.Debug("File that caused the error: " + pathStatic)
			return err
		}

		// Skip node_modules and .edge directories
		if info.IsDir() && (strings.Contains(path, "node_modules") || strings.Contains(path, ".edge")) {
			logger.Debug("Skipping directory: " + path)
			return filepath.SkipDir
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

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.UploadStart, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadStart)

	noOfWorkers := 5
	var currentFile int64
	Jobs := make(chan contracts.FileOps, totalFiles)
	results := make(chan error, noOfWorkers)

	// Create worker goroutines
	for i := 1; i <= noOfWorkers; i++ {
		go remote.Worker(Jobs, results, &currentFile, clientUpload, conf, bucket)
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

	if err := cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip node_modules and .edge directories
		if info.IsDir() && (strings.Contains(path, "node_modules") || strings.Contains(path, ".edge")) {
			logger.Debug("Skipping directory: " + path)
			return filepath.SkipDir
		}

		if !info.IsDir() {
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
	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.UploadSuccessful, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadSuccessful)

	return nil
}
