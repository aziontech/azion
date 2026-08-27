package manifest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	"go.uber.org/zap"
)

// TimingCallback is a callback function type for reporting timing
type TimingCallback func(name string, duration time.Duration)

// GlobalTimingCallback is the global callback for timing reports
var GlobalTimingCallback TimingCallback

// CacheIds and RuleIds carry the resources tracked by azion.json into
// deleteResources, which removes the ones the manifest no longer declares.
// Every other id map lives on ResourceContext for the duration of a run.
var (
	CacheIds         map[string]int64
	RuleIds          map[string]contracts.RuleIdsStruct
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

func (man *ManifestInterpreter) CreateResources(conf *contracts.AzionApplicationOptions, manifest *contracts.ManifestV4, f *cmdutil.Factory, projectConf string, msgs *[]string) error {
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
		start := time.Now()
		if err := rc.ApplyFunctions(manifest.Functions); err != nil {
			return err
		}
		if GlobalTimingCallback != nil {
			GlobalTimingCallback("ManifestFunctions", time.Since(start))
		}
	}

	if len(manifest.Applications) > 0 && len(manifest.Applications[0].FunctionsInstances) > 0 {
		logger.Debug("Applying function instances")
		start := time.Now()
		if err := rc.ApplyFunctionInstances(manifest.Applications[0].FunctionsInstances); err != nil {
			return err
		}
		if GlobalTimingCallback != nil {
			GlobalTimingCallback("ManifestFunctionInstances", time.Since(start))
		}
	}

	if len(manifest.Applications) > 0 {
		edgeappman := manifest.Applications[0]
		logger.Debug("Applying edge application")
		start := time.Now()
		if err := rc.ApplyEdgeApplication(edgeappman); err != nil {
			return err
		}
		if GlobalTimingCallback != nil {
			GlobalTimingCallback("ManifestEdgeApplication", time.Since(start))
		}

		if len(edgeappman.CacheSettings) > 0 {
			logger.Debug("Applying cache settings")
			start := time.Now()
			if err := rc.ApplyCacheSettings(edgeappman.CacheSettings); err != nil {
				return err
			}
			if GlobalTimingCallback != nil {
				GlobalTimingCallback("ManifestCacheSettings", time.Since(start))
			}
		}

		if len(manifest.Connectors) > 0 {
			logger.Debug("Applying connectors")
			start := time.Now()
			if err := rc.ApplyConnectors(manifest.Connectors); err != nil {
				return err
			}
			if GlobalTimingCallback != nil {
				GlobalTimingCallback("ManifestConnectors", time.Since(start))
			}
		}

		if len(edgeappman.Rules) > 0 {
			logger.Debug("Applying rules engine")
			start := time.Now()
			if err := rc.ApplyRulesEngine(edgeappman.Rules); err != nil {
				return err
			}
			if GlobalTimingCallback != nil {
				GlobalTimingCallback("ManifestRulesEngine", time.Since(start))
			}
		}
	}

	if len(manifest.Workloads) > 0 {
		logger.Debug("Applying workloads")
		start := time.Now()
		if err := rc.ApplyWorkloads(manifest.Workloads); err != nil {
			return err
		}
		if GlobalTimingCallback != nil {
			GlobalTimingCallback("ManifestWorkloads", time.Since(start))
		}
	}

	if len(manifest.WorkloadDeployments) > 0 {
		logger.Debug("Applying workload deployments")
		start := time.Now()
		if err := rc.ApplyWorkloadDeployments(manifest.WorkloadDeployments); err != nil {
			return err
		}
		if GlobalTimingCallback != nil {
			GlobalTimingCallback("ManifestWorkloadDeployments", time.Since(start))
		}
	}

	if len(manifest.Firewalls) > 0 {
		logger.Debug("Applying firewalls")
		start := time.Now()
		if err := rc.ApplyFirewalls(manifest.Firewalls); err != nil {
			return err
		}
		if GlobalTimingCallback != nil {
			GlobalTimingCallback("ManifestFirewalls", time.Since(start))
		}
	}

	if len(manifest.Purge) > 0 {
		logger.Debug("Applying purge")
		start := time.Now()
		if err := rc.ApplyPurge(manifest.Purge); err != nil {
			return err
		}
		if GlobalTimingCallback != nil {
			GlobalTimingCallback("ManifestPurge", time.Since(start))
		}
	}

	// Hand the tracked ids to deleteResources, minus the ones consumed during
	// this run (a cache setting still referenced by a rule is not an orphan).
	CacheIds = rc.CacheIds
	RuleIds = rc.RuleIds

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

// unmarshalJsonArgs reads the args file referenced by azion.json. The recorded
// path is relative to the project root, so it is resolved against the working
// directory rather than the process' current directory.
//
// The boolean report tells whether the file was there at all. Instances are
// updated with PATCH, so a project without an args file must leave the field
// out of the request instead of sending an empty object, which would wipe
// arguments defined elsewhere.
func unmarshalJsonArgs(argsPath string) (map[string]interface{}, bool, error) {
	resolvedPath := argsPath
	if !filepath.IsAbs(resolvedPath) {
		workingDir, err := utils.GetWorkingDir()
		if err != nil {
			return nil, false, err
		}
		resolvedPath = filepath.Join(workingDir, argsPath)
	}

	marshalledArgs, err := os.ReadFile(resolvedPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// The args file is optional: a function without arguments is valid
			logger.Debug("Args file not found, no arguments will be sent",
				zap.String("path", resolvedPath))
			return map[string]interface{}{}, false, nil
		}
		logger.Debug("Error while reading args file",
			zap.String("path", resolvedPath), zap.Error(err))
		return nil, false, err
	}

	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		return nil, false, fmt.Errorf("%s: %w", msg.ErrorUnmarshalArgsFile, err)
	}
	return args, true, nil
}
