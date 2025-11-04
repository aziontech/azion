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
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/cmd/deploy"
	"github.com/aziontech/azion-cli/pkg/cmd/dev"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/command"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/github"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/node"
	"github.com/aziontech/azion-cli/pkg/output"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	SAMPLESURL = "https://github.com/aziontech/azion-samples.git"
	APIURL     = "https://api.azion.com/v4/utils/project_samples"
)

type initCmd struct {
	name                  string
	preset                string
	auto                  bool
	sync                  bool
	local                 bool
	SkipFramework         bool
	packageManager        string
	pathWorkingDir        string
	f                     *cmdutil.Factory
	git                   github.Github
	getWorkDir            func() (string, error)
	fileReader            func(path string) ([]byte, error)
	isDirEmpty            func(dirpath string) (bool, error)
	cleanDir              func(dirpath string) error
	writeFile             func(filename string, data []byte, perm fs.FileMode) error
	openFile              func(name string) (*os.File, error)
	removeAll             func(path string) error
	rename                func(oldpath string, newpath string) error
	envLoader             func(path string) ([]string, error)
	stat                  func(path string) (fs.FileInfo, error)
	mkdir                 func(path string, perm os.FileMode) error
	gitPlainClone         func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	commandRunner         func(envVars []string, comm string) (string, int, error)
	commandRunnerOutput   func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
	commandRunInteractive func(f *cmdutil.Factory, comm string) error
	deployCmd             func(f *cmdutil.Factory) *deploy.DeployCmd
	devCmd                func(f *cmdutil.Factory) *dev.DevCmd
	changeDir             func(dir string) error
	askOne                func(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
	load                  func(filenames ...string) (err error)
	dir                   func() config.DirPath
	mkdirTemp             func(dir, pattern string) (string, error)
	readAll               func(r io.Reader) ([]byte, error)
	get                   func(url string) (resp *http.Response, err error)
	marshalIndent         func(v any, prefix, indent string) ([]byte, error)
	unmarshal             func(data []byte, v any) error
	DetectPackageManager  func(pathWorkDir string) string
}

func NewInitCmd(f *cmdutil.Factory) *initCmd {
	return &initCmd{
		f:                     f,
		getWorkDir:            utils.GetWorkingDir,
		fileReader:            os.ReadFile,
		isDirEmpty:            utils.IsDirEmpty,
		cleanDir:              utils.CleanDirectory,
		writeFile:             os.WriteFile,
		openFile:              os.Open,
		removeAll:             os.RemoveAll,
		rename:                os.Rename,
		mkdirTemp:             os.MkdirTemp,
		envLoader:             utils.LoadEnvVarsFromFile,
		stat:                  os.Stat,
		mkdir:                 os.MkdirAll,
		gitPlainClone:         git.PlainClone,
		devCmd:                dev.NewDevCmd,
		deployCmd:             deploy.NewDeployCmd,
		changeDir:             os.Chdir,
		commandRunner:         command.RunCommandWithOutput,
		commandRunInteractive: command.CommandRunInteractive,
		commandRunnerOutput:   command.CommandRunInteractiveWithOutput,
		askOne:                survey.AskOne,
		load:                  godotenv.Load,
		dir:                   config.Dir,
		readAll:               io.ReadAll,
		get:                   http.Get,
		marshalIndent:         json.MarshalIndent,
		unmarshal:             json.Unmarshal,
		git:                   *github.NewGithub(),
		DetectPackageManager:  node.DetectPackageManager,
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
	cmd.Flags().BoolVar(&init.sync, "sync", false, msg.FLAG_SYNC)
	cmd.Flags().BoolVar(&init.local, "local", false, msg.FLAG_LOCAL)
	cmd.Flags().BoolVar(&init.SkipFramework, "skip-framework-build", false, msg.SkipFrameworkBuild)
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
	if cmd.f.GlobalFlagAll && cmd.name == "" {
		cmd.name = thoth.GenerateName()
	} else {
		// if name was not sent we ask for input, otherwise info.Name already has the value
		if cmd.name == "" {
			projName, err := cmd.askForInput(msg.InitProjectQuestion, thoth.GenerateName())
			if err != nil {
				return err
			}
			cmd.name = projName
		}
	}

	cmd.pathWorkingDir = path.Join(pathWorkingDirHere, cmd.name)
	resp, err := cmd.get(APIURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code %d but got %d", http.StatusOK, resp.StatusCode)
	}

	body, err := cmd.readAll(resp.Body)
	if err != nil {
		return err
	}

	templateMap := make(map[string][]Item)

	var templates []Template
	err = cmd.unmarshal(body, &templates)
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
	err = cmd.askOne(prompt, &answer)
	if err != nil {
		return err
	}

	templateOptions := []string{}
	templateOptionsMap := make(map[string]Item)
	for _, value := range templateMap[answer] {
		if value.RequiresAdditionalBuild != nil && *value.RequiresAdditionalBuild {
			continue
		}
		templateOptions = append(templateOptions, value.Name)
		templateOptionsMap[value.Name] = value
	}

	promptTemplate := &survey.Select{
		Message:  "Choose a template:",
		Options:  templateOptions,
		PageSize: len(templateOptions),
	}

	var answerTemplate string
	err = cmd.askOne(promptTemplate, &answerTemplate)
	if err != nil {
		return err
	}

	dirPath := cmd.dir()

	// Create a temporary directory
	tempDir, err := cmd.mkdirTemp(dirPath.Dir, "tempclonesamples")
	if err != nil {
		return err
	}

	// Defer deletion of the temporary directory
	defer func() {
		err := cmd.removeAll(tempDir)
		if err != nil {
			log.Fatal(err)
		}
	}()

	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Suffix = " Fetching selected template..."
	s.FinalMSG = "Template successfully fetched\n"
	if !cmd.f.Debug {
		s.Start() // Start the spinner
	}

	// options := &git.CloneOptions{}
	options := &git.CloneOptions{
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName("dev"),
	}
	err = cmd.git.Clone(options, SAMPLESURL, tempDir)
	if err != nil {
		logger.Debug("Error while cloning the repository", zap.Error(err))
		return err
	}

	s.Stop()

	oldPath := path.Join(tempDir, "templates", templateOptionsMap[answerTemplate].Path)
	newPath := path.Join(pathWorkingDirHere, cmd.name)

	//move contents from temporary directory into final destination
	err = cmd.rename(oldPath, newPath)
	if err != nil {
		logger.Debug("Error move contents directory", zap.Error(err))
		return utils.ErrorMovingFiles
	}

	cmd.preset = strings.ToLower(templateOptionsMap[answerTemplate].Preset)
	if err = cmd.createTemplateAzion(); err != nil {
		return err
	}

	// Handle optional extras from the Templates API
	selectedItem := templateOptionsMap[answerTemplate]
	if selectedItem.Extras != nil && len(selectedItem.Extras.Inputs) > 0 {
		inputs := make([]utils.EnvInput, 0, len(selectedItem.Extras.Inputs))
		for _, in := range selectedItem.Extras.Inputs {
			inputs = append(inputs, utils.EnvInput{Key: in.Key, Text: in.Text, IsSecret: in.IsSecret})
		}
		switch strings.ToLower(selectedItem.Extras.Type) {
		case "env":
			if err := utils.CollectEnvInputsAndWriteFile(inputs, newPath); err != nil {
				return err
			}
		case "args":
			if err := utils.CollectArgsInputsAndWriteFile(inputs, path.Join(newPath, "azion")); err != nil {
				return err
			}
		default:
			return fmt.Errorf("invalid extra type: %s", selectedItem.Extras.Type)
		}
	}

	logger.FInfoFlags(cmd.f.IOStreams.Out, msg.WebAppInitCmdSuccess, cmd.f.Format, cmd.f.Out)
	msgs = append(msgs, msg.WebAppInitCmdSuccess)

	err = cmd.changeDir(cmd.pathWorkingDir)
	if err != nil {
		logger.Debug("Error while changing to new working directory", zap.Error(err))
		return msg.ErrorWorkingDir
	}

	vul := vulcanPkg.NewVulcan()
	if err := cmd.deps(c, msg.AskInstallDepsBuild, &msgs); err != nil {
		return err
	}
	err = cmd.selectVulcanTemplates(vul)
	if err != nil {
		return msg.ErrorGetProjectInfo
	}

	if cmd.auto || !utils.Confirm(cmd.f.GlobalFlagAll, msg.AskLocalDev, false) {
		logger.FInfoFlags(cmd.f.IOStreams.Out, msg.InitDevCommand, cmd.f.Format, cmd.f.Out)
		logger.FInfoFlags(cmd.f.IOStreams.Out, msg.ChangeWorkingDir, cmd.f.Format, cmd.f.Out)
		msgs = append(msgs, msg.InitDevCommand)
		msgs = append(msgs, msg.ChangeWorkingDir)
	} else {
		if err := cmd.deps(c, msg.AskInstallDepsDev, &msgs); err != nil {
			return err
		}
		logger.Debug("Running dev command from init command")
		dev := cmd.devCmd(cmd.f)
		err = dev.ExternalRun(cmd.f, cmd.SkipFramework)
		if err != nil {
			logger.Debug("Error while running dev command called by init command", zap.Error(err))
			return err
		}
	}

	if cmd.auto || !utils.Confirm(cmd.f.GlobalFlagAll, msg.AskDeploy, false) {
		logger.FInfoFlags(cmd.f.IOStreams.Out, msg.InitDeployCommand, cmd.f.Format, cmd.f.Out)
		logger.FInfoFlags(cmd.f.IOStreams.Out, msg.ChangeWorkingDir, cmd.f.Format, cmd.f.Out)
		msgs = append(msgs, msg.InitDeployCommand)
		msgs = append(msgs, msg.ChangeWorkingDir)
		msgEdgeAppInitSuccessFull := fmt.Sprintf(msg.EdgeApplicationsInitSuccessful, cmd.name)
		logger.FInfoFlags(cmd.f.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsInitSuccessful, cmd.name),
			cmd.f.Format, cmd.f.Out)
		msgs = append(msgs, msgEdgeAppInitSuccessFull)
	} else {
		if err := cmd.deps(c, msg.AskInstallDepsDeploy, &msgs); err != nil {
			return err
		}
		logger.Debug("Running deploy command from init command")
		deploy := cmd.deployCmd(cmd.f)
		err = deploy.ExternalRun(cmd.f, "azion", cmd.sync, cmd.local, cmd.SkipFramework)
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

func (cmd *initCmd) deps(c *cobra.Command, m string, msgs *[]string) error {
	if !c.Flags().Changed("package-manager") {
		if !utils.Confirm(cmd.f.GlobalFlagAll, m, true) {
			return nil
		}

		pathWorkDir, err := cmd.getWorkDir()
		if err != nil {
			return err
		}

		cmd.packageManager = node.DetectPackageManager(pathWorkDir)
	}

	logger.FInfoFlags(cmd.f.IOStreams.Out, msg.InstallDeps, cmd.f.Format, cmd.f.Out)
	*msgs = append(*msgs, msg.InstallDeps)

	if err := cmd.depsInstall(); err != nil {
		logger.Debug("Failed to install project dependencies")
		return err
	}

	return nil
}
