package cmdutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
)

func WriteDetailsToScreenOrFile(data []byte, out bool, outPath string, writer io.Writer) error {

	if out {
		err := os.MkdirAll(filepath.Dir(outPath), os.ModePerm)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(outPath, data, 0644)
		if err != nil {
			return err
		}
	} else {
		_, err := writer.Write(data[:])
		if err != nil {
			return err
		}

	}
	return nil
}

func UnmarshallFunctionUpdateIntoFile(inPath string, object *api.UpdateRequest) error {

	jsonFile, err := ioutil.ReadFile(inPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonFile, &object)
	if err != nil {
		return err
	}

	return nil

}
