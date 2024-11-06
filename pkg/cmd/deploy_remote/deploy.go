package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	apiEdgeApplications "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmd/sync"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
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
	ProjectConf string
	Sync        bool
	Env         string
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

func (cmd *DeployCmd) ExternalRun(f *cmdutil.Factory, configPath string, env string, shouldSync, auto, skipBuild bool) error {
	ProjectConf = configPath
	Sync = shouldSync
	Env = env
	Auto = auto
	SkipBuild = skipBuild
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
		err := buildCmd.ExternalRun(&contracts.BuildInfo{}, ProjectConf, &msgs)
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

	versionID := cmd.VersionID()

	conf.Prefix = versionID

	err = checkArgsJson(cmd, ProjectConf)
	if err != nil {
		return err
	}

	clients := NewClients(f)
	interpreter := cmd.Interpreter()

	pathManifest, err := interpreter.ManifestPath()
	if err != nil {
		return err
	}

	err = cmd.doApplication(clients.EdgeApplication, context.Background(), conf, &msgs)
	if err != nil {
		return err
	}

	singleOriginId, err := cmd.doOriginSingle(clients.Origin, ctx, conf, &msgs)
	if err != nil {
		return err
	}

	err = cmd.doBucket(clients.Bucket, ctx, conf, &msgs)
	if err != nil {
		return err
	}

	// Check if directory exists; if not, we skip uploading static files
	if _, err := os.Stat(PathStatic); os.IsNotExist(err) {
		logger.Debug(msg.SkipUpload)
	} else {
		err = cmd.uploadFiles(f, conf, &msgs)
		if err != nil {
			return err
		}
	}

	conf.Function.File = ".edge/worker.js"
	err = cmd.doFunction(clients, ctx, conf, &msgs)
	if err != nil {
		return err
	}

	if !conf.NotFirstRun {
		ruleDefaultID, err := clients.EdgeApplication.GetRulesDefault(ctx, conf.Application.ID, "request")
		if err != nil {
			logger.Debug("Error while getting default rules engine", zap.Error(err))
			return err
		}
		behaviors := make([]sdk.RulesEngineBehaviorEntry, 0)

		var behString sdk.RulesEngineBehaviorString
		behString.SetName("set_origin")

		behString.SetTarget(strconv.Itoa(int(singleOriginId)))

		behaviors = append(behaviors, sdk.RulesEngineBehaviorEntry{
			RulesEngineBehaviorString: &behString,
		})

		reqUpdateRulesEngine := apiEdgeApplications.UpdateRulesEngineRequest{
			IdApplication: conf.Application.ID,
			Phase:         "request",
			Id:            ruleDefaultID,
		}

		reqUpdateRulesEngine.SetBehaviors(behaviors)

		_, err = clients.EdgeApplication.UpdateRulesEngine(ctx, &reqUpdateRulesEngine)
		if err != nil {
			logger.Debug("Error while updating default rules engine", zap.Error(err))
			return err
		}
	}

	manifestStructure, err := interpreter.ReadManifest(pathManifest, f, &msgs)
	if err != nil {
		return err
	}

	if len(conf.RulesEngine.Rules) == 0 {
		err = cmd.doRulesDeploy(ctx, conf, clients.EdgeApplication, &msgs)
		if err != nil {
			return err
		}
	}

	err = interpreter.CreateResources(conf, manifestStructure, f, ProjectConf, &msgs)
	if err != nil {
		return err
	}

	if manifestStructure.Domain.Name == "" {
		err = cmd.doDomain(clients.Domain, ctx, conf, &msgs)
		if err != nil {
			return err
		}
	}

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.DeploySuccessful, f.Format, f.Out)
	msgs = append(msgs, msg.DeploySuccessful)

	msgfOutputDomainSuccess := fmt.Sprintf(msg.DeployOutputDomainSuccess, conf.Domain.Url)
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
