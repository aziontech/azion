package dryrun

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/dryrun"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"go.uber.org/zap"
)

//TODO: FIX EVERYTHING

type DryrunStruct struct {
	Io                    *iostreams.IOStreams
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent   func(pathConfig string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confConf string) error
	EnvLoader             func(path string) ([]string, error)
	BuildCmd              func(f *cmdutil.Factory) *build.BuildCmd
	Open                  func(name string) (*os.File, error)
	FilepathWalk          func(root string, fn filepath.WalkFunc) error
	F                     *cmdutil.Factory
	Unmarshal             func(data []byte, v interface{}) error
	Interpreter           func() *manifestInt.ManifestInterpreter
	VersionID             func() string
	Stat                  func(name string) (fs.FileInfo, error)
}

var skip bool

func NewDryrunCmd(f *cmdutil.Factory) *DryrunStruct {
	return &DryrunStruct{
		Io:                    f.IOStreams,
		GetWorkDir:            utils.GetWorkingDir,
		FileReader:            os.ReadFile,
		WriteFile:             os.WriteFile,
		EnvLoader:             utils.LoadEnvVarsFromFile,
		BuildCmd:              build.NewBuildCmd,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		Open:                  os.Open,
		FilepathWalk:          filepath.Walk,
		Unmarshal:             json.Unmarshal,
		F:                     f,
		Interpreter:           manifestInt.NewManifestInterpreter,
		VersionID:             utils.Timestamp,
		Stat:                  os.Stat,
	}
}

