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
	"github.com/aziontech/azion-cli/pkg/contracts"
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
	CommandRunInteractive func(f *cmdutil.Factory, envVars []string, comm string) error
	ShouldConfigure       func(info *InitInfo) (bool, error)
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
		ShouldConfigure: shouldConfigure,
		ShouldDevDeploy: shouldDevDeploy,
		DevCmd:          dev.NewDevCmd,
		DeployCmd:       deploy.NewDeployCmd,
		ChangeDir:       os.Chdir,
		CommandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
		CommandRunInteractive: func(f *cmdutil.Factory, envVars []string, comm string) error {
			return utils.CommandRunInteractive(f, envVars, comm)
		},
	}
}

func NewCobraCmd(init *InitCmd, f *cmdutil.Factory) *cobra.Command {
	options := &contracts.AzionApplicationOptions{}
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
		$ azion init --name "thisisatest" --template nextjs
		$ azion init --name "thisisatest" --template static
		$ azion init --name "thisisatest" --template hexo --mode deliver
		$ azion init --name "thisisatest" --template hexo --mode deliver --auto
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			info.GlobalFlagAll = f.GlobalFlagAll
			return init.run(info, options, cmd)
		},
	}

	cobraCmd.Flags().StringVar(&info.Name, "name", "", msg.EdgeApplicationsInitFlagName)
	cobraCmd.Flags().StringVar(&info.Template, "template", "", msg.EdgeApplicationsInitFlagTemplate)
	cobraCmd.Flags().StringVar(&info.Mode, "mode", "", msg.EdgeApplicationsInitFlagMode)
	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewInitCmd(f), f)
}

func (cmd *InitCmd) run(info *InitInfo, options *contracts.AzionApplicationOptions, c *cobra.Command) error {
	logger.Debug("Running init command")

	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}
	info.PathWorkingDir = path

	switch info.Template {
	case "simple":
		return initSimple(cmd, path, info, c)
	case "static":
		return initStatic(cmd, info, options, c)
	}

	// Checks for global --yes flag and that name flag was not sent
	if info.GlobalFlagAll && !c.Flags().Changed("name") {
		info.Name = thoth.GenerateName()
	} else {
		// if name was not sent we ask for input, otherwise info.Name already has the value
		if !c.Flags().Changed("name") {
			projName, err := askForInput(msg.InitProjectQuestion, thoth.GenerateName())
			if err != nil {
				return err
			}

			info.Name = projName
		}
	}

	if !c.Flags().Changed("template") || !c.Flags().Changed("mode") {
		err = cmd.selectVulcanTemplates(info)
		if err != nil {
			return err
		}
	}

	info.PathWorkingDir = info.PathWorkingDir + "/" + info.Name

	if err = cmd.createTemplateAzion(info); err != nil {
		return err
	}

	logger.FInfo(cmd.Io.Out, msg.WebAppInitCmdSuccess)

	err = cmd.ChangeDir(info.PathWorkingDir)
	if err != nil {
		logger.Debug("Error while changing to new working directory", zap.Error(err))
		return msg.ErrorDeps
	}

	shouldDev, err := cmd.ShouldDevDeploy(info, "Do you want to start a local development server?")
	if err != err {
		return err
	}
	if shouldDev {
		logger.Debug("Running dev command from init command")

		err = yarnInstall(cmd)
		if err != nil {
			logger.Debug("Failed to install project dependencies")
			return err
		}

		// Run build command
		dev := cmd.DevCmd(cmd.F)
		err := dev.Run(cmd.F)
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
		logger.Debug("Running deploy command from init command")

		err = yarnInstall(cmd)
		if err != nil {
			logger.Debug("Failed to install project dependencies")
			return err
		}

		// Run build command
		deploy := cmd.DeployCmd(cmd.F)
		err := deploy.Run(cmd.F)
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
