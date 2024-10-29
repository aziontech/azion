package deploy

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

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

const maxZipSize = 1 * 1024 * 1024 // 1MB

func criarZipEmMemoria(arquivos []contracts.FileOps) ([]byte, error) {
	// Cria um buffer para armazenar o ZIP em memória
	var buffer bytes.Buffer

	// Cria um novo writer para o ZIP
	zipWriter := zip.NewWriter(&buffer)
	defer func() {
		if err := zipWriter.Close(); err != nil {
			fmt.Println("Erro ao fechar o writer do ZIP:", err)
		}
	}()

	// Itera sobre os arquivos e os adiciona ao ZIP
	for _, arquivo := range arquivos {
		// Cria o cabeçalho do arquivo no ZIP
		header := &zip.FileHeader{
			Name:   filepath.Base(arquivo.Path),
			Method: zip.Deflate, // Método de compressão
		}

		// Define o MIME type como comentário (opcional)
		if arquivo.MimeType != "" {
			header.Comment = arquivo.MimeType
		} else {
			// Tenta detectar o MIME type com base na extensão
			mimeType := mime.TypeByExtension(filepath.Ext(arquivo.Path))
			header.Comment = mimeType
		}

		// Cria o writer para o arquivo dentro do ZIP
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return nil, fmt.Errorf("não foi possível criar o writer para o arquivo %s: %v", arquivo.Path, err)
		}

		// Reposiciona o cursor do arquivo para o início
		_, err = arquivo.FileContent.Seek(0, io.SeekStart)
		if err != nil {
			return nil, fmt.Errorf("não foi possível reposicionar o cursor do arquivo %s: %v", arquivo.Path, err)
		}

		// Copia o conteúdo do arquivo para o writer do ZIP
		_, err = io.Copy(writer, arquivo.FileContent)
		if err != nil {
			return nil, fmt.Errorf("não foi possível copiar o conteúdo do arquivo %s para o ZIP: %v", arquivo.Path, err)
		}
	}

	return buffer.Bytes(), nil
}

func ReadAllFiles(pathStatic string, cmd *DeployCmd) ([]contracts.FileOps, error) {
	var listFiles []contracts.FileOps

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
	cfg, err := s3.New(settings.S3AccessKey, settings.S3SecreKey)
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
			for Retries < 5 {
				atomic.AddInt64(&Retries, 1)
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

			if Retries >= 5 {
				logger.Debug("There have been 5 retries already, quitting upload")
				results <- err
				return
			}
		}

		atomic.AddInt64(currentFile, 1)
		results <- nil
	}
}
