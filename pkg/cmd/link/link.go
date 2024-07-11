package link

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/link"
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

type LinkInfo struct {
	Name           string
	Preset         string
	Mode           string
	packageManager string
	PathWorkingDir string
	GlobalFlagAll  bool
	remote         string
	Auto           bool
	projectPath    string
}

type LinkCmd struct {
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
	CommandRunner         func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
	CommandRunInteractive func(f *cmdutil.Factory, comm string) error
	ShouldConfigure       func(info *LinkInfo) bool
	ShouldDevDeploy       func(info *LinkInfo, msg string, defaultYes bool) bool
	DeployCmd             func(f *cmdutil.Factory) *deploy.DeployCmd
	DevCmd                func(f *cmdutil.Factory) *dev.DevCmd
	F                     *cmdutil.Factory
}

func NewLinkCmd(f *cmdutil.Factory) *LinkCmd {
	return &LinkCmd{
		Io:              f.IOStreams,
		F:               f,
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
		CommandRunner: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
			return utils.CommandRunInteractiveWithOutput(f, comm, envVars)
		},
		CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
			return utils.CommandRunInteractive(f, comm)
		},
	}
}

func NewCobraCmd(link *LinkCmd, f *cmdutil.Factory) *cobra.Command {
	info := &LinkInfo{}
	cobraCmd := &cobra.Command{
		Use:           msg.EdgeApplicationsLinkUsage,
		Short:         msg.EdgeApplicationsLinkShortDescription,
		Long:          msg.EdgeApplicationsLinkLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azion link
		$ azion link --help
		$ azion link --name "thisisatest" --preset hexo --mode deliver
		$ azion link --preset astro --mode deliver
		$ azion link --name "thisisatest" --preset nextjs
		$ azion link --name "thisisatest" --preset static
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			info.GlobalFlagAll = f.GlobalFlagAll
			return link.run(cmd, info)
		},
	}

	cobraCmd.Flags().StringVar(&info.Name, "name", "", msg.EdgeApplicationsLinkFlagName)
	cobraCmd.Flags().StringVar(&info.Preset, "preset", "", msg.EdgeApplicationsLinkFlagTemplate)
	cobraCmd.Flags().StringVar(&info.packageManager, "package-manager", "", msg.FLAG_PACKAGE_MANAGE)
	cobraCmd.Flags().StringVar(&info.Mode, "mode", "", msg.EdgeApplicationsLinkFlagMode)
	cobraCmd.Flags().BoolVar(&info.Auto, "auto", false, msg.LinkFlagAuto)
	cobraCmd.Flags().StringVar(&info.remote, "remote", "", msg.FLAG_REMOTE)
	cobraCmd.Flags().StringVar(&info.projectPath, "config-dir", "azion", msg.FLAGPATHCONF)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewLinkCmd(f), f)
}

