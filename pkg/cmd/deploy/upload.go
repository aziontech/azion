package deploy

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/schollz/progressbar/v3"
)

const maxZipSize = 1 * 1024 * 1024 // 1MB

var (
	region     = "us-east"
	endpoint   = "https://s3.us-east-005.azionstorage.net"
	numWorkers = 5
)

type CustomEndpointResolver struct {
	URL           string
	SigningRegion string
}

// ResolveEndpoint is the method that defines the custom endpoint
func (e *CustomEndpointResolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) { // nolint
	return aws.Endpoint{ //nolint
		URL:           e.URL,
		SigningRegion: e.SigningRegion,
	}, nil
}

// getFilesFromDir func that reads all the files in a directory and returns a slice with the full paths
func getFilesFromDir(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// addFileToZip func to add an individual file to the zip archive, preserving the relative path
func addFileToZip(zipWriter *zip.Writer, filePath, baseDir string) (os.FileInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Create an entry in the zip with the relative path
	relativePath := strings.TrimPrefix(filePath, baseDir)
	writer, err := zipWriter.Create(relativePath)
	if err != nil {
		return nil, err
	}

	// Copy the contents of the file to the zip
	_, err = io.Copy(writer, file)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

// uploadFile func to upload a file to the bucket
func uploadFile(ctx context.Context, cfg aws.Config, bucketName, filePath string) error {
	s3Client := s3.NewFromConfig(cfg)

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	fileName := filepath.Base(filePath)

	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName), // Name of the file in the bucket
		Body:   file,
	}

	_, err = s3Client.PutObject(ctx, uploadInput)
	if err != nil {
		return fmt.Errorf("failed to upload file to bucket %s: %w", bucketName, err)
	}

	return nil
}

// Worker responsible for zipping and uploading files
func worker(ctx context.Context, cfg aws.Config, filesChan <-chan string,
	wg *sync.WaitGroup, baseDir string, workerID int, atomicCounter *int32,
	bucketName string, progressBar *progressbar.ProgressBar) {
	defer wg.Done()

	var (
		currentSize int64
		filesToZip  []string
	)

	zipFileName := fmt.Sprintf("project_part_%d.zip", workerID)

	for file := range filesChan {
		filesToZip = append(filesToZip, file)

		// Get the file size and add it to the current size
		fileInfo, err := os.Stat(file)
		if err != nil {
			log.Printf("Failed to get file info for %s: %v", file, err)
			continue
		}
		currentSize += fileInfo.Size()

		// If the current size exceeds the limit (1MB), zip and upload
		if currentSize >= maxZipSize {
			err = zipAndUpload(ctx, cfg, filesToZip, zipFileName, baseDir,
				atomicCounter, bucketName, progressBar)
			if err != nil {
				log.Printf("Failed to zip and upload files: %v", err)
			}

			// Reset the size counter and clear the file list
			currentSize = 0
			filesToZip = nil
		}
	}

	// If there are still files at the end, zip them up and upload them
	if len(filesToZip) > 0 {
		err := zipAndUpload(ctx, cfg, filesToZip, zipFileName, baseDir, atomicCounter, bucketName, progressBar)
		if err != nil {
			log.Printf("Failed to zip and upload files: %v", err)
		}
	}
}

func zipAndUpload(ctx context.Context, cfg aws.Config, filesToZip []string,
	zipFileName, baseDir string, atomicCounter *int32,
	bucketName string, progressBar *progressbar.ProgressBar) error {

	// Increment atomic counter to generate unique zip file name
	counter := atomic.AddInt32(atomicCounter, 1)
	zipFileName = fmt.Sprintf("%s_%d.zip", strings.TrimSuffix(zipFileName, ".zip"), counter)

	zipFilePath := filepath.Join(baseDir, zipFileName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file %s: %w", zipFileName, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files to the zip
	for _, file := range filesToZip {
		_, err = addFileToZip(zipWriter, file, baseDir)
		if err != nil {
			return fmt.Errorf("failed to add file to zip %s: %w", file, err)
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return fmt.Errorf("failed to close zip file %s: %w", zipFileName, err)
	}

	// Check if the zip file was created
	if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
		return fmt.Errorf("zip file %s does not exist", zipFileName)
	}

	err = uploadFile(ctx, cfg, bucketName, zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to upload zip file %s: %w", zipFileName, err)
	}

	if progressBar != nil {
		err := progressBar.Set(int(len(filesToZip)))
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadFiles(f *cmdutil.Factory, msgs *[]string, pathStatic string, settings token.Settings) error {
	files, err := getFilesFromDir(pathStatic)
	if err != nil {
		return err
	}

	endpointResolver := &CustomEndpointResolver{
		URL:           endpoint,
		SigningRegion: region,
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(settings.S3AccessKey, settings.S3SecreKey, "")),
		config.WithEndpointResolver(endpointResolver), // nolint
	)
	if err != nil {
		return errors.New("unable to load SDK config, " + err.Error())
	}

	ctx := context.TODO()

	bar := progressbar.NewOptions(
		len(files),
		progressbar.OptionSetDescription("Uploading files"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWriter(f.IOStreams.Out),
		progressbar.OptionClearOnFinish(),
	)

	if f.Silent {
		bar = nil
	}

	logger.FInfoFlags(f.IOStreams.Out, msg.UploadStart, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadStart)

	filesChan := make(chan string)

	var wg sync.WaitGroup

	// Atomic counter to guarantee unique zip file names
	var atomicCounter int32

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, cfg, filesChan, &wg, pathStatic, i+1, &atomicCounter, settings.S3Bucket, bar)
	}

	for _, file := range files {
		filesChan <- file
	}

	close(filesChan)
	wg.Wait()

	logger.FInfoFlags(f.IOStreams.Out, msg.UploadSuccessful, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadSuccessful)
	return nil
}
