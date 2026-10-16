package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	"github.com/aziontech/azion-cli/pkg/cmd/build/v3"
	delete "github.com/aziontech/azion-cli/pkg/cmd/delete/v3/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmd/sync"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	manifestInt "github.com/aziontech/azion-cli/pkg/v3manifest"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeployCmd struct {
	Io                    *iostreams.IOStreams
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent   func(pathConfig string) (*contracts.AzionApplicationOptionsV3, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptionsV3, confConf string) error
	EnvLoader             func(path string) ([]string, error)
	BuildCmd              func(f *cmdutil.Factory) *build.BuildCmd
	Open                  func(name string) (*os.File, error)
	FilepathWalk          func(root string, fn filepath.WalkFunc) error
	F                     *cmdutil.Factory
	Unmarshal             func(data []byte, v interface{}) error
	Interpreter           func() *manifestInt.ManifestInterpreter
	VersionID             func() string
	Auto                  bool
	Env                   string
	NoPrompt              bool
	Path                  string
	ProjectConf           string
	SkipBuild             bool
	Sync                  bool
	AliasEnv              bool
}

func NewDeployCmd(f *cmdutil.Factory) *DeployCmd {
	return &DeployCmd{
		ProjectConf:           "azion",
		Env:                   ".edge/.env",
		Io:                    f.IOStreams,
		GetWorkDir:            utils.GetWorkingDir,
		FileReader:            os.ReadFile,
		WriteFile:             os.WriteFile,
		EnvLoader:             utils.LoadEnvVarsFromFile,
		BuildCmd:              build.NewBuildCmd,
		GetAzionJsonContent:   utils.GetAzionJsonContentV3,
		WriteAzionJsonContent: utils.WriteAzionJsonContentV3,
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
		Annotations: map[string]string{
			"Category": "skip",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy.Run(deploy.F)
		},
	}
	deployCmd.Flags().BoolP("help", "h", false, msg.DeployFlagHelp)
	deployCmd.Flags().StringVar(&deploy.Path, "path", "", msg.EdgeApplicationDeployPathFlag)
	deployCmd.Flags().BoolVar(&deploy.Auto, "auto", false, msg.DeployFlagAuto)
	deployCmd.Flags().BoolVar(&deploy.NoPrompt, "no-prompt", false, msg.DeployFlagNoPrompt)
	deployCmd.Flags().BoolVar(&deploy.SkipBuild, "skip-build", false, msg.DeployFlagSkipBuild)
	deployCmd.Flags().StringVar(&deploy.ProjectConf, "config-dir", "azion", msg.EdgeApplicationDeployProjectConfFlag)
	deployCmd.Flags().BoolVar(&deploy.Sync, "sync", false, msg.EdgeApplicationDeploySync)
	deployCmd.Flags().StringVar(&deploy.Env, "env", ".edge/.env", msg.EnvFlag)
	deployCmd.Flags().BoolVar(&deploy.AliasEnv, "alias-env", false, msg.AliasEnvFlag)
	return deployCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeployCmd(f))
}

func (cmd *DeployCmd) ExternalRun(f *cmdutil.Factory, configPath string, env string, shouldSync, auto, skipBuild, aliasEnv bool) error {
	cmd.AliasEnv = aliasEnv
	cmd.ProjectConf = configPath
	cmd.Sync = shouldSync
	cmd.Env = env
	cmd.Auto = auto
	cmd.SkipBuild = skipBuild
	return cmd.Run(f)
}

func (cmd *DeployCmd) Run(f *cmdutil.Factory) error {
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

	if !cmd.SkipBuild {
		buildCmd := cmd.BuildCmd(f)
		err := buildCmd.ExternalRun(&contracts.BuildInfoV3{AliasEnv: cmd.AliasEnv}, cmd.ProjectConf, &msgs)
		if err != nil {
			logger.Debug("Error while running build command called by deploy command", zap.Error(err))
			return err
		}
	}

	conf, err := cmd.GetAzionJsonContent(cmd.ProjectConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	if conf.Prefix == "" || conf.RotatePrefix == nil || *conf.RotatePrefix == true {
		conf.Prefix = cmd.VersionID()
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
	}

	if _, err := pipeline.Run(ctx, V3DeployPipeline{}, state, &pipeline.Plan{Factory: f}); err != nil {
		return err
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

func callDeleteCascade(f *cmdutil.Factory) error {
	del := delete.NewDeleteCmd(f)
	del.Io.In = f.IOStreams.In
	del.Io.Out = f.IOStreams.Out
	del.Io.Err = f.IOStreams.Err
	del.Io = f.IOStreams

	cmd := delete.NewCobraCmd(del)

	cmd.SetArgs([]string{"--cascade"})
	_, err := cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}
