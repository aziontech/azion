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
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/cheggaaa/pb/v3"
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

// ResolveEndpoint é o método que define o endpoint customizado
func (e *CustomEndpointResolver) ResolveEndpoint(service, region string) (aws.Endpoint, error) {
	return aws.Endpoint{
		URL:           e.URL,
		SigningRegion: e.SigningRegion,
	}, nil
}

// Função que lê todos os arquivos de um diretório e retorna um slice com os caminhos completos
func getFilesFromDir(dirPath string) ([]string, error) {
	var files []string

	// Função que será aplicada a cada item encontrado no diretório
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Verifica se o item é um arquivo (e não um diretório)
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

// Função para zipar arquivos com limite de tamanho de 1MB
func zipFilesInChunks(baseDir string, files []string, baseZipName string) ([]string, error) {
	var zipFiles []string
	var zipFileName string
	var zipFile *os.File
	var zipWriter *zip.Writer
	var currentSize int64

	for _, file := range files {
		// Se o arquivo zip não existir ou excedeu o tamanho, cria um novo
		if zipFile == nil || currentSize >= maxZipSize {
			if zipFile != nil {
				zipWriter.Close()
				zipFile.Close()
				zipFiles = append(zipFiles, zipFileName)
			}

			zipFileName = fmt.Sprintf("%s_part%d.zip", baseZipName, len(zipFiles)+1)
			var err error
			zipFile, err = os.Create(zipFileName)
			if err != nil {
				return nil, err
			}
			zipWriter = zip.NewWriter(zipFile)
			currentSize = 0
		}

		// Adiciona o arquivo ao zip
		fileInfo, err := addFileToZip(zipWriter, file, baseDir)
		if err != nil {
			return nil, err
		}

		// Atualiza o tamanho atual
		currentSize += fileInfo.Size()
	}

	// Fecha o último zip e adiciona à lista
	if zipFile != nil {
		zipWriter.Close()
		zipFile.Close()
		zipFiles = append(zipFiles, zipFileName)
	}

	return zipFiles, nil
}

// Função para adicionar um arquivo individual ao arquivo zip, preservando o caminho relativo
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

	// Criar uma entrada no zip com o caminho relativo
	relativePath := strings.TrimPrefix(filePath, baseDir)
	writer, err := zipWriter.Create(relativePath)
	if err != nil {
		return nil, err
	}

	// Copiar o conteúdo do arquivo para o zip
	_, err = io.Copy(writer, file)
	if err != nil {
		return nil, err
	}

	return fileInfo, nil
}

// Função para fazer o upload de um arquivo para o bucket
func uploadFile(ctx context.Context, cfg aws.Config, bucketName, filePath string) error {
	// Criar o cliente S3 com a configuração fornecida
	s3Client := s3.NewFromConfig(cfg)

	// Abrir o arquivo para upload
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	// Obter o nome do arquivo
	fileName := filepath.Base(filePath)

	// Configurar os parâmetros do upload, sem a ACL
	uploadInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileName), // Nome do arquivo no bucket
		Body:   file,                 // Conteúdo do arquivo
	}

	// Fazer o upload do arquivo
	_, err = s3Client.PutObject(ctx, uploadInput)
	if err != nil {
		return fmt.Errorf("failed to upload file to bucket %s: %w", bucketName, err)
	}

	// fmt.Printf("File %s uploaded successfully to bucket %s\n", fileName, bucketName)
	return nil
}

