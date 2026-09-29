package deploy

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

// CreateZipsInBatches Function to split files into batches of up to 1MB and
// create ZIPs keeping the directory structure
func CreateZipsInBatches(files []contracts.FileOps) error {
	const maxBatchSize = 1 * 1024 * 1024 // 1MB em bytes
	var currentBatch []contracts.FileOps
	var currentSize int64 = 0
	tempDir := os.TempDir()

	batchNumber := 1

	for _, fileOp := range files {
		info, err := fileOp.FileContent.Stat()
		if err != nil {
			return fmt.Errorf("error getting info from file %s: %v", fileOp.Path, err)
		}
		fileSize := info.Size()

		// Check if adding this file exceeds the maximum batch size
		if currentSize+fileSize > maxBatchSize && len(currentBatch) > 0 {
			// Create ZIP for the current batch
			if err := createZip(currentBatch, tempDir, batchNumber); err != nil {
				return err
			}

			batchNumber++

			// reset batch
			currentBatch = []contracts.FileOps{}
			currentSize = 0
		}

		// Add the file to the current batch
		currentBatch = append(currentBatch, fileOp)
		currentSize += fileSize
	}

	// Create ZIP for any remaining files
	if len(currentBatch) > 0 {
		if err := createZip(currentBatch, tempDir, batchNumber); err != nil {
			return err
		}
	}

	return nil
}

// createZip Helper function to create a ZIP file from a batch of FileOps while maintaining the directory structure
func createZip(batch []contracts.FileOps, destDir string, batchNumber int) error {
	zipFileName := fmt.Sprintf("batch_%d.zip", batchNumber)
	zipFilePath := filepath.Join(destDir, zipFileName)

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf(msg.ErrorCreateZip, zipFilePath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, fileOp := range batch {
		relPath, err := filepath.Rel(destDir, fileOp.Path)
		if err != nil {
			return fmt.Errorf(msg.ErrorRelPath, fileOp.Path, err)
		}

		// Replace directory separators for ZIP compatibility
		zipPath := filepath.ToSlash(relPath)

		// Add a file to the ZIP with the relative path
		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return fmt.Errorf(msg.ErrorCreateZip, zipPath, err)
		}

		if _, err := fileOp.FileContent.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf(msg.ErrorResetPointFile, fileOp.Path, err)
		}

		_, err = io.Copy(writer, fileOp.FileContent)
		if err != nil {
			return fmt.Errorf(msg.ErrorCopyContentFile, fileOp.Path, err)
		}
	}

	logger.Debug("Create ZIP file", zap.String("path", zipFilePath), zap.Int("batch", len(batch)))
	return nil
}

// ReadZip reads the ZIP files in the pathStatic directory that start with
// prefix and end with .zip.  Returns a list of contracts.FileOps or an error.
func ReadZip() ([]contracts.FileOps, error) {
	var listZIP []contracts.FileOps

	tempDir := os.TempDir()

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return []contracts.FileOps{}, err
	}

	for _, f := range files {
		if strings.ToLower(filepath.Ext(f.Name())) == ".zip" &&
			strings.HasPrefix(f.Name(), "batch") {
			pathZIP := filepath.Join(tempDir, f.Name())
			logger.Debug("Processing ZIP file " + pathZIP)

			f, err := os.Open(pathZIP)
			if err != nil {
				logger.Debug("Error opening ZIP "+pathZIP, zap.Error(err))
				continue
			}

			fileString := strings.TrimPrefix(pathZIP, tempDir)
			fileOptions := contracts.FileOps{
				Path:        fileString,
				FileContent: f,
			}
			//
			listZIP = append(listZIP, fileOptions)

		}

	}

	return listZIP, nil
}
