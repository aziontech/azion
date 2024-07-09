package deploy

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	apidom "github.com/aziontech/azion-cli/pkg/api/domain"
	apipurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

type Data struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

func (cmd *DeployCmd) PurgeWildcard(domain []string, path string) error {
	purgeDomains := make([]string, len(domain))
	for i := 0; i < len(domain); i++ {
		purgeDomains[i] = domain[i] + path
	}
	ctx := context.Background()
	clipurge := apipurge.NewClient(cmd.F.HttpClient, cmd.F.Config.GetString("api_url"), cmd.F.Config.GetString("token"))
	err := clipurge.PurgeWildcard(ctx, purgeDomains)
	if err != nil {
		logger.Debug("Error while purging wildcard domain", zap.Error(err))
		return err
	}
	return nil
}

func (cmd *DeployCmd) PurgeUrls(domain []string, path string) error {
	purgeDomains := make([]string, len(domain))
	for i := 0; i < len(domain); i++ {
		purgeDomains[i] = domain[i] + path
	}
	ctx := context.Background()
	clipurge := apipurge.NewClient(cmd.F.HttpClient, cmd.F.Config.GetString("api_url"), cmd.F.Config.GetString("token"))
	err := clipurge.PurgeUrls(ctx, purgeDomains)
	if err != nil {
		logger.Debug("Error while purging urls domain", zap.Error(err))
		return err
	}
	return nil
}

func PurgeForUpdatedFiles(cmd *DeployCmd, domain apidom.DomainResponse, confPath string, msgs *[]string) error {
	listURLsDomains := domain.GetCnames()
	if !domain.GetCnameAccessOnly() {
		listURLsDomains = append(listURLsDomains, domain.GetDomainName())
	}

	currentDataMap, err := ReadFilesJSONL()
	if err != nil {
		return err
	}

	if currentDataMap == nil {
		wildCard := "/*"
		for _, v := range listURLsDomains {
			if err := cmd.PurgeWildcard([]string{v}, wildCard); err != nil {
				logger.Debug("Error purge path domain", zap.String("wildCard", wildCard), zap.Error(err))
			}
			msgsf := fmt.Sprintf(msg.DeployOutputCachePurgeWildCard, v)
			logger.FInfoFlags(cmd.F.IOStreams.Out, msgsf, cmd.F.Format, cmd.F.Out)
			*msgs = append(*msgs, msgsf)
		}
	}

	newData, err := ReadFilesEdgeStorage()
	if err != nil {
		return err
	}

	if currentDataMap != nil {
		newDataMap := make(map[string]Data)
		for _, newDataItem := range newData {
			newDataMap[newDataItem.Name] = newDataItem
		}

		for _, current := range currentDataMap {
			if newDataItem, exists := newDataMap[current.Name]; exists {
				if current.Hash != newDataItem.Hash {
					path := strings.TrimPrefix(current.Name, ".edge/storage")
					if err := cmd.PurgeUrls(listURLsDomains, path); err != nil {
						logger.Debug("Error purge path domain", zap.String("path", path), zap.Error(err))
					}
					logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputCachePurgeUrl, current.Name))
				}
			}
		}
	}

	jsonl, err := json.MarshalIndent(newData, "  ", " ")
	if err != nil {
		return err
	}

	jsonlFile, err := os.Create(path.Join(confPath, "files.json"))
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
