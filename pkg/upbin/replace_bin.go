package upbin

import (
	"fmt"
	"github.com/aziontech/azion-cli/utils"
	"io"
	"net/http"
	"os"
)

const urlDownloadAzioncli string = "https://downloads.azion.com/%s/%s/azioncli"

var (
	downloadFileFunc  = downloadFile
	replaceFileFunc   = replaceFile
	replaceBinaryFunc = replaceBinary
)

func replaceBinary(filePath string) (err error) {
	err = downloadFileFunc(filePath)
	if err != nil {
		return err
	}

	err = replaceFileFunc(filePath)
	if err != nil {
		return err
	}

	return
}

func downloadFile(filePath string) error {
	fileURL, err := prepareURLFunc()
	if err != nil {
		return err
	}

	response, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func replaceFile(filePath string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	err = os.Rename(filePath, exe)
	if err != nil {
		return err
	}

	return nil
}

func prepareURL() (string, error) {
	sysAr := GetInfoSystem()
	if sysAr.System == unknown || sysAr.Arch == unknown {
		return "", utils.ErrorUnknownSystem
	}
	return fmt.Sprintf(urlDownloadAzioncli, sysAr.System, sysAr.Arch), nil
}
