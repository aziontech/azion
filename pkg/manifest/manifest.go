package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	msgcache "github.com/aziontech/azion-cli/messages/cache_setting"
	msgrule "github.com/aziontech/azion-cli/messages/delete/rules_engine"
	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiApplications "github.com/aziontech/azion-cli/pkg/api/applications"
	apiCache "github.com/aziontech/azion-cli/pkg/api/cache_setting"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/briandowns/spinner"
)

var (
	CacheIds         map[string]int64
	CacheIdsBackup   map[string]int64
	RuleIds          map[string]contracts.RuleIdsStruct
	ConnectorIds     map[string]int64
	DeploymentIds    map[string]int64
	FunctionIds      map[string]contracts.AzionJsonDataFunction
	manifestFilePath = "/.edge/manifest.json"
)

type ManifestInterpreter struct {
	FileReader            func(path string) ([]byte, error)
	GetWorkDir            func() (string, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
}

func NewManifestInterpreter() *ManifestInterpreter {
	return &ManifestInterpreter{
		FileReader:            os.ReadFile,
		GetWorkDir:            utils.GetWorkingDir,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
	}
}

func (man *ManifestInterpreter) ManifestPath() (string, error) {
	pathWorkingDir, err := man.GetWorkDir()
	if err != nil {
		return "", err
	}
	return utils.Concat(pathWorkingDir, manifestFilePath), nil
}

func (man *ManifestInterpreter) ReadManifest(path string, f *cmdutil.Factory, msgs *[]string) (*contracts.ManifestV4, error) {
	logger.FInfoFlags(f.IOStreams.Out, msg.ReadingManifest, f.Format, f.Out)
	*msgs = append(*msgs, msg.ReadingManifest)
	manifest := &contracts.ManifestV4{}

	byteManifest, err := man.FileReader(path)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorReadManifest, err)
	}

	err = json.Unmarshal(byteManifest, &manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (man *ManifestInterpreter) CreateResources(conf *contracts.AzionApplicationOptions, manifest *contracts.ManifestV4, functions map[string]contracts.AzionJsonDataFunction, f *cmdutil.Factory, projectConf string, msgs *[]string) error {
	logger.Debug("Applying manifest resources")
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Suffix = " " + msg.CreatingManifest
	s.FinalMSG = "\n"
	if !f.Debug {
		s.Start()
	}
	defer s.Stop()

	rc := NewResourceContext(f, conf, manifest, projectConf, msgs, man.WriteAzionJsonContent)

	if len(manifest.Functions) > 0 {
		if err := rc.ApplyFunctions(manifest.Functions); err != nil {
			return err
		}
	}

	if len(manifest.Applications) > 0 && len(manifest.Applications[0].FunctionsInstances) > 0 {
		logger.Debug("Applying function instances")
		if err := rc.ApplyFunctionInstances(manifest.Applications[0].FunctionsInstances); err != nil {
			return err
		}
	}

	if len(manifest.Applications) > 0 {
		edgeappman := manifest.Applications[0]
		logger.Debug("Applying edge application")
		if err := rc.ApplyEdgeApplication(edgeappman); err != nil {
			return err
		}

		if len(edgeappman.CacheSettings) > 0 {
			logger.Debug("Applying cache settings")
			if err := rc.ApplyCacheSettings(edgeappman.CacheSettings); err != nil {
				return err
			}
		}

		if len(manifest.Connectors) > 0 {
			logger.Debug("Applying connectors")
			if err := rc.ApplyConnectors(manifest.Connectors); err != nil {
				return err
			}
		}

		if len(edgeappman.Rules) > 0 {
			logger.Debug("Applying rules engine")
			if err := rc.ApplyRulesEngine(edgeappman.Rules); err != nil {
				return err
			}
		}
	}

	if len(manifest.Workloads) > 0 {
		logger.Debug("Applying workloads")
		if err := rc.ApplyWorkloads(manifest.Workloads); err != nil {
			return err
		}
	}

	if len(manifest.WorkloadDeployments) > 0 {
		logger.Debug("Applying workload deployments")
		if err := rc.ApplyWorkloadDeployments(manifest.WorkloadDeployments); err != nil {
			return err
		}
	}

	if len(manifest.Firewalls) > 0 {
		logger.Debug("Applying firewalls")
		if err := rc.ApplyFirewalls(manifest.Firewalls); err != nil {
			return err
		}
	}

	if len(manifest.Purge) > 0 {
		logger.Debug("Applying purge")
		if err := rc.ApplyPurge(manifest.Purge); err != nil {
			return err
		}
	}

	CacheIds = rc.CacheIds
	CacheIdsBackup = rc.CacheIdsBackup
	RuleIds = rc.RuleIds
	ConnectorIds = rc.ConnectorIds
	DeploymentIds = rc.DeploymentIds
	FunctionIds = rc.FunctionIds

	if err := rc.DeleteOrphanedResources(); err != nil {
		return err
	}

	return nil
}

func deleteResources(ctx context.Context, f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string) error {
	client := apiApplications.NewClient(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	clientCache := apiCache.NewClientV4(f.HttpClient, f.Config.GetString("api_v4_url"), f.Config.GetString("token"))
	// clientOrigin := apiOrigin.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	if conf.SkipDeletion != nil && *conf.SkipDeletion {
		logger.FInfoFlags(f.IOStreams.Out, msg.SkipDeletion, f.Format, f.Out)
		*msgs = append(*msgs, msg.SkipDeletion)
		return nil
	}

	for _, value := range RuleIds {
		//since until [UXE-3599] was carried out we'd only cared about "request" phase, this check guarantees that if Phase is empty
		// we are probably dealing with a rule engine from a previous version
		phase := "request"
		if value.Phase != "" {
			phase = value.Phase
		}
		var statusInt int
		var err error
		switch phase {
		case "request":
			statusInt, err = client.DeleteRulesEngineRequest(ctx, conf.Application.ID, phase, value.Id)
		case "response":
			statusInt, err = client.DeleteRulesEngineResponse(ctx, conf.Application.ID, phase, value.Id)
		default:
			return msgrule.ErrorInvalidPhase
		}

		if statusInt == 404 {
			logger.Debug("Rule Engine not found. Skipping delete")
			continue
		}
		if err != nil {
			return err
		}
		msgf := fmt.Sprintf(msgrule.DeleteOutputSuccess+"\n", value.Id)
		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
		*msgs = append(*msgs, msgf)
	}

	for _, value := range CacheIds {
		status, err := clientCache.Delete(ctx, conf.Application.ID, value)
		if status == 404 {
			logger.Debug("Cache Setting not found. Skipping delete")
			continue
		}
		if err != nil {
			return err
		}
		msgf := fmt.Sprintf(msgcache.DeleteOutputSuccess+"\n", value)
		logger.FInfoFlags(f.IOStreams.Out, msgf, f.Format, f.Out)
		*msgs = append(*msgs, msgf)
	}

	return nil
}

func unmarshalJsonArgs(argsPath string) (map[string]interface{}, error) {
	marshalledArgs, err := os.ReadFile(argsPath)
	if err != nil {
		// If args.json file doesn't exist, return empty map as default
		return map[string]interface{}{}, nil
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		return nil, fmt.Errorf("%s: %w", msg.ErrorUnmarshalArgsFile, err)
	}
	return args, nil
}
