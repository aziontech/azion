package deploy

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aziontech/azion-cli/pkg/contracts"
)

// Função para dividir os arquivos em lotes de até 1MB e criar ZIPs mantendo a estrutura de diretórios
func CreateZipsInBatches(files []contracts.FileOps) error {
	const maxBatchSize = 1 * 1024 * 1024 // 1MB em bytes
	var currentBatch []contracts.FileOps
	var currentSize int64 = 0
	tempDir := os.TempDir()

	batchNumber := 1

	for _, fileOp := range files {
		// Obter o tamanho do arquivo
		info, err := fileOp.FileContent.Stat()
		if err != nil {
			return fmt.Errorf("erro ao obter info do arquivo %s: %v", fileOp.Path, err)
		}
		fileSize := info.Size()

		// Verificar se adicionar este arquivo excede o tamanho máximo do lote
		if currentSize+fileSize > maxBatchSize && len(currentBatch) > 0 {
			// Criar ZIP para o lote atual
			if err := createZip(currentBatch, tempDir, batchNumber); err != nil {
				return err
			}

			batchNumber++
			// Resetar o lote atual
			currentBatch = []contracts.FileOps{}
			currentSize = 0
		}

		// Adicionar o arquivo ao lote atual
		currentBatch = append(currentBatch, fileOp)
		currentSize += fileSize
	}

	// Criar ZIP para quaisquer arquivos restantes
	if len(currentBatch) > 0 {
		if err := createZip(currentBatch, tempDir, batchNumber); err != nil {
			return err
		}
	}

	return nil
}

// Função auxiliar para criar um arquivo ZIP a partir de um lote de FileOps mantendo a estrutura de diretórios
func createZip(batch []contracts.FileOps, destDir string, batchNumber int) error {
	zipFileName := fmt.Sprintf("batch_%d.zip", batchNumber)
	zipFilePath := filepath.Join(destDir, zipFileName)

	// Criar o arquivo ZIP
	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo ZIP %s: %v", zipFilePath, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, fileOp := range batch {
		relPath, err := filepath.Rel(destDir, fileOp.Path)
		if err != nil {
			return fmt.Errorf("erro ao determinar caminho relativo para %s: %v", fileOp.Path, err)
		}

		// Substituir separadores de diretório para compatibilidade com ZIP
		zipPath := filepath.ToSlash(relPath)

		// Adicionar um arquivo ao ZIP com o caminho relativo
		writer, err := zipWriter.Create(zipPath)
		if err != nil {
			return fmt.Errorf("erro ao adicionar arquivo %s ao ZIP: %v", zipPath, err)
		}

		if _, err := fileOp.FileContent.Seek(0, io.SeekStart); err != nil {
			return fmt.Errorf("erro ao resetar ponteiro do arquivo %s: %v", fileOp.Path, err)
		}

		// Copiar o conteúdo do arquivo para o ZIP
		_, err = io.Copy(writer, fileOp.FileContent)
		if err != nil {
			return fmt.Errorf("erro ao copiar conteúdo do arquivo %s para o ZIP: %v", fileOp.Path, err)
		}
	}

	fmt.Printf("Criado %s com %d arquivos\n", zipFilePath, len(batch))
	return nil
}

// ReadZip lê os arquivos ZIP no diretório pathStatic que começam com prefix e terminam com .zip.
// Retorna uma lista de contracts.FileOps ou um erro.
func ReadZip() ([]contracts.FileOps, error) {
	var listZIP []contracts.FileOps

	tempDir := os.TempDir()

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return []contracts.FileOps{}, fmt.Errorf("erro ao ler o diretório: %v", err)
	}

	for _, f := range files {
		if strings.ToLower(filepath.Ext(f.Name())) == ".zip" &&
			strings.HasPrefix(f.Name(), "batch") {
			pathZIP := filepath.Join(tempDir, f.Name())
			fmt.Printf("Processando arquivo ZIP: %s\n", pathZIP)

			fmt.Println(">>> ", tempDir)
			fmt.Println(">>> ", pathZIP)

			f, err := os.Open(pathZIP)
			if err != nil {
				log.Printf("Erro ao abrir o ZIP %s: %v\n", pathZIP, err)
				continue
			}

			fileString := strings.TrimPrefix(pathZIP, tempDir)
			// mimeType, err := mimemagic.MatchFilePath(fileString, -1)
			// // if err != nil {
			// // 	return []contracts.FileOps{}, err
			// // }
			//
			fileOptions := contracts.FileOps{
				Path: fileString,
				// 	MimeType:    mimeType.MediaType(),
				FileContent: f,
			}
			//
			listZIP = append(listZIP, fileOptions)

		}

	}

	return listZIP, nil
}
