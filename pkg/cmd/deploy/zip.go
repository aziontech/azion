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

// CreateZipsFromFileInfos creates ZIP files in batches from FileInfo slices
// This function opens files in batches to avoid "too many open files" error
func CreateZipsFromFileInfos(fileInfos []FileInfo) ([]string, error) {
	const maxBatchSize = 1 * 1024 * 1024 // 1MB in bytes
	const maxOpenFiles = 100             // Maximum files to keep open at once
	var createdZips []string
	tempDir := os.TempDir()

	batchNumber := 1

	for i := 0; i < len(fileInfos); {
		var currentBatch []FileInfo
		var currentSize int64 = 0

		// Build a batch respecting both size and file count limits
		for j := i; j < len(fileInfos) && len(currentBatch) < maxOpenFiles; j++ {
			fileInfo := fileInfos[j]
			if currentSize+fileInfo.Size > maxBatchSize && len(currentBatch) > 0 {
				break
			}
			currentBatch = append(currentBatch, fileInfo)
			currentSize += fileInfo.Size
		}

		if len(currentBatch) == 0 {
			break
		}

		// Create ZIP for the current batch
		zipPath, err := createZipFromInfos(currentBatch, tempDir, batchNumber)
		if err != nil {
			return nil, err
		}
		createdZips = append(createdZips, zipPath)

		batchNumber++
		i += len(currentBatch)
	}

	return createdZips, nil
}

// createZipFromInfos creates a ZIP file from FileInfo entries, opening and closing files as needed
func createZipFromInfos(batch []FileInfo, destDir string, batchNumber int) (string, error) {
	zipFileName := fmt.Sprintf("batch_%d.zip", batchNumber)
	zipFilePath := filepath.Join(destDir, zipFileName)
	logger.Debug("Creating ZIP file", zap.String("path", zipFilePath), zap.Int("files", len(batch)))

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf(msg.ErrorCreateZip, zipFilePath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, fileInfo := range batch {
		// Open file, copy content, and close immediately to limit open file handles
		file, err := os.Open(fileInfo.AbsolutePath)
		if err != nil {
			return "", fmt.Errorf("error opening file %s: %v", fileInfo.AbsolutePath, err)
		}

		// Add a file to the ZIP with the relative path
		writer, err := zipWriter.Create(fileInfo.Path)
		if err != nil {
			file.Close()
			return "", fmt.Errorf(msg.ErrorCreateZip, fileInfo.Path, err)
		}

		_, err = io.Copy(writer, file)
		file.Close() // Close immediately after copying

		if err != nil {
			return "", fmt.Errorf(msg.ErrorCopyContentFile, fileInfo.Path, err)
		}
	}

	logger.Debug("Created ZIP file successfully", zap.String("path", zipFilePath), zap.Int("files", len(batch)))
	return zipFilePath, nil
}

func CreateZipsInBatches(files []contracts.FileOps) ([]string, error) {
	const maxBatchSize = 1 * 1024 * 1024 // 1MB em bytes
	var currentBatch []contracts.FileOps
	var currentSize int64 = 0
	var createdZips []string
	tempDir := os.TempDir()

	batchNumber := 1

	for _, fileOp := range files {
		info, err := fileOp.FileContent.Stat()
		if err != nil {
			return nil, fmt.Errorf("error getting info from file %s: %v", fileOp.Path, err)
		}
		fileSize := info.Size()

		// Check if adding this file exceeds the maximum batch size
		if currentSize+fileSize > maxBatchSize && len(currentBatch) > 0 {
			// Create ZIP for the current batch
			zipPath, err := createZip(currentBatch, tempDir, batchNumber)
			if err != nil {
				return nil, err
			}
			createdZips = append(createdZips, zipPath)

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
		zipPath, err := createZip(currentBatch, tempDir, batchNumber)
		if err != nil {
			return nil, err
		}
		createdZips = append(createdZips, zipPath)
	}

	return createdZips, nil
}

func createZip(batch []contracts.FileOps, destDir string, batchNumber int) (string, error) {
	zipFileName := fmt.Sprintf("batch_%d.zip", batchNumber)
	zipFilePath := filepath.Join(destDir, zipFileName)
	logger.Debug("Creating ZIP file", zap.String("path", zipFilePath), zap.Int("batch", len(batch)))

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return "", fmt.Errorf(msg.ErrorCreateZip, zipFilePath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, fileOp := range batch {
		relPath, err := filepath.Rel(destDir, fileOp.Path)
		if err != nil {
			return "", fmt.Errorf(msg.ErrorRelPath, fileOp.Path, err)
		}

		// Replace directory separators for ZIP compatibility
		zipPath := filepath.ToSlash(relPath)

		// Add a file to the ZIP with the relative path
		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return "", fmt.Errorf(msg.ErrorCreateZip, zipPath, err)
		}

		if _, err := fileOp.FileContent.Seek(0, io.SeekStart); err != nil {
			return "", fmt.Errorf(msg.ErrorResetPointFile, fileOp.Path, err)
		}

		_, err = io.Copy(writer, fileOp.FileContent)
		if err != nil {
			return "", fmt.Errorf(msg.ErrorCopyContentFile, fileOp.Path, err)
		}
	}

	logger.Debug("Create ZIP file successfully", zap.String("path", zipFilePath), zap.Int("batch", len(batch)))
	return zipFilePath, nil
}

func isBatchZipFile(name string) bool {
	return strings.ToLower(filepath.Ext(name)) == ".zip" && strings.HasPrefix(name, "batch")
}

func ReadZip() ([]contracts.FileOps, error) {
	var listZIP []contracts.FileOps

	tempDir := os.TempDir()

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return []contracts.FileOps{}, err
	}

	for _, f := range files {
		if isBatchZipFile(f.Name()) {
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
			listZIP = append(listZIP, fileOptions)

		}

	}

	return listZIP, nil
}

func CleanupZipFiles() error {
	logger.Debug("Running cleanup for batch zip files")
	tempDir := os.TempDir()

	files, err := os.ReadDir(tempDir)
	if err != nil {
		logger.Debug("Error reading temp directory for cleanup", zap.Error(err))
		return err
	}

	var cleanupErrors []error
	for _, f := range files {
		if isBatchZipFile(f.Name()) {
			zipPath := filepath.Join(tempDir, f.Name())
			logger.Debug("Removing temporary ZIP file", zap.String("path", zipPath))

			if err := os.Remove(zipPath); err != nil {
				logger.Debug("Error removing ZIP file", zap.String("path", zipPath), zap.Error(err))
				cleanupErrors = append(cleanupErrors, err)
			}
		}
	}

	if len(cleanupErrors) > 0 {
		logger.Debug("Some ZIP files could not be removed during cleanup", zap.Int("count", len(cleanupErrors)))
		return cleanupErrors[0]
	}

	return nil
}
