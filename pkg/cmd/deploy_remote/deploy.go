package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmd/sync"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/command"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/output"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeployCmd struct {
	Io                    *iostreams.IOStreams
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent   func(pathConfig string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confConf string) error
	commandRunInteractive func(f *cmdutil.Factory, comm string) error
	commandRunnerOutput   func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
	WriteManifest         func(manifest *contracts.ManifestV4, pathMan string) error
	EnvLoader             func(path string) ([]string, error)
	BuildCmd              func(f *cmdutil.Factory) *build.BuildCmd
	Open                  func(name string) (*os.File, error)
	FilepathWalk          func(root string, fn filepath.WalkFunc) error
	F                     *cmdutil.Factory
	Unmarshal             func(data []byte, v interface{}) error
	Interpreter           func() *manifestInt.ManifestInterpreter
	VersionID             func() string
}

var (
	Path        string
	Auto        bool
	NoPrompt    bool
	SkipBuild   bool
  SkipFramework bool
	ProjectConf string
	Sync        bool
	Env         string
	FunctionIds map[string]contracts.AzionJsonDataFunction
  WriteBucket   bool
)

func NewDeployCmd(f *cmdutil.Factory) *DeployCmd {
	return &DeployCmd{
		Io:                    f.IOStreams,
		GetWorkDir:            utils.GetWorkingDir,
		FileReader:            os.ReadFile,
		WriteFile:             os.WriteFile,
		EnvLoader:             utils.LoadEnvVarsFromFile,
		BuildCmd:              build.NewBuildCmd,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		commandRunInteractive: command.CommandRunInteractive,
		commandRunnerOutput:   command.CommandRunInteractiveWithOutput,
		WriteManifest:         WriteManifest,
		Open:                  os.Open,
		FilepathWalk:          filepath.Walk,
		Unmarshal:             json.Unmarshal,
		F:                     f,
		Interpreter:           manifestInt.NewManifestInterpreter,
		VersionID:             utils.Timestamp,
	}
}

func NewCobraCmd(deploy *DeployCmd) *cobra.Command {
	deployCmd := &cobra.Command{
		Use:           msg.DeployUsage,
		Short:         msg.DeployShortDescription,
		Long:          msg.DeployLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy.Run(deploy.F)
		},
	}
	deployCmd.Flags().BoolP("help", "h", false, msg.DeployFlagHelp)
	deployCmd.Flags().StringVar(&Path, "path", "", msg.EdgeApplicationDeployPathFlag)
	deployCmd.Flags().BoolVar(&Auto, "auto", false, msg.DeployFlagAuto)
	deployCmd.Flags().BoolVar(&NoPrompt, "no-prompt", false, msg.DeployFlagNoPrompt)
	deployCmd.Flags().BoolVar(&SkipBuild, "skip-build", false, msg.DeployFlagSkipBuild)
	deployCmd.Flags().StringVar(&ProjectConf, "config-dir", "azion", msg.EdgeApplicationDeployProjectConfFlag)
	deployCmd.Flags().BoolVar(&Sync, "sync", false, msg.EdgeApplicationDeploySync)
	deployCmd.Flags().StringVar(&Env, "env", ".edge/.env", msg.EnvFlag)
	return deployCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeployCmd(f))
}

func (cmd *DeployCmd) ExternalRun(f *cmdutil.Factory, configPath string, env string, shouldSync, auto, skipBuild, writeBucket, skipFramework bool) error {
	ProjectConf = configPath
	Sync = shouldSync
	Env = env
	Auto = auto
	SkipBuild = skipBuild
	SkipFramework = skipFramework
	WriteBucket = writeBucket
	return cmd.Run(f)
}

