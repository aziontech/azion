package cmdutil

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
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
		writer.Write(data[:])
	}
	return nil
}
