package deploy

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
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

		if info.IsDir() && (strings.Contains(path, "node_modules") || strings.Contains(path, ".edge")) {
			logger.Debug("Skipping directory: " + path)
			return filepath.SkipDir
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
func addFileToZip(zipWriter *zip.Writer, fileNameZip, baseDir string) (os.FileInfo, error) {
	file, err := os.Open(fileNameZip)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Create an entry in the zip with the relative path
	relativePath := strings.TrimPrefix(fileNameZip, baseDir)
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
func uploadFile(ctx context.Context, cfg aws.Config, bucketName, filePath, prefix string, fileInfo fs.FileInfo) error {
	s3Client := s3.NewFromConfig(cfg)

	file, err := os.Open(filePath)

	if err != nil {
		return fmt.Errorf(msg.ErrorOpenFile, filePath, err)
	}
	defer file.Close()

	fileName := filepath.Join(prefix, fileInfo.Name())

	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName), // Name of the file in the bucket
		Body:   file,
	}

	_, err = s3Client.PutObject(ctx, uploadInput)
	if err != nil {
		return fmt.Errorf(msg.ErrorUploadFileBucket, bucketName, err)
	}

	return nil
}

// Worker responsible for zipping and uploading files
func worker(ctx context.Context, cancel context.CancelFunc, cfg aws.Config,
	filesChan <-chan string, wg *sync.WaitGroup, baseDir string, workerID int,
	atomicCounter *int32, bucketName string, progressBar *progressbar.ProgressBar,
	errChan chan error, prefix string) {
	defer wg.Done()

	var (
		currentSize int64
		filesToZip  []string
	)

	zipFileName := fmt.Sprintf("project_part_%d.zip", workerID)

	for {
		select {
		case file, ok := <-filesChan:
			if !ok {
				if len(filesToZip) > 0 {
					err := zipAndUpload(ctx, cfg, filesToZip, zipFileName,
						baseDir, atomicCounter, bucketName, prefix)
					if err != nil {
						errChan <- err
						cancel()
						return
					}
					if progressBar != nil {
						progressBar.Add(len(filesToZip)) // nolint
					}
				}
				return
			}

			filesToZip = append(filesToZip, file)

			fileInfo, err := os.Stat(file)
			if err != nil {
				errChan <- fmt.Errorf(msg.ErrorGetFileInfo, file, err)
				cancel()
				return
			}
			currentSize += fileInfo.Size()

			// If the current size exceeds the limit (1MB), zip and upload
			if currentSize >= maxZipSize {
				err = zipAndUpload(ctx, cfg, filesToZip, zipFileName, baseDir,
					atomicCounter, bucketName, prefix)
				if err != nil {
					errChan <- err
					cancel()
					return
				}

				if progressBar != nil {
					progressBar.Add(len(filesToZip)) // nolint
				}

				// Reset the size counter and clear the file list
				currentSize = 0
				filesToZip = nil
			}
		case <-ctx.Done():
			continue
		}
	}
}

func zipAndUpload(ctx context.Context, cfg aws.Config, filesToZip []string,
	zipFileName, baseDir string, atomicCounter *int32,
	bucketName string, prefix string) error {

	// Increment atomic counter to generate unique zip file name
	counter := atomic.AddInt32(atomicCounter, 1)
	zipFileName = fmt.Sprintf("%s_%d.zip", strings.TrimSuffix(zipFileName, ".zip"), counter)

	zipFilePath := filepath.Join(baseDir, zipFileName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf(msg.ErrorCreateZip, zipFileName, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add files to the zip
	for _, file := range filesToZip {
		_, err = addFileToZip(zipWriter, file, baseDir)
		if err != nil {
			return fmt.Errorf(msg.ErrorAddFileZip, file, err)
		}
	}

	err = zipWriter.Close()
	if err != nil {
		return fmt.Errorf(msg.ErrorCloseFileZip, zipFileName, err)
	}

	// Check if the zip file was created
	fileInfo, err := os.Stat(zipFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf(msg.ErrorZipNotExist, zipFileName)
	}

	err = uploadFile(ctx, cfg, bucketName, zipFilePath, prefix, fileInfo)
	if err != nil {
		return fmt.Errorf(msg.ErrorUploadZip, zipFileName, err)
	}

	err = os.Remove(zipFilePath)
	if err != nil {
		return fmt.Errorf(msg.ErrorDelFileZip, zipFileName, err)
	}

	return nil
}

func uploadFiles(f *cmdutil.Factory, msgs *[]string, prefix string, pathStatic string, settings token.Settings) error {
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
		return errors.New(msg.ErrorUnableSDKConfig + err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)

	var bar *progressbar.ProgressBar = nil
	if !f.Silent {
		bar = progressbar.NewOptions(
			len(files),
			progressbar.OptionSetDescription("Uploading files"),
			progressbar.OptionShowCount(),
			progressbar.OptionSetWriter(f.IOStreams.Out),
			progressbar.OptionClearOnFinish(),
		)
	}

	logger.FInfoFlags(f.IOStreams.Out, msg.UploadStart, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadStart)

	filesChan := make(chan string, len(files))

	var wg sync.WaitGroup
	var errWg sync.WaitGroup

	// Atomic counter to guarantee unique zip file names
	var atomicCounter int32

	var uploadErr error
	errWg.Add(1)
	go func() {
		defer errWg.Done()
		for err := range errChan {
			if err != nil {
				uploadErr = err
				cancel()
				break
			}
		}
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, cancel, cfg, filesChan, &wg, pathStatic, i+1, &atomicCounter, settings.S3Bucket, bar, errChan, prefix)
	}

	go func() {
		defer close(filesChan)
		for _, file := range files {
			filesChan <- file
		}
	}()

	wg.Wait()
	close(errChan)
	errWg.Wait()

	if uploadErr != nil {
		return uploadErr
	}

	logger.FInfoFlags(f.IOStreams.Out, msg.UploadSuccessful, f.Format, f.Out)
	*msgs = append(*msgs, msg.UploadSuccessful)

	return nil
}