func (cmd *DeployCmd) Run(f *cmdutil.Factory) error {
	msgs := []string{}
	logger.FInfoFlags(cmd.F.IOStreams.Out, "Running deploy command\n", cmd.F.Format, cmd.F.Out)
	msgs = append(msgs, "Running deploy command")
	ctx := context.Background()

	if Sync {
		sync.ProjectConf = ProjectConf
		syncCmd := sync.NewSyncCmd(f)
		syncCmd.EnvPath = Env
		if err := sync.Run(syncCmd); err != nil {
			logger.Debug("Error while synchronizing local resources with remove resources", zap.Error(err))
			return err
		}
	}

	if !SkipBuild {
		buildCmd := cmd.BuildCmd(f)
		err := buildCmd.ExternalRun(&contracts.BuildInfo{}, ProjectConf, &msgs, SkipFramework)
		if err != nil {
			logger.Debug("Error while running build command called by deploy command", zap.Error(err))
			return err
		}
	}

	conf, err := cmd.GetAzionJsonContent(ProjectConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	defer func() {
		if err := cmd.WriteAzionJsonContent(conf, ProjectConf); err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
		}
	}()

	versionID := cmd.VersionID()
	var oldprefix string

	oldprefix, conf.Prefix = conf.Prefix, versionID

	err = checkArgsJson(cmd, ProjectConf)
	if err != nil {
		return err
	}

	clients := NewClients(f)
	interpreter := cmd.Interpreter()

	if !SkipBuild && conf.NotFirstRun {
		cmdStr := fmt.Sprintf("config replace -k '%s' -v '%s'", oldprefix, conf.Prefix)
		vul := vulcanPkg.NewVulcan()
		command := vul.Command("", cmdStr, cmd.F)
		logger.Debug("Running the following command", zap.Any("Command", command))

		err := cmd.commandRunInteractive(cmd.F, command)
		if err != nil {
			return err
		}
		buildCmd := cmd.BuildCmd(f)
		err = buildCmd.ExternalRun(&contracts.BuildInfo{Preset: conf.Preset}, ProjectConf, &msgs)
		if err != nil {
			logger.Debug("Error while running build command called by deploy command", zap.Error(err))
			return err
		}
	}

	pathManifest, err := interpreter.ManifestPath()
	if err != nil {
		return err
	}

	err = cmd.doApplication(clients.EdgeApplication, context.Background(), conf, &msgs)
	if err != nil {
		return err
	}

	manifestStructure, err := interpreter.ReadManifest(pathManifest, f, &msgs)
	if err != nil {
		return err
	}

	// Check if directory exists; if not, we skip creating bucket
	if len(manifestStructure.EdgeStorage) == 0 {
		logger.Debug(msg.SkipBucket)
	} else {
		err = cmd.doBucket(clients.Bucket, ctx, conf, &msgs)
		if err != nil {
			return err
		}
	}

	if !conf.NotFirstRun {
		err = cmd.callBundlerInit(conf)
		if err != nil {
			return err
		}
		buildCmd := cmd.BuildCmd(f)
		err = buildCmd.ExternalRun(&contracts.BuildInfo{}, ProjectConf, &msgs)
		if err != nil {
			logger.Debug("Error while running build command called by deploy command", zap.Error(err))
			return err
		}
	}

	manifestStructure, err = interpreter.ReadManifest(pathManifest, f, &msgs)
	if err != nil {
		return err
	}

	// Check if directory exists; if not, we skip uploading static files
	if _, err := os.Stat(PathStatic); os.IsNotExist(err) {
		logger.Debug(msg.SkipUpload)
	} else {
		for _, storage := range manifestStructure.EdgeStorage {
			err = cmd.uploadFiles(f, conf, &msgs, storage.Dir)
			if err != nil {
				return err
			}
		}
	}

	if len(conf.RulesEngine.Rules) == 0 && !conf.NotFirstRun {
		err = cmd.doRulesDeploy(ctx, conf, clients.EdgeApplication, &msgs)
		if err != nil {
			return err
		}
	}

	conf.NotFirstRun = true

	err = interpreter.CreateResources(conf, manifestStructure, FunctionIds, f, ProjectConf, &msgs)
	if err != nil {
		return err
	}

	if len(manifestStructure.Workloads) == 0 || manifestStructure.Workloads[0].Name == "" {
		err = cmd.doWorkload(clients.Workload, ctx, conf, &msgs)
		if err != nil {
			return err
		}
	}

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.DeploySuccessful, f.Format, f.Out)
	msgs = append(msgs, msg.DeploySuccessful)

	msgfOutputDomainSuccess := fmt.Sprintf(msg.DeployOutputWorkloadSuccess, conf.Workloads.Url)
	logger.FInfoFlags(cmd.F.IOStreams.Out, msgfOutputDomainSuccess, f.Format, f.Out)
	msgs = append(msgs, msgfOutputDomainSuccess)

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.DeployPropagation, f.Format, f.Out)
	msgs = append(msgs, msg.DeployPropagation)

	outSlice := output.SliceOutput{
		Messages: msgs,
		GeneralOutput: output.GeneralOutput{
			Out:   cmd.F.IOStreams.Out,
			Flags: cmd.F.Flags,
		},
	}

	return output.Print(&outSlice)
}