func (cmd *LinkCmd) run(c *cobra.Command, info *LinkInfo) error {
	logger.Debug("Running link command")

	msgs := []string{}
	nodeManager := node.NewNode()
	err := nodeManager.NodeVer(nodeManager)
	if err != nil {
		return err
	}

	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}
	info.PathWorkingDir = path

	git := github.NewGithub()

	if len(info.remote) > 0 {
		logger.Debug("flag remote", zap.Any("repository", info.remote))
		urlFull, _ := regexp.MatchString(`^https?://(?:www\.)?(?:github\.com|gitlab\.com)/[\w.-]+/[\w.-]+(\.git)?$`, info.remote)
		if !urlFull {
			info.remote = fmt.Sprintf("https://github.com/%s.git", info.remote)
		}
		nameRepo := git.GetNameRepo(info.remote)
		info.PathWorkingDir = filepath.Join(info.PathWorkingDir, nameRepo)
		err = git.Clone(info.remote, filepath.Join(path, nameRepo))
		if err != nil {
			logger.Debug("Error while cloning the repository", zap.Error(err))
			return err
		}
	}

	shouldLink := cmd.ShouldConfigure(info)
	if !shouldLink {
		return nil
	}

	shouldFetchTemplates, err := shouldFetch(cmd, info)
	if err != nil {
		return err
	}

	if shouldFetchTemplates {
		// Checks for global --yes flag and that name flag was not sent
		if (info.GlobalFlagAll || info.Auto) && info.Name == "" {
			info.Name = thoth.GenerateName()
		} else {
			// if name was not sent we ask for input, otherwise info.Name already has the value
			if info.Name == "" {
				projName, err := askForInput(msg.LinkProjectQuestion, thoth.GenerateName())
				if err != nil {
					return err
				}
				info.Name = projName
			}
		}

		if info.Preset == "" || info.Mode == "" {
			err = cmd.selectVulcanMode(info)
			if err != nil {
				return err
			}
		}

		if err = cmd.createTemplateAzion(info); err != nil {
			return err
		}

		logger.FInfoFlags(cmd.Io.Out, msg.WebAppLinkCmdSuccess, cmd.F.Format, cmd.F.Out)
		msgs = append(msgs, msg.WebAppLinkCmdSuccess)

		//asks if user wants to add files to .gitignore
		gitignore, err := git.CheckGitignore(info.PathWorkingDir)
		if err != nil {
			return msg.ErrorReadingGitignore
		}

		if !gitignore && (info.Auto || info.GlobalFlagAll || utils.Confirm(info.GlobalFlagAll, msg.AskGitignore, true)) {
			if err := git.WriteGitignore(info.PathWorkingDir); err != nil {
				return msg.ErrorWritingGitignore
			}
			logger.FInfoFlags(cmd.Io.Out, msg.WrittenGitignore, cmd.F.Format, cmd.F.Out)
			msgs = append(msgs, msg.WrittenGitignore)
		}

		if !info.Auto {
			if cmd.ShouldDevDeploy(info, msg.AskLocalDev, false) {
				if err := deps(c, cmd, info, msg.AskInstallDepsDev); err != nil {
					return err
				}

				logger.Debug("Running dev command from link command")
				dev := cmd.DevCmd(cmd.F)
				err = dev.Run(cmd.F)
				if err != nil {
					logger.Debug("Error while running deploy command called by link command", zap.Error(err))
					return err
				}
			} else {
				logger.FInfoFlags(cmd.Io.Out, msg.LinkDevCommand, cmd.F.Format, cmd.F.Out)
				msgs = append(msgs, msg.LinkDevCommand)
			}

			if cmd.ShouldDevDeploy(info, msg.AskDeploy, false) {
				if err := deps(c, cmd, info, msg.AskInstallDepsDeploy); err != nil {
					return err
				}

				logger.Debug("Running deploy command from link command")
				deploy := cmd.DeployCmd(cmd.F)
				err = deploy.ExternalRun(cmd.F, info.projectPath)
				if err != nil {
					logger.Debug("Error while running deploy command called by link command", zap.Error(err))
					return err
				}
			} else {
				logger.FInfoFlags(cmd.Io.Out, msg.LinkDeployCommand, cmd.F.Format, cmd.F.Out)
				msgs = append(msgs, msg.LinkDeployCommand)
				logger.FInfoFlags(cmd.Io.Out, fmt.Sprintf(msg.EdgeApplicationsLinkSuccessful, info.Name), cmd.F.Format, cmd.F.Out)
				msgs = append(msgs, fmt.Sprintf(msg.EdgeApplicationsLinkSuccessful, info.Name))
			}
		}

	}

	initOut := output.SliceOutput{
		GeneralOutput: output.GeneralOutput{
			Out:   cmd.F.IOStreams.Out,
			Flags: cmd.F.Flags,
		},
		Messages: msgs,
	}

	return output.Print(&initOut)
}

func deps(c *cobra.Command, cmd *LinkCmd, info *LinkInfo, m string) error {
	pacManIsInformed := c.Flags().Changed("package-manager")
	if pacManIsInformed || cmd.ShouldDevDeploy(info, m, false) {
		var err error
		pacMan := info.packageManager
		if !pacManIsInformed {
			pacMan, err = utils.GetPackageManager()
			if err != nil {
				return err
			}
		}
		err = depsInstall(cmd, pacMan)
		if err != nil {
			logger.Debug("Error while installing project dependencies", zap.Error(err))
			return msg.ErrorDeps
		}
	}
	return nil
}
