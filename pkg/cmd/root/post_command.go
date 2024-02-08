package root

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/spf13/cobra"
)

const (
	metricsFilename = "metrics.json"
	total           = "total-commands-executed"
)

func saveMetrict(cmd *cobra.Command) error {
	//1 = authorize; anything different than 1 means that the user did not authorize metrics collection, or did not answer the question yet
	if globalSettings.AuthorizeMetricsCollection != 1 {
		return nil
	}
	dir, err := config.Dir()
	if err != nil {
		return err
	}

	ignoredWords := map[string]bool{
		"__complete": true,
	}
	if ignoredWords[cmd.Parent().Name()] || ignoredWords[cmd.Name()] {
		return nil
	}

	metricsLocation := filepath.Join(dir, metricsFilename)

	file, err := os.OpenFile(metricsLocation, os.O_RDWR|os.O_CREATE, 0666)
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

	commandRun := cmd.Parent().Name() + "-" + cmd.Name()
	rewrittenCommand := strings.TrimPrefix(commandRun, "azion-")

	data[rewrittenCommand]++
	data[total]++

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