// Worker responsável por zipar e fazer o upload de arquivos
func worker(ctx context.Context, cfg aws.Config, filesChan <-chan string, wg *sync.WaitGroup, baseDir string, workerID int, atomicCounter *int32, bucketName string, progressBar *pb.ProgressBar) {
	defer wg.Done()

	var (
		currentSize int64
		filesToZip  []string
	)

	zipFileName := fmt.Sprintf("project_part_%d.zip", workerID)

	for file := range filesChan {
		// Adiciona o arquivo à lista de arquivos a serem zipados
		filesToZip = append(filesToZip, file)

		// Obter o tamanho do arquivo e adicionar ao tamanho atual
		fileInfo, err := os.Stat(file)
		if err != nil {
			log.Printf("Failed to get file info for %s: %v", file, err)
			continue
		}
		currentSize += fileInfo.Size()

		// Se o tamanho atual exceder o limite (1MB), zipar e fazer upload
		if currentSize >= maxZipSize {
			// Criar o arquivo zip e fazer o upload
			err = zipAndUpload(ctx, cfg, filesToZip, zipFileName, baseDir, atomicCounter, bucketName, progressBar)
			if err != nil {
				log.Printf("Failed to zip and upload files: %v", err)
			}

			// Reiniciar o contador de tamanho e limpar a lista de arquivos
			currentSize = 0
			filesToZip = nil
		}
	}

	// Se ainda houver arquivos no final, zipar e fazer upload
	if len(filesToZip) > 0 {
		err := zipAndUpload(ctx, cfg, filesToZip, zipFileName, baseDir, atomicCounter, bucketName, progressBar)
		if err != nil {
			log.Printf("Failed to zip and upload files: %v", err)
		}
	}
}

func zipAndUpload(ctx context.Context, cfg aws.Config, filesToZip []string, zipFileName, baseDir string, atomicCounter *int32, bucketName string, progressBar *pb.ProgressBar) error {
	// Incrementar contador atomic para gerar nome de arquivo zip único
	counter := atomic.AddInt32(atomicCounter, 1)
	zipFileName = fmt.Sprintf("%s_%d.zip", strings.TrimSuffix(zipFileName, ".zip"), counter)

	// Criar o arquivo zip
	zipFilePath := filepath.Join(baseDir, zipFileName)
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to create zip file %s: %w", zipFileName, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Adicionar arquivos ao zip
	for _, file := range filesToZip {
		_, err = addFileToZip(zipWriter, file, baseDir)
		if err != nil {
			return fmt.Errorf("failed to add file to zip %s: %w", file, err)
		}
	}

	// Fechar o zipWriter antes de tentar fazer o upload
	err = zipWriter.Close()
	if err != nil {
		return fmt.Errorf("failed to close zip file %s: %w", zipFileName, err)
	}

	// Verificar se o arquivo zip foi criado
	if _, err := os.Stat(zipFilePath); os.IsNotExist(err) {
		return fmt.Errorf("zip file %s does not exist", zipFileName)
	}

	// Fazer o upload do arquivo zipado
	err = uploadFile(ctx, cfg, bucketName, zipFilePath)
	if err != nil {
		return fmt.Errorf("failed to upload zip file %s: %w", zipFileName, err)
	}

	progressBar.Increment()

	return nil
}

func uploadFiles(pathStatic string, settings token.Settings) error {
	files, err := getFilesFromDir(pathStatic)
	if err != nil {
		return errors.New("")
	}

	// Criar a configuração AWS com o endpoint correto
	endpointResolver := &CustomEndpointResolver{
		URL:           endpoint,
		SigningRegion: region,
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(settings.S3AccessKey, settings.S3SecreKey, "")),
		config.WithEndpointResolver(endpointResolver),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ctx := context.TODO()

	progressBar := pb.StartNew(len(files))

	// Canal para enviar arquivos para os workers
	filesChan := make(chan string)

	// WaitGroup para aguardar todos os workers terminarem
	var wg sync.WaitGroup

	// Atomic counter para garantir nomes únicos dos arquivos zip
	var atomicCounter int32

	// Iniciar os workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, cfg, filesChan, &wg, pathStatic, i+1, &atomicCounter, settings.S3Bucket, progressBar)
	}

	// Enviar arquivos para os workers
	for _, file := range files {
		filesChan <- file
	}

	// Fechar o canal e aguardar os workers
	close(filesChan)
	wg.Wait()

	progressBar.Finish()

	fmt.Println("Upload completed... ")
	return nil
}
