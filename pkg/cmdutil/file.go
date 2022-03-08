package cmdutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteDetailsToFile(data []byte, outPath string, writer io.Writer) error {

	err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outPath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func UnmarshallJsonFromReader(file io.Reader, object interface{}) error {
	jsonFile, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonFile, &object)
	if err != nil {
		return err
	}

	return nil
}
