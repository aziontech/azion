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
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk/edge-api"
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
		Modules: &edgesdk.EdgeApplicationModulesRequest{
			EdgeFunctions: &edgesdk.EdgeFunctionModuleRequest{
				Enabled: &truePointer,
			},
		},
		Active: &truePointer,
	}

	storageType := edgesdk.EdgeConnectorStorageAttributesRequest{
		Bucket: conf.Bucket,
		Prefix: &conf.Prefix,
	}
	storageConnector := edgesdk.EdgeConnectorStorageRequest{
		Name:       conf.Name,
		Active:     &truePointer,
		Attributes: storageType,
	}
	connectorManifest := edgesdk.EdgeConnectorPolymorphicRequest{
		EdgeConnectorStorageRequest: &storageConnector,
	}

	functionMan := contracts.EdgeFunction{
		Name:     conf.Name,
		Argument: ".edge/worker.js",
		Bindings: contracts.FunctionBindings{
			Storage: contracts.StorageBinding{
				Bucket: conf.Bucket,
				Prefix: conf.Prefix,
			},
		},
	}

	storageMan := contracts.EdgeStorageManifest{
		Name:       conf.Bucket,
		EdgeAccess: "read_only",
		Dir:        conf.Prefix,
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
