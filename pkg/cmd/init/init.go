package init

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/cmd/deploy"
	"github.com/aziontech/azion-cli/pkg/cmd/dev"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/github"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	SAMPLESURL = "https://github.com/aziontech/azion-samples.git"
	APIURL     = "https://os4bzngzt0.map.azionedge.net/api/templates"
)

type initCmd struct {
	name                  string
	preset                string
	auto                  bool
	mode                  string
	packageManager        string
	pathWorkingDir        string
	globalFlagAll         bool
	f                     *cmdutil.Factory
	io                    *iostreams.IOStreams
	getWorkDir            func() (string, error)
	fileReader            func(path string) ([]byte, error)
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
	cmd.Flags().BoolVar(&init.auto, "auto", false, msg.FLAG_AUTO)
	return cmd
}

func (cmd *initCmd) Run(c *cobra.Command, _ []string) error {
	pathWorkingDirHere, err := cmd.getWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}

	msgs := []string{}

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

	cmd.pathWorkingDir = path.Join(pathWorkingDirHere, cmd.name)

	resp, err := http.Get(APIURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	templateMap := make(map[string][]Item)

	var templates []Template
	err = json.Unmarshal(body, &templates)
	if err != nil {
		return err
	}

	listTemplates := make([]string, len(templates))

	for number, value := range templates {
		templateMap[value.Name] = value.Items
		listTemplates[number] = value.Name
	}

	prompt := &survey.Select{
		Message:  "Choose a preset:",
		Options:  listTemplates,
		PageSize: len(listTemplates),
	}

	var answer string
	err = survey.AskOne(prompt, &answer)
	if err != nil {
		return err
	}

	templateOptions := make([]string, len(templateMap[answer]))
	templateOptionsMap := make(map[string]Item)
	for number, value := range templateMap[answer] {
		templateOptions[number] = value.Name
		templateOptionsMap[value.Name] = value
	}

	promptTemplate := &survey.Select{
		Message:  "Choose a template:",
		Options:  templateOptions,
		PageSize: len(templateOptions),
	}

	var answerTemplate string
	err = survey.AskOne(promptTemplate, &answerTemplate)
	if err != nil {
		return err
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "tempclonesamples")
	if err != nil {
		return err
	}

	// Defer deletion of the temporary directory
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Fatal(err)
		}
	}()

	git := github.NewGithub()
	err = git.Clone(SAMPLESURL, tempDir)
	if err != nil {
		logger.Debug("Error while cloning the repository", zap.Error(err))
		return err
	}

	//move contents from temporary directory into final destination
	err = cmd.rename(path.Join(tempDir, "templates", templateOptionsMap[answerTemplate].Path), path.Join(pathWorkingDirHere, cmd.name))
	if err != nil {
		fmt.Println(err.Error())
		return utils.ErrorMovingFiles
	}

	cmd.preset = strings.ToLower(templateOptionsMap[answerTemplate].Preset)
	cmd.mode = strings.ToLower(templateOptionsMap[answerTemplate].Mode)

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

	err = cmd.selectVulcanTemplates()
	if err != nil {
		return err
	}

	if cmd.auto || !cmd.shouldDevDeploy(msg.AskLocalDev, cmd.globalFlagAll, false) {
		logger.FInfoFlags(cmd.io.Out, msg.InitDevCommand, cmd.f.Format, cmd.f.Out)
		msgs = append(msgs, msg.InitDevCommand)
	} else {
		if err := deps(c, cmd, msg.AskInstallDepsDev); err != nil {
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
		if err := deps(c, cmd, msg.AskInstallDepsDeploy); err != nil {
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

func deps(c *cobra.Command, cmd *initCmd, m string) error {
	pacManIsInformed := c.Flags().Changed("package-manager")
	var err error
	if pacManIsInformed || cmd.shouldDevDeploy(m, cmd.globalFlagAll, false) {
		pacMan := cmd.packageManager
		if !pacManIsInformed {
			pacMan, err = utils.GetPackageManager()
			if err != nil {
				return err
			}
		}
		err = depsInstall(cmd, pacMan)
		if err != nil {
			logger.Debug("Failed to install project dependencies")
			return err
		}
	}
	return nil
}
