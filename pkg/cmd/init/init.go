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
	"github.com/aziontech/azion-cli/pkg/github"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/node"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type initCmd struct {
	name                  string
	preset                string
	template              string
	auto                  bool
	mode                  string
	packageManager        string
	pathWorkingDir        string
	globalFlagAll         bool
	f                     *cmdutil.Factory
	io                    *iostreams.IOStreams
	getWorkDir            func() (string, error)
	fileReader            func(path string) ([]byte, error)
	lookPath              func(bin string) (string, error)
	isDirEmpty            func(dirpath string) (bool, error)
	cleanDir              func(dirpath string) error
	writeFile             func(filename string, data []byte, perm fs.FileMode) error
	openFile              func(name string) (*os.File, error)
	removeAll             func(path string) error
	rename                func(oldpath string, newpath string) error
	createTempDir         func(dir string, pattern string) (string, error)
	envLoader             func(path string) ([]string, error)
	stat                  func(path string) (fs.FileInfo, error)
	mkdir                 func(path string, perm os.FileMode) error
	gitPlainClone         func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	commandRunner         func(cmd string, envvars []string) (string, int, error)
	commandRunnerOutput   func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
	commandRunInteractive func(f *cmdutil.Factory, comm string) error
	shouldDevDeploy       func(msg string, globalFlagAll, defaultYes bool) bool
	deployCmd             func(f *cmdutil.Factory) *deploy.DeployCmd
	devCmd                func(f *cmdutil.Factory) *dev.DevCmd
	changeDir             func(dir string) error
}

func NewInitCmd(f *cmdutil.Factory) *initCmd {
	return &initCmd{
		f:               f,
		io:              f.IOStreams,
		getWorkDir:      utils.GetWorkingDir,
		fileReader:      os.ReadFile,
		lookPath:        exec.LookPath,
		isDirEmpty:      utils.IsDirEmpty,
		cleanDir:        utils.CleanDirectory,
		writeFile:       os.WriteFile,
		openFile:        os.Open,
		removeAll:       os.RemoveAll,
		rename:          os.Rename,
		createTempDir:   os.MkdirTemp,
		envLoader:       utils.LoadEnvVarsFromFile,
		stat:            os.Stat,
		mkdir:           os.MkdirAll,
		gitPlainClone:   git.PlainClone,
		shouldDevDeploy: shouldDevDeploy,
		devCmd:          dev.NewDevCmd,
		deployCmd:       deploy.NewDeployCmd,
		changeDir:       os.Chdir,
		commandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
		commandRunInteractive: func(f *cmdutil.Factory, comm string) error {
			return utils.CommandRunInteractive(f, comm)
		},
		commandRunnerOutput: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
			return utils.CommandRunInteractiveWithOutput(f, comm, envVars)
		},
	}
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	init := NewInitCmd(f)
	cmd := &cobra.Command{
		Use:           msg.USAGE,
		Short:         msg.SHORT_DESCRIPTION,
		Long:          msg.LONG_DESCRIPTION,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example:       heredoc.Doc(msg.EXAMPLE),
		RunE:          init.Run,
	}
	cmd.Flags().StringVar(&init.name, "name", "", msg.FLAG_NAME)
	cmd.Flags().StringVar(&init.packageManager, "package-manager", "", msg.FLAG_PACKAGE_MANAGE)
	cmd.Flags().StringVar(&init.preset, "preset", "", msg.FLAG_PRESET)
	cmd.Flags().StringVar(&init.template, "template", "", msg.FLAG_TEMPLATE)
	cmd.Flags().BoolVar(&init.auto, "auto", false, msg.FLAG_AUTO)
	return cmd
}

