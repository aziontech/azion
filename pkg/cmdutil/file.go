package cmdutil

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/pkg/contracts"
	"gopkg.in/yaml.v2"
)

func WriteToFile(data interface{}, opts *contracts.DescribeOptions) error {

	switch opts.Format {
	case "json":
		file, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(opts.OutPath), os.ModePerm)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(opts.OutPath, file, 0644)
		if err != nil {
			return err
		}

	case "yaml":
		file, err := yaml.Marshal(data)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(opts.OutPath), os.ModePerm)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(opts.OutPath, file, 0644)
		if err != nil {
			return err
		}
	default:
		return errors.New("Format not supported. Use --help for more information")

	}
	return nil
}
