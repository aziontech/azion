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
		CacheIds := make(map[string]int64)
		CacheIdsBackup := make(map[string]int64)
		RuleIds := make(map[string]contracts.RuleIdsStruct)
		OriginKeys := make(map[string]string)
		OriginIds := make(map[string]int64)

		for _, cacheConf := range conf.CacheSettings {
			CacheIds[cacheConf.Name] = cacheConf.Id
		}

		for _, ruleConf := range conf.RulesEngine.Rules {
			RuleIds[ruleConf.Name] = contracts.RuleIdsStruct{
				Id:    ruleConf.Id,
				Phase: ruleConf.Phase,
			}
		}

		for _, originConf := range conf.Origin {
			OriginKeys[originConf.Name] = originConf.OriginKey
			OriginIds[originConf.Name] = originConf.OriginId
		}

		if manifestStructure.Domain.Name != "" {
			skip = true
			logger.Debug("", zap.Any("Domain Payload", manifestStructure.Domain))
			if conf.Domain.Id > 0 {
				msgf := fmt.Sprintf(msg.UpdateDomain, conf.Domain.Id, manifestStructure.Domain.Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			} else {
				msgf := fmt.Sprintf(msg.CreateDomain, manifestStructure.Domain.Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			}
		}

		for _, origin := range manifestStructure.Origins {
			logger.Debug("", zap.Any("Origin Payload", origin))
			if id := OriginIds[origin.Name]; id > 0 {
				msgf := fmt.Sprintf(msg.UpdateOrigin, id, OriginKeys[origin.Name], origin.Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			} else {
				msgf := fmt.Sprintf(msg.CreateOrigin, origin.Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			}
		}

		for _, cache := range manifestStructure.CacheSettings {
			logger.Debug("", zap.Any("Cache Setting Payload", cache))
			if id := CacheIds[*cache.Name]; id > 0 {
				msgf := fmt.Sprintf(msg.UpdateCacheSetting, id, cache.Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			} else {
				msgf := fmt.Sprintf(msg.CreateCacheSetting, cache.Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
			}
		}

		//backup cache ids
		for k, v := range CacheIds {
			CacheIdsBackup[k] = v
		}

		for _, rule := range manifestStructure.Rules {
			logger.Debug("", zap.Any("Rule Engine Payload", rule))
			if r := RuleIds[rule.Name]; r.Id > 0 {
				msgf := fmt.Sprintf(msg.UpdateRule, r.Id, rule.Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
				delete(RuleIds, rule.Name)
				for _, v := range rule.Behaviors {
					if v.RulesEngineBehaviorString != nil {
						if v.RulesEngineBehaviorString.Name == "set_cache_policy" {
							if id := CacheIdsBackup[v.RulesEngineBehaviorString.Target]; id > 0 {
								delete(CacheIds, v.RulesEngineBehaviorString.Target)
							} else if v.RulesEngineBehaviorString.Name == "set_origin" {
								if id := OriginIds[v.RulesEngineBehaviorString.Target]; id > 0 {
									delete(OriginKeys, v.RulesEngineBehaviorString.Target)
								}
							}
						}
					}
				}

			} else {
				msgf := fmt.Sprintf(msg.CreateRule, rule.Name)
				msgs = append(msgs, msgf)
				logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
				for _, v := range rule.Behaviors {
					if v.RulesEngineBehaviorString != nil {
						if v.RulesEngineBehaviorString.Name == "set_cache_policy" {
							if id := CacheIdsBackup[v.RulesEngineBehaviorString.Target]; id > 0 {
								delete(CacheIds, v.RulesEngineBehaviorString.Target)
							} else if v.RulesEngineBehaviorString.Name == "set_origin" {
								if id := OriginIds[v.RulesEngineBehaviorString.Target]; id > 0 {
									delete(OriginKeys, v.RulesEngineBehaviorString.Target)
								}
							}
						}

					}
				}
			}
		}

		for key, value := range RuleIds {
			msgf := fmt.Sprintf(msg.DeletingRuleEngine, value.Id, key)
			msgs = append(msgs, msgf)
			logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		}

		for key, value := range OriginKeys {
			if strings.Contains(key, "_single") {
				continue
			}
			msgf := fmt.Sprintf(msg.DeletingOrigin, OriginIds[key], value, key)
			msgs = append(msgs, msgf)
			logger.FInfoFlags(dry.Io.Out, msgf, dry.F.Format, dry.F.Out)
		}

		for key, value := range CacheIds {
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
