package init

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/cmd/deploy"
	"github.com/aziontech/azion-cli/pkg/cmd/dev"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	REPO string = "https://github.com/aziontech/azioncli-template.git"
)

type InitInfo struct {
	Name           string
	Template       string
	Mode           string
	PathWorkingDir string
	GlobalFlagAll  bool
}

type InitCmd struct {
	F                     *cmdutil.Factory
	Io                    *iostreams.IOStreams
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	LookPath              func(bin string) (string, error)
	IsDirEmpty            func(dirpath string) (bool, error)
	CleanDir              func(dirpath string) error
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	OpenFile              func(name string) (*os.File, error)
	RemoveAll             func(path string) error
	Rename                func(oldpath string, newpath string) error
	CreateTempDir         func(dir string, pattern string) (string, error)
	EnvLoader             func(path string) ([]string, error)
	Stat                  func(path string) (fs.FileInfo, error)
	Mkdir                 func(path string, perm os.FileMode) error
	GitPlainClone         func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	CommandRunner         func(cmd string, envvars []string) (string, int, error)
	CommandRunInteractive func(f *cmdutil.Factory, comm string) error
	ShouldDevDeploy       func(info *InitInfo, msg string) (bool, error)
	DeployCmd             func(f *cmdutil.Factory) *deploy.DeployCmd
	DevCmd                func(f *cmdutil.Factory) *dev.DevCmd
	ChangeDir             func(dir string) error
}

func NewInitCmd(f *cmdutil.Factory) *InitCmd {
	return &InitCmd{
		F:               f,
		Io:              f.IOStreams,
		GetWorkDir:      utils.GetWorkingDir,
		FileReader:      os.ReadFile,
		LookPath:        exec.LookPath,
		IsDirEmpty:      utils.IsDirEmpty,
		CleanDir:        utils.CleanDirectory,
		WriteFile:       os.WriteFile,
		OpenFile:        os.Open,
		RemoveAll:       os.RemoveAll,
		Rename:          os.Rename,
		CreateTempDir:   os.MkdirTemp,
		EnvLoader:       utils.LoadEnvVarsFromFile,
		Stat:            os.Stat,
		Mkdir:           os.MkdirAll,
		GitPlainClone:   git.PlainClone,
		ShouldDevDeploy: shouldDevDeploy,
		DevCmd:          dev.NewDevCmd,
		DeployCmd:       deploy.NewDeployCmd,
		ChangeDir:       os.Chdir,
		CommandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
		CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
			return utils.CommandRunInteractive(f, comm)
		},
	}
}

func NewCobraCmd(init *InitCmd, f *cmdutil.Factory) *cobra.Command {
	info := &InitInfo{}
	cobraCmd := &cobra.Command{
		Use:           msg.EdgeApplicationsInitUsage,
		Short:         msg.EdgeApplicationsInitShortDescription,
		Long:          msg.EdgeApplicationsInitLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion init
		$ azion init --help
		$ azion init --name testproject
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			info.GlobalFlagAll = f.GlobalFlagAll
			return init.Run(info)
		},
	}

	cobraCmd.Flags().StringVar(&info.Name, "name", "", msg.EdgeApplicationsInitFlagName)
	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewInitCmd(f), f)
}

func (cmd *InitCmd) Run(info *InitInfo) error {
	logger.Debug("Running init command")

	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}
	info.PathWorkingDir = path

	// Checks for global --yes flag and that name flag was not sent
	if info.GlobalFlagAll && info.Name == "" {
		info.Name = thoth.GenerateName()
	} else {
		// if name was not sent we ask for input, otherwise info.Name already has the value
		if info.Name == "" {

			projName, err := askForInput(msg.InitProjectQuestion, thoth.GenerateName())
			if err != nil {
				return err
			}

			info.Name = projName
		}
	}

	info.PathWorkingDir = info.PathWorkingDir + "/" + info.Name

	err = cmd.selectVulcanTemplates(info)
	if err != nil {
		return err
	}

	if err = cmd.createTemplateAzion(info); err != nil {
		return err
	}

	logger.FInfo(cmd.Io.Out, msg.WebAppInitCmdSuccess)

	err = cmd.ChangeDir(info.PathWorkingDir)
	if err != nil {
		logger.Debug("Error while changing to new working directory", zap.Error(err))
		return msg.ErrorWorkingDir
	}

	shouldDev, err := cmd.ShouldDevDeploy(info, "Do you want to start a local development server?")
	if err != err {
		return err
	}
	if shouldDev {
		shouldDeps, err := cmd.ShouldDevDeploy(info, "Do you want to install project dependencies? This may be required to start local development server")
		if err != err {
			return err
		}

		if shouldDeps {
			answer, err := utils.GetPackageManager()
			if err != nil {
				return err
			}
			err = depsInstall(cmd, answer)
			if err != nil {
				logger.Debug("Error while installing project dependencies", zap.Error(err))
				return msg.ErrorDeps
			}
		}

		logger.Debug("Running dev command from init command")
		dev := cmd.DevCmd(cmd.F)
		err = dev.Run(cmd.F)
		if err != nil {
			logger.Debug("Error while running dev command called by init command", zap.Error(err))
			return err
		}
	} else {
		logger.FInfo(cmd.Io.Out, msg.InitDevCommand)
	}

	shouldDeploy, err := cmd.ShouldDevDeploy(info, "Do you want to deploy your project?")
	if err != err {
		return err
	}
	if shouldDeploy {
		shouldDeps, err := cmd.ShouldDevDeploy(info, "Do you want to install project dependencies? This may be required to deploy your project")
		if err != err {
			return err
		}

		if shouldDeps {
			answer, err := utils.GetPackageManager()
			if err != nil {
				return err
			}
			err = depsInstall(cmd, answer)
			if err != nil {
				logger.Debug("Failed to install project dependencies")
				return err
			}
		}

		logger.Debug("Running deploy command from init command")
		deploy := cmd.DeployCmd(cmd.F)
		err = deploy.Run(cmd.F)
		if err != nil {
			logger.Debug("Error while running deploy command called by init command", zap.Error(err))
			return err
		}
	} else {
		logger.FInfo(cmd.Io.Out, msg.InitDeployCommand)
		logger.FInfo(cmd.Io.Out, fmt.Sprintf(msg.EdgeApplicationsInitSuccessful, info.Name))
	}

	return nil
}
