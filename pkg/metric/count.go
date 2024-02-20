package metric

import (
	"encoding/json"
	"github.com/aziontech/azion-cli/pkg/cmd/version"
	"github.com/aziontech/azion-cli/pkg/github"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/spf13/cobra"
)

const (
	metricsFilename = "metrics.json"
)

type command struct {
	TotalSuccess  int
	TotalFailed   int
	ExecutionTime float64
	VersionCLI    string
	VersionVulcan string
}

func TotalCommandsCount(cmd *cobra.Command, commandName string, executionTime float64, success bool) error {
	if commandName == "" {
		return nil
	}

	dir, err := config.Dir()
	if err != nil {
		return err
	}

	ignoredWords := map[string]bool{
		"__complete": true,
		"completion": true,
	}

	if ignoredWords[cmd.Name()] {
		return nil
	}

	metricsLocation := filepath.Join(dir, metricsFilename)

	file, err := os.OpenFile(metricsLocation, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	var data map[string]*command

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil && err != io.EOF {
		return err
	}

	// If EOF is encountered or the file is empty, initialize data as an empty map
	if data == nil {
		data = make(map[string]*command)
	}

	if data[commandName] == nil {
		data[commandName] = &command{}
	}

	data[commandName].ExecutionTime = executionTime
	tagName, err := github.GetVersionGitHub("vulcan")
	if err != nil {
		return err
	}
	data[commandName].VersionCLI = version.BinVersion
	data[commandName].VersionVulcan = tagName[1:]
	if success {
		data[commandName].TotalSuccess++
	} else {
		data[commandName].TotalFailed++
	}

	// Reset file offset to the beginning
	_, err = file.Seek(0, 0)
	if err != nil {
		return err
	}
	// Truncate the file in case the new content is smaller than the previous one
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	// Encode and write the updated map back to the file
	if err := json.NewEncoder(file).Encode(data); err != nil {
		log.Fatal(err)
	}

	return nil
}
