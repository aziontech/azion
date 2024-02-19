package metric

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/spf13/cobra"
)

const (
	metricsFilename = "metrics.json"
	total           = "total-commands-executed"
	totalSuccess    = "total-commands-executed-successfully"
	totalFailed     = "total-commands-executed-unsuccessfully"
	failSuffix      = "-failed"
)

func TotalCommandsCount(cmd *cobra.Command, commandName string, success bool) error {
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

	var data map[string]int

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil && err != io.EOF {
		return err
	}

	// If EOF is encountered or the file is empty, initialize data as an empty map
	if data == nil {
		data = make(map[string]int)
	}

	if !success {
		commandName = commandName + failSuffix
	}

	data[commandName]++
	data[total]++
	if success {
		data[totalSuccess]++
	} else {
		data[totalFailed]++
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
