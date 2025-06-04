package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	sdkstorage "github.com/aziontech/azionapi-go-sdk/storage"
	"github.com/aziontech/azionapi-v4-go-sdk/edge"
	"go.uber.org/zap"
)

func WriteManifest(manifest *contracts.ManifestV4, pathMan string) error {

	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		logger.Debug("Error marshalling response", zap.Error(err))
		return msg.ERRORMARSHALMANIFEST
	}

	err = os.WriteFile(path.Join(pathMan, "manifesttoconvert.json"), data, 0644)
	if err != nil {
		logger.Debug("Error writing file", zap.Error(err))
		return msg.ERRORWRITEMANIFEST
	}

	return nil
}

func (cmd *DeployCmd) firstRunManifestToConfig(conf *contracts.AzionApplicationOptions) error {

	truePointer := true
	appManifest := contracts.EdgeApplications{
		Name: conf.Name,
		Modules: &contracts.Modules{
			EdgeFunctionsEnabled: &truePointer,
		},
		Active: &truePointer,
	}

	storageType := edge.EdgeConnectorStorageTypePropertiesRequest{
		Bucket: &conf.Bucket,
		Prefix: &conf.Prefix,
	}
	storageConnector := edge.EdgeConnectorStorageTypedRequest{
		Name:           conf.Name,
		Active:         &truePointer,
		TypeProperties: storageType,
	}
	connectorManifest := edge.EdgeConnectorPolymorphicRequest{
		EdgeConnectorStorageTypedRequest: &storageConnector,
	}

	functionMan := contracts.EdgeFunction{
		Name:   conf.Name,
		Target: ".edge/worker.js",
		Bindings: contracts.Bindings{
			EdgeStorage: contracts.EdgeStorage{
				Bucket: conf.Bucket,
				Prefix: conf.Prefix,
			},
		},
	}

	storageMan := sdkstorage.BucketCreate{
		Name:       conf.Bucket,
		EdgeAccess: sdkstorage.READ_ONLY,
	}

	manifestToConfig := &contracts.ManifestV4{}
	manifestToConfig.EdgeConnectors = append(manifestToConfig.EdgeConnectors, connectorManifest)
	manifestToConfig.EdgeApplications = append(manifestToConfig.EdgeApplications, appManifest)
	manifestToConfig.EdgeFunctions = append(manifestToConfig.EdgeFunctions, functionMan)
	manifestToConfig.EdgeStorage = append(manifestToConfig.EdgeStorage, storageMan)

	err := cmd.WriteManifest(manifestToConfig, "")
	if err != nil {
		return err
	}
	defer os.Remove("manifesttoconvert.json")

	vul := vulcanPkg.NewVulcan()
	command := vul.Command("", "manifest -o %s transform %s", cmd.F)
	format, err := findAzionConfig()
	if err != nil {
		format = ".mjs"
	}
	fileName := fmt.Sprintf("azion.config%s", format)
	err = cmd.commandRunInteractive(cmd.F, fmt.Sprintf(command, fileName, "manifesttoconvert.json"))
	if err != nil {
		return err
	}
	err = cmd.callBundlerInit(conf)
	if err != nil {
		return nil
	}

	return nil
}

func findAzionConfig() (string, error) {
	extensions := []string{".cjs", ".mjs", ".js"}
	baseName := "azion.config"

	for _, ext := range extensions {
		filename := baseName + ext
		if _, err := os.Stat(filename); err == nil {
			return ext, nil
		} else if !os.IsNotExist(err) {
			return "", fmt.Errorf("error checking file %s: %w", filename, err)
		}
	}
	return "", fmt.Errorf("no azion.config file found")
}
