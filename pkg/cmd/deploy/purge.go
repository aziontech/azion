package deploy

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	apipurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

type Data struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

func (cmd *DeployCmd) Purge(domain []string, path string) error {
	purgeDomains := make([]string, len(domain))
	for i := 0; i < len(domain); i++ {
		purgeDomains[i] = domain[i] + path
	}
	ctx := context.Background()
	clipurge := apipurge.NewClient(cmd.F.HttpClient, cmd.F.Config.GetString("api_url"), cmd.F.Config.GetString("token"))
	err := clipurge.PurgeWildcard(ctx, purgeDomains)
	if err != nil {
		logger.Debug("Error while purging domain", zap.Error(err))
		return err
	}
	return nil
}

func PurgeForUpdatedFiles(cmd *DeployCmd, domain []string) error {
	currentDataMap, err := ReadFilesJSONL()
	if err != nil {
		return err
	}

	newData, err := ReadFilesEdgeStorage()
	if err != nil {
		return err
	}

	newDataMap := make(map[string]Data)
	for _, newDataItem := range newData {
		newDataMap[newDataItem.Name] = newDataItem
	}

	for _, current := range currentDataMap {
		if newDataItem, exists := newDataMap[current.Name]; exists {
			if current.Hash != newDataItem.Hash {
				if err := cmd.Purge(domain, current.Name); err != nil {
					logger.Debug("Error purge path domain", zap.String("path", current.Name), zap.Error(err))
				}
				logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputCachePurgePath, current.Name))
			}
		}
	}

	jsonl, err := json.MarshalIndent(newData, "  ", " ")
	if err != nil {
		return err
	}

	jsonlFile, err := os.Create("./azion/files.json")
	if err != nil {
		return err
	}
	defer jsonlFile.Close()

	if _, err := jsonlFile.Write(jsonl); err != nil {
		return err
	}

	return nil
}

func ReadFilesJSONL() ([]Data, error) {
	var dt []Data
	file, err := os.Open("./azion/files.json")
	if os.IsNotExist(err) {
		return dt, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&dt); err != nil {
		return nil, err
	}
	return dt, nil
}

func ReadFilesEdgeStorage() ([]Data, error) {
	var data []Data
	err := filepath.Walk(".edge/storage", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			content, err := os.ReadFile(path)
			if err != nil {
				logger.Debug("Error read file", zap.Error(err))
				return err
			}
			dt := Data{
				Name: path,
				Hash: fmt.Sprintf("%x", sha256.Sum256(content)),
			}
			data = append(data, dt)
		}
		return nil
	})
	if err != nil {
		logger.Debug("Error filepath walk directory", zap.Error(err))
		return nil, err
	}
	return data, nil
}