func (cmd *initCmd) Run(c *cobra.Command, _ []string) error {
	logger.Debug("Running init command")

	msgs := []string{}
	nodeManager := node.NewNode()
	err := nodeManager.NodeVer(nodeManager)
	if err != nil {
		return err
	}

	cmd.globalFlagAll = cmd.f.GlobalFlagAll

	path, err := cmd.getWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}
	cmd.pathWorkingDir = path

	// Checks for global --yes flag and that name flag was not sent
	if cmd.globalFlagAll && cmd.name == "" {
		cmd.name = thoth.GenerateName()
	} else {
		// if name was not sent we ask for input, otherwise info.Name already has the value
		if cmd.name == "" {
			projName, err := askForInput(msg.InitProjectQuestion, thoth.GenerateName())
			if err != nil {
				return err
			}
			cmd.name = projName
		}
	}

	cmd.pathWorkingDir = cmd.pathWorkingDir + "/" + cmd.name
	err = cmd.selectVulcanTemplates()
	if err != nil {
		return err
	}

	if err = cmd.createTemplateAzion(); err != nil {
		return err
	}
	logger.FInfoFlags(cmd.io.Out, msg.WebAppInitCmdSuccess, cmd.f.Format, cmd.f.Out)
	msgs = append(msgs, msg.WebAppInitCmdSuccess)

	err = cmd.changeDir(cmd.pathWorkingDir)
	if err != nil {
		logger.Debug("Error while changing to new working directory", zap.Error(err))
		return msg.ErrorWorkingDir
	}

	git := github.NewGithub()

	gitignore, err := git.CheckGitignore(cmd.pathWorkingDir)
	if err != nil {
		return msg.ErrorReadingGitignore
	}
	if !gitignore && (cmd.auto || cmd.f.GlobalFlagAll || utils.Confirm(cmd.f.GlobalFlagAll, msg.AskGitignore, true)) {
		if err := git.WriteGitignore(cmd.pathWorkingDir); err != nil {
			return msg.ErrorWritingGitignore
		}
		logger.FInfoFlags(cmd.f.IOStreams.Out, msg.WrittenGitignore, cmd.f.Format, cmd.f.Out)
		msgs = append(msgs, msg.WrittenGitignore)
	}

	if cmd.auto || !cmd.shouldDevDeploy(msg.AskLocalDev, cmd.globalFlagAll, false) {
		logger.FInfoFlags(cmd.io.Out, msg.InitDevCommand, cmd.f.Format, cmd.f.Out)
		msgs = append(msgs, msg.InitDevCommand)
	} else {
		if err := deps(c, cmd, msg.AskInstallDepsDev, &msgs); err != nil {
			return err
		}
		logger.Debug("Running dev command from init command")
		dev := cmd.devCmd(cmd.f)
		err = dev.Run(cmd.f)
		if err != nil {
			logger.Debug("Error while running dev command called by init command", zap.Error(err))
			return err
		}
	}

	if cmd.auto || !cmd.shouldDevDeploy(msg.AskDeploy, cmd.globalFlagAll, false) {
		logger.FInfoFlags(cmd.io.Out, msg.InitDeployCommand, cmd.f.Format, cmd.f.Out)
		msgs = append(msgs, msg.InitDeployCommand)
		msgEdgeAppInitSuccessFul := fmt.Sprintf(msg.EdgeApplicationsInitSuccessful, cmd.name)
		logger.FInfoFlags(cmd.io.Out, fmt.Sprintf(msg.EdgeApplicationsInitSuccessful, cmd.name),
			cmd.f.Format, cmd.f.Out)
		msgs = append(msgs, msgEdgeAppInitSuccessFul)
	} else {
		if err := deps(c, cmd, msg.AskInstallDepsDeploy, &msgs); err != nil {
			return err
		}

		logger.Debug("Running deploy command from init command")
		deploy := cmd.deployCmd(cmd.f)
		err = deploy.Run(cmd.f)
		if err != nil {
			logger.Debug("Error while running deploy command called by init command", zap.Error(err))
			return err
		}
	}

	initOut := output.SliceOutput{
		GeneralOutput: output.GeneralOutput{
			Out:   cmd.f.IOStreams.Out,
			Flags: cmd.f.Flags,
		},
		Messages: msgs,
	}
	return output.Print(&initOut)
}

func deps(c *cobra.Command, cmd *initCmd, m string, msgs *[]string) error {
	if !c.Flags().Changed("package-manager") {
		if !cmd.shouldDevDeploy(m, cmd.globalFlagAll, true) {
			return nil
		}

		pathWorkDir, err := cmd.getWorkDir()
		if err != nil {
			return err
		}

		cmd.packageManager = node.DetectPackageManager(pathWorkDir)
	}

	logger.FInfoFlags(cmd.io.Out, msg.InstallDeps, cmd.f.Format, cmd.f.Out)
	*msgs = append(*msgs, msg.InstallDeps)

	if err := depsInstall(cmd, cmd.packageManager); err != nil {
		logger.Debug("Failed to install project dependencies")
		return err
	}

	return nil
}

