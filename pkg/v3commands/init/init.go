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
	"github.com/aziontech/azion-cli/pkg/command"
	"github.com/aziontech/azion-cli/pkg/config"
	"github.com/aziontech/azion-cli/pkg/github"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/node"
	"github.com/aziontech/azion-cli/pkg/output"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/aziontech/azion-cli/utils"
	thoth "github.com/aziontech/go-thoth"
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
	packageManager        string
	pathWorkingDir        string
	globalFlagAll         bool
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
	return cmd
}

func (cmd *initCmd) Run(c *cobra.Command, _ []string) error {
	// Set globalFlagAll from factory
	cmd.globalFlagAll = cmd.f.GlobalFlagAll
	
	pathWorkingDirHere, err := cmd.getWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}

	msgs := []string{}

	// Phase 1: Welcome message
	if !cmd.auto {
		cmd.showWelcome()
	}

	// Phase 2: Template selection (moved before name prompt)
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

	var templates []Template
	err = cmd.unmarshal(body, &templates)
	if err != nil {
		return err
	}

	// Flatten all templates into a single list
	var allTemplates []Item
	for _, template := range templates {
		allTemplates = append(allTemplates, template.Items...)
	}

	// Build template options with descriptions
	templateOptions := make([]string, len(allTemplates))
	templateOptionsMap := make(map[string]Item)
	for number, value := range allTemplates {
		// Format template option with description if available
		optionText := value.Name
		if value.Message != "" {
			optionText = fmt.Sprintf("%s - %s", value.Name, value.Message)
		}
		templateOptions[number] = optionText
		templateOptionsMap[optionText] = value
	}

	promptTemplate := &survey.Select{
		Message:  "Choose a template:",
		Options:  templateOptions,
		PageSize: 15, // Show 15 items at a time for better UX
	}

	var answerTemplate string
	err = cmd.askOne(promptTemplate, &answerTemplate)
	if err != nil {
		return err
	}

	// Store selected template info
	selectedTemplate := templateOptionsMap[answerTemplate]
	cmd.preset = strings.ToLower(selectedTemplate.Preset)

	// Phase 3: Project configuration (name comes after template selection)
	// Checks for global --yes flag and that name flag was not sent
	if cmd.globalFlagAll && cmd.name == "" {
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

	// Phase 4: Show configuration summary and confirm
	if !cmd.auto {
		if !cmd.showSummaryAndConfirm(answerTemplate, selectedTemplate) {
			fmt.Fprintln(cmd.f.IOStreams.Out, "\nOperation cancelled.")
			return nil
		}
	}

	// Phase 5: Execute project creation with progress indicators
	if !cmd.auto {
		fmt.Fprintln(cmd.f.IOStreams.Out, "\nCreating your project...")
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

	options := &git.CloneOptions{
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName("v3"),
	}
	err = cmd.git.Clone(options, SAMPLESURL, tempDir)
	if err != nil {
		logger.Debug("Error while cloning the repository", zap.Error(err))
		return err
	}

	oldPath := path.Join(tempDir, "templates", templateOptionsMap[answerTemplate].Path)
	newPath := path.Join(pathWorkingDirHere, cmd.name)

	//move contents from temporary directory into final destination
	err = cmd.rename(oldPath, newPath)
	if err != nil {
		logger.Debug("Error move contents directory", zap.Error(err))
		return utils.ErrorMovingFiles
	}

	if !cmd.auto {
		fmt.Fprintln(cmd.f.IOStreams.Out, "  ✓ Template downloaded")
		fmt.Fprintln(cmd.f.IOStreams.Out, "  ✓ Files extracted")
	}

	if err = cmd.createTemplateAzion(); err != nil {
		return err
	}

	if !cmd.auto {
		fmt.Fprintln(cmd.f.IOStreams.Out, "  ✓ Configuration generated")
	}

	logger.FInfoFlags(cmd.f.IOStreams.Out, msg.WebAppInitCmdSuccess, cmd.f.Format, cmd.f.Out)
	msgs = append(msgs, msg.WebAppInitCmdSuccess)

	err = cmd.changeDir(cmd.pathWorkingDir)
	if err != nil {
		logger.Debug("Error while changing to new working directory", zap.Error(err))
		return msg.ErrorWorkingDir
	}

	vul := vulcanPkg.NewVulcanV3()
	err = cmd.selectVulcanTemplates(vul)
	if err != nil {
		return msg.ErrorGetProjectInfo
	}

	// Phase 6: Show next steps
	cmd.showNextSteps()

	// Phase 7: Optional dev server
	if cmd.auto || !utils.Confirm(cmd.globalFlagAll, msg.AskLocalDev, false) {
		msgs = append(msgs, msg.InitDevCommand)
		msgs = append(msgs, msg.ChangeWorkingDir)
	} else {
		if err := cmd.deps(c, msg.AskInstallDepsDev, &msgs); err != nil {
			return err
		}
		logger.Debug("Running dev command from init command")
		dev := cmd.devCmd(cmd.f)
		err = dev.Run(cmd.f)
		if err != nil {
			logger.Debug("Error while running dev command called by init command", zap.Error(err))
			return err
		}
		return nil
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
		if !utils.Confirm(cmd.globalFlagAll, m, true) {
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