func (dry *DryrunStruct) SimulateDeploy(workingDir, projConf string) error {
	msgs := []string{}
	conf, err := dry.GetAzionJsonContent(projConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	if conf.Application.ID == 0 {
		msgf := fmt.Sprintf(msg.CreateEdgeApp, conf.Name)
		msgs = append(msgs, msgf)
		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
	} else {
		msgf := fmt.Sprintf(msg.UpdateEdgeApp, conf.Application.ID, conf.Name)
		msgs = append(msgs, msgf)
		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
	}

	if !conf.NotFirstRun {
		msgf := fmt.Sprintf(msg.CreateOriginSingle, utils.Concat(conf.Name, "_single"))
		msgs = append(msgs, msgf)
		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
	}

	if conf.Bucket == "" || !(conf.Preset == "javascript" || conf.Preset == "typescript") {
		msgf := fmt.Sprintf(msg.CreateBucket, conf.Name)
		msgs = append(msgs, msgf)
		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
	}

	interpreter := dry.Interpreter()

	pathManifest, err := interpreter.ManifestPath()
	if err != nil {
		return err
	}

	//FUNCTION AND INSTANCE

	if !conf.NotFirstRun {
		msgf := fmt.Sprintf(msg.UpdateDefaultRule, utils.Concat(conf.Name, "_single"))
		msgs = append(msgs, msgf)
		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
	}

	if len(conf.RulesEngine.Rules) == 0 && !conf.NotFirstRun {
		msgs = append(msgs, msg.CreateRulesCache)
		logger.Debug("", zap.Any("Cache Setting information", msg.AskCreateCacheSettings))
		logger.FInfoFlags(dry.Io.Out, msg.CreateRulesCache, dry.F.Format, dry.F.Out)
	}

	var skipManifest bool
	manifestStructure, err := interpreter.ReadManifest(pathManifest, dry.F, &msgs)
	if err != nil {
		skipManifest = true
	}

	if _, err := dry.Stat(pathManifest); os.IsNotExist(err) {
		logger.FInfoFlags(dry.Io.Out, msg.SkipManifest, dry.F.Format, dry.F.Out)
		msgs = append(msgs, msg.SkipManifest)
	} else if !skipManifest {
		// Initialize maps to track resources
		manifestInt.CacheIds = make(map[string]int64)
		manifestInt.CacheIdsBackup = make(map[string]int64)
		manifestInt.RuleIds = make(map[string]contracts.RuleIdsStruct)
		manifestInt.ConnectorIds = make(map[string]int64)
		manifestInt.DeploymentIds = make(map[string]int64)
		manifestInt.FunctionIds = make(map[string]contracts.AzionJsonDataFunction)

		// Local maps for origin tracking (not in manifest package)
		OriginKeys := make(map[string]string)
		OriginIds := make(map[string]int64)

		for _, cacheConf := range conf.CacheSettings {
			manifestInt.CacheIds[cacheConf.Name] = cacheConf.Id
		}

		for _, ruleConf := range conf.RulesEngine.Rules {
			manifestInt.RuleIds[ruleConf.Name] = contracts.RuleIdsStruct{
				Id:    ruleConf.Id,
				Phase: ruleConf.Phase,
			}
		}

		for _, funcConf := range conf.Function {
			manifestInt.FunctionIds[funcConf.Name] = funcConf
		}

		for _, deploymentConf := range conf.Workloads.Deployments {
			manifestInt.DeploymentIds[deploymentConf.Name] = deploymentConf.Id
		}

		for _, originConf := range conf.Origin {
			OriginKeys[originConf.Name] = originConf.OriginKey
			OriginIds[originConf.Name] = originConf.OriginId
			// Also initialize ConnectorIds with the same data
			manifestInt.ConnectorIds[originConf.Name] = originConf.OriginId
		}

		if len(manifestStructure.Workloads) > 0 && manifestStructure.Workloads[0].Name != "" {
			skip = true
			logger.Debug("", zap.Any("Workload Payload", manifestStructure.Workloads))
			if conf.Domain.Id > 0 {
				msgf := fmt.Sprintf(msg.UpdateDomain, conf.Domain.Id, manifestStructure.Workloads[0].Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			} else {
				msgf := fmt.Sprintf(msg.CreateDomain, manifestStructure.Workloads[0].Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			}
		}

		// for _, connector := range manifestStructure.EdgeConnectors {
		// 	logger.Debug("", zap.Any("Edge Connector Payload", connector))
		// 	if id, ok := manifestInt.ConnectorIds[connector.Name]; ok && id > 0 {
		// 		msgf := fmt.Sprintf(msg.UpdateOrigin, id, OriginKeys[connector.Name], connector.Name)
		// 		msgs = append(msgs, msgf)
		// 		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		// 	} else {
		// 		msgf := fmt.Sprintf(msg.CreateOrigin, connector.Name)
		// 		msgs = append(msgs, msgf)
		// 		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		// 	}
		// }

		// for _, app := range manifestStructure.EdgeApplications {
		// 	for _, cache := range app.Cache {
		// 		logger.Debug("", zap.Any("Cache Setting Payload", cache))
		// 		if id, ok := manifestInt.CacheIds[cache.Name]; ok && id > 0 {
		// 			msgf := fmt.Sprintf(msg.UpdateCacheSetting, id, cache.Name)
		// 			msgs = append(msgs, msgf)
		// 			logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		// 		} else {
		// 			msgf := fmt.Sprintf(msg.CreateCacheSetting, cache.Name)
		// 			msgs = append(msgs, msgf)
		// 			logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		// 		}
		// 	}
		// }

		//backup cache ids
		for k, v := range manifestInt.CacheIds {
			manifestInt.CacheIdsBackup[k] = v
		}

		// Process rules from manifest
		// for _, rule := range manifestStructure.Rules {
		// 	logger.Debug("", zap.Any("Rule Engine Payload", rule))
		// 	if r, ok := manifestInt.RuleIds[rule.Name]; ok && r.Id > 0 {
		// 		msgf := fmt.Sprintf(msg.UpdateRule, r.Id, rule.Name)
		// 		msgs = append(msgs, msgf)
		// 		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		// 		delete(manifestInt.RuleIds, rule.Name)
		// 		for _, v := range rule.Behaviors {
		// 			if v.RulesEngineBehaviorString != nil {
		// 				if v.RulesEngineBehaviorString.Name == "set_cache_policy" {
		// 					if id, ok := manifestInt.CacheIdsBackup[v.RulesEngineBehaviorString.Target]; ok && id > 0 {
		// 						delete(manifestInt.CacheIds, v.RulesEngineBehaviorString.Target)
		// 					}
		// 				} else if v.RulesEngineBehaviorString.Name == "set_origin" {
		// 					if id, ok := OriginIds[v.RulesEngineBehaviorString.Target]; ok && id > 0 {
		// 						delete(OriginKeys, v.RulesEngineBehaviorString.Target)
		// 					}
		// 				}
		// 			}
		// 		}
		// 	} else {
		// 		msgf := fmt.Sprintf(msg.CreateRule, rule.Name)
		// 		msgs = append(msgs, msgf)
		// 		logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		// 		for _, v := range rule.Behaviors {
		// 			if v.RulesEngineBehaviorString != nil {
		// 				if v.RulesEngineBehaviorString.Name == "set_cache_policy" {
		// 					if id, ok := manifestInt.CacheIdsBackup[v.RulesEngineBehaviorString.Target]; ok && id > 0 {
		// 						delete(manifestInt.CacheIds, v.RulesEngineBehaviorString.Target)
		// 					}
		// 				} else if v.RulesEngineBehaviorString.Name == "set_origin" {
		// 					if id, ok := OriginIds[v.RulesEngineBehaviorString.Target]; ok && id > 0 {
		// 						delete(OriginKeys, v.RulesEngineBehaviorString.Target)
		// 					}
		// 				}
		// 			}
		// 		}
		// 	}
		// }

		for key, value := range manifestInt.RuleIds {
			msgf := fmt.Sprintf(msg.DeletingRuleEngine, value.Id, key)
			msgs = append(msgs, msgf)
			logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		}

		for key, value := range OriginKeys {
			if strings.Contains(key, "_single") {
				continue
			}
			if id, ok := OriginIds[key]; ok {
				msgf := fmt.Sprintf(msg.DeletingOrigin, id, value, key)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			}
		}

		for key, value := range manifestInt.CacheIds {
			msgf := fmt.Sprintf(msg.DeletingCacheSetting, value, key)
			msgs = append(msgs, msgf)
			logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		}

	}

	if !skip {
		if conf.Domain.Id > 0 {
			msgf := fmt.Sprintf(msg.UpdateDomain, conf.Domain.Id, conf.Name)
			msgs = append(msgs, msgf)
			logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		} else {
			msgf := fmt.Sprintf(msg.CreateDomain, conf.Name)
			msgs = append(msgs, msgf)
			logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		}
	}

	outSlice := output.SliceOutput{
		Messages: msgs,
		GeneralOutput: output.GeneralOutput{
			Out:   dry.F.IOStreams.Out,
			Flags: dry.F.Flags,
		},
	}

	return output.Print(&outSlice)
}
