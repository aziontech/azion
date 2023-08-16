package deploy

import (
	"context"
	"os"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/schollz/progressbar/v3"
	"github.com/zRedShift/mimemagic"
	"go.uber.org/zap"
)

func (cmd *DeployCmd) uploadFiles(pathStatic string, versionID string) error {
	// Get total amount of files to display progress
	totalFiles := 0
	if err := cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Debug("Error while reading files to be uploaded", zap.Error(err))
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

	currentFile := 0
	bar := progressbar.NewOptions(
		totalFiles,
		progressbar.OptionSetDescription("Uploading files"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWriter(cmd.F.IOStreams.Out),
		progressbar.OptionClearOnFinish(),
	)
	if err := cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
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

			if err = clientUpload.Upload(context.Background(), versionID, fileString, mimeType.MediaType(), fileContent); err != nil {
				logger.Debug("Error while uploading file to storage api", zap.Error(err))
				return err
			}

			currentFile++
			bar.Set(currentFile)
		}
		return nil
	}); err != nil {
		logger.Debug("Error while reading files to be uploaded", zap.Error(err))
		return err
	}

	logger.FInfo(cmd.F.IOStreams.Out, msg.UploadSuccessful)

	return nil
}
