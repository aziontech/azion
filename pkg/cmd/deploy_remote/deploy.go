package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

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
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeployCmd struct {
	Io                       *iostreams.IOStreams
	GetWorkDir               func() (string, error)
	FileReader               func(path string) ([]byte, error)
	WriteFile                func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent      func(pathConfig string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent    func(conf *contracts.AzionApplicationOptions, confConf string) error
	commandRunInteractive    func(f *cmdutil.Factory, comm string) error
	commandRunnerOutput      func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
	WriteManifest            func(manifest *contracts.ManifestV4, pathMan string) error
	EnvLoader                func(path string) ([]string, error)
	BuildCmd                 func(f *cmdutil.Factory) *build.BuildCmd
	Open                     func(name string) (*os.File, error)
	FilepathWalk             func(root string, fn filepath.WalkFunc) error
	F                        *cmdutil.Factory
	Unmarshal                func(data []byte, v interface{}) error
	Interpreter              func() *manifestInt.ManifestInterpreter
	VersionID                func() string
	ReadSettings             func(path string) (token.Settings, error)
	GetCredentialsForBucket  func(path string, bucketName string) (token.S3Credentials, bool, error)
	SaveCredentialsForBucket func(path string, bucketName string, creds token.S3Credentials) error
	CreateBucketCredentials  func(ctx context.Context, bucketName string, f *cmdutil.Factory, subdir string) (token.S3Credentials, error)
	Auto                     bool
	Env                      string
	NoPrompt                 bool
	Path                     string
	ProjectConf              string
	Sync                     bool
	SkipBuild                bool
	SkipFramework            bool
	WriteBucket              bool
	Workers                  int
	AliasEnv                 bool
}

func NewDeployCmd(f *cmdutil.Factory) *DeployCmd {
	return &DeployCmd{
		ProjectConf:              "azion",
		Env:                      ".edge/.env",
		Io:                       f.IOStreams,
		GetWorkDir:               utils.GetWorkingDir,
		FileReader:               os.ReadFile,
		WriteFile:                os.WriteFile,
		EnvLoader:                utils.LoadEnvVarsFromFile,
		BuildCmd:                 build.NewBuildCmd,
		GetAzionJsonContent:      utils.GetAzionJsonContent,
		WriteAzionJsonContent:    utils.WriteAzionJsonContent,
		commandRunInteractive:    command.CommandRunInteractive,
		commandRunnerOutput:      command.CommandRunInteractiveWithOutput,
		WriteManifest:            WriteManifest,
		Open:                     os.Open,
		FilepathWalk:             filepath.Walk,
		Unmarshal:                json.Unmarshal,
		F:                        f,
		Interpreter:              manifestInt.NewManifestInterpreter,
		VersionID:                utils.Timestamp,
		ReadSettings:             token.ReadSettings,
		GetCredentialsForBucket:  token.GetCredentialsForBucket,
		SaveCredentialsForBucket: token.SaveCredentialsForBucket,
		CreateBucketCredentials:  CreateBucketCredentials,
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
	deployCmd.Flags().StringVar(&deploy.Path, "path", "", msg.EdgeApplicationDeployPathFlag)
	deployCmd.Flags().BoolVar(&deploy.Auto, "auto", false, msg.DeployFlagAuto)
	deployCmd.Flags().BoolVar(&deploy.NoPrompt, "no-prompt", false, msg.DeployFlagNoPrompt)
	deployCmd.Flags().StringVar(&deploy.ProjectConf, "config-dir", "azion", msg.EdgeApplicationDeployProjectConfFlag)
	deployCmd.Flags().BoolVar(&deploy.Sync, "sync", false, msg.EdgeApplicationDeploySync)
	deployCmd.Flags().StringVar(&deploy.Env, "env", ".edge/.env", msg.EnvFlag)
	deployCmd.Flags().BoolVar(&deploy.AliasEnv, "alias-env", false, msg.AliasEnvFlag)
	return deployCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeployCmd(f))
}

func (cmd *DeployCmd) ExternalRun(f *cmdutil.Factory, configPath string, env string, shouldSync, auto, skipBuild, writeBucket, skipFramework, aliasEnv bool, workers int) error {
	cmd.AliasEnv = aliasEnv
	cmd.ProjectConf = configPath
	cmd.Sync = shouldSync
	cmd.Env = env
	cmd.Auto = auto
	cmd.SkipBuild = skipBuild
	cmd.SkipFramework = skipFramework
	cmd.WriteBucket = writeBucket
	cmd.Workers = workers
	return cmd.Run(f)
}

func (cmd *DeployCmd) Run(f *cmdutil.Factory) error {
	InitTimingSummary()
	totalStart := time.Now()

	// Set up callback for manifest package timing
	manifestInt.GlobalTimingCallback = HandleManifestTimingCallback

	msgs := []string{}
	logger.FInfoFlags(cmd.F.IOStreams.Out, "Running deploy command\n", cmd.F.Format, cmd.F.Out)
	msgs = append(msgs, "Running deploy command")
	ctx := context.Background()

	if cmd.Sync {
		syncCmd := sync.NewSyncCmd(f)
		syncCmd.ProjectConf = cmd.ProjectConf
		syncCmd.EnvPath = cmd.Env
		if err := sync.Run(syncCmd); err != nil {
			logger.Debug("Error while synchronizing local resources with remove resources", zap.Error(err))
			return err
		}
	}

	conf, err := cmd.GetAzionJsonContent(cmd.ProjectConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	defer func() {
		if err := cmd.WriteAzionJsonContent(conf, cmd.ProjectConf); err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
		}
	}()

	var oldprefix, newprefix string

	if conf.Prefix == "" || conf.RotatePrefix == nil || *conf.RotatePrefix == true {
		versionID := cmd.VersionID()
		oldprefix, newprefix = conf.Prefix, versionID
	} else {
		oldprefix, newprefix = conf.Prefix, conf.Prefix
	}

	err = checkArgsJson(cmd, cmd.ProjectConf)
	if err != nil {
		return err
	}

	clients := NewClients(f)
	interpreter := cmd.Interpreter()

	state := &DeployState{
		Cmd: cmd, F: f, Ctx: ctx, Conf: conf, Msgs: &msgs,
		Clients:     clients,
		Interpreter: interpreter,
		OldPrefix:   oldprefix, NewPrefix: newprefix,
	}

	if _, err := pipeline.Run(ctx, V4DeployPipeline{}, state, &pipeline.Plan{
		Factory: f,
		Timing:  HandleDeployTimingCallback,
	}); err != nil {
		return err
	}

	// Calculate total deploy time
	GlobalTimingSummary.TotalDeployTime = time.Since(totalStart)

	// Print timing summary before success messages
	GlobalTimingSummary.PrintSummary()

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
