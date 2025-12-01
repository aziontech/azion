package deploy

import (
	"os"
	"path"
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
	Retries    int64
)

func (cmd *DeployCmd) uploadFiles(
	f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string, dir string) error {
	logger.Debug("Path to be uploaded: " + dir)

	clientUpload := storage.NewClient(cmd.F.HttpClient, cmd.F.Config.GetString("storage_url"), cmd.F.Config.GetString("token"))

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.UploadStart, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadStart)

	noOfWorkers := 5
	var currentFile int64

	// Collect all files in a single walk to avoid double traversal
	var fileOps []contracts.FileOps
	if err := cmd.FilepathWalk(dir, func(pathDir string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink != 0 {
			logger.Debug("Skipping symlink file", zap.Any("File name", pathDir))
			return nil
		}
		logger.Debug("Reading the following file", zap.Any("File name", pathDir))
		if !info.IsDir() {
			fileContent, err := cmd.Open(pathDir)
			if err != nil {
				logger.Debug("Error while trying to read file <"+pathDir+"> about to be uploaded", zap.Error(err))
				return err
			}

			fileString := strings.TrimPrefix(pathDir, path.Clean(dir))
			mimeType, err := mimemagic.MatchFilePath(pathDir, -1)
			if err != nil {
				logger.Debug("Error while matching file path", zap.Error(err))
				return err
			}
			fileOptions := contracts.FileOps{
				Path:        fileString,
				MimeType:    mimeType.MediaType(),
				FileContent: fileContent,
			}

			fileOps = append(fileOps, fileOptions)
		}
		return nil
	}); err != nil {
		logger.Debug("Error while reading files to be uploaded", zap.Error(err))
		return err
	}

	totalFiles := len(fileOps)
	if totalFiles == 0 {
		logger.FInfoFlags(cmd.F.IOStreams.Out, msg.UploadSuccessful, f.Format, f.Out)
		*msgs = append(*msgs, msg.UploadSuccessful)
		return nil
	}

	// Create channels and workers after we know the file count
	jobsChan := make(chan contracts.FileOps, totalFiles)
	results := make(chan error, noOfWorkers)

	// Create worker goroutines
	for i := 1; i <= noOfWorkers; i++ {
		go Worker(jobsChan, results, &currentFile, clientUpload, conf, conf.Bucket)
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

	// Queue all files for upload
	for _, fileOp := range fileOps {
		jobsChan <- fileOp
	}
	close(jobsChan)

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
