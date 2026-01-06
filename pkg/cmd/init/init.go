package init

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/charmbracelet/huh"
	huhspinner "github.com/charmbracelet/huh/spinner"
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
	runForm               func(form *huh.Form) error
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
		runForm: func(form *huh.Form) error {
			return form.Run()
		},
		load: godotenv.Load,
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

	// Prepare data structures for dynamic form
	var selectedCategory string
	var selectedFramework string
	var answerTemplate string
	var projectName string
	var confirmProceed bool = true
	templateOptionsMap := make(map[string]Item)

	// Get unique frameworks from templates
	frameworkSet := make(map[string]bool)
	for _, template := range templates {
		for _, item := range template.Items {
			if item.RequiresAdditionalBuild != nil && *item.RequiresAdditionalBuild {
				continue
			}
			if cmd.isFrameworkTemplate(item.Preset, item.Name) {
				fw := cmd.normalizeFrameworkName(item.Preset)
				frameworkSet[fw] = true
			}
		}
	}

	// Build framework options list
	frameworks := make([]string, 0, len(frameworkSet))
	for fw := range frameworkSet {
		frameworks = append(frameworks, fw)
	}
	sort.Strings(frameworks)

	// Create a single form with all steps to preserve answers on terminal
	categories := []string{
		"Simple Hello World",
		"JavaScript",
		"TypeScript",
		"Frameworks",
	}

	// Styles for displaying answers
	labelStyle := GetAzionLabelStyle()   // Purple (#b5b1f4)
	answerStyle := GetAzionAnswerStyle() // Orange (#f3652b)
	
	// Step 1: Category selection
	err = huh.NewSelect[string]().
		Title("Choose a category:").
		Options(huh.NewOptions(categories...)...).
		Value(&selectedCategory).
		WithTheme(ThemeAzion()).
		Run()
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, "%s %s\n", labelStyle.Render("Category:"), answerStyle.Render(selectedCategory))
	
	// Step 2: Framework selection (only if Frameworks category)
	if selectedCategory == "Frameworks" {
		err = huh.NewSelect[string]().
			Title("Choose a framework:").
			Options(huh.NewOptions(frameworks...)...).
			Height(10).
			Value(&selectedFramework).
			WithTheme(ThemeAzion()).
			Run()
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.f.IOStreams.Out, "%s %s\n", labelStyle.Render("Framework:"), answerStyle.Render(selectedFramework))
	}
	
	// Step 3: Template selection (filtered based on category/framework)
	var filteredTemplates []Item
	for _, template := range templates {
		for _, item := range template.Items {
			if item.RequiresAdditionalBuild != nil && *item.RequiresAdditionalBuild {
				continue
			}

			if cmd.matchesCategory(selectedCategory, item, template.Name) {
				if selectedCategory == "Frameworks" {
					if cmd.matchesFramework(selectedFramework, item) {
						filteredTemplates = append(filteredTemplates, item)
					}
				} else {
					filteredTemplates = append(filteredTemplates, item)
				}
			}
		}
	}

	// Build template options
	templateOptions := make([]string, len(filteredTemplates))
	for number, value := range filteredTemplates {
		templateName := value.Name
		
		// Add (Edge-Runtime) tag for NextJs templates (not OpenNext)
		if strings.Contains(strings.ToLower(value.Preset), "next") && 
		   !strings.Contains(strings.ToLower(value.Preset), "open") {
			templateName = fmt.Sprintf("%s (Edge-Runtime)", value.Name)
		}
		
		templateOptions[number] = templateName
		templateOptionsMap[templateName] = value
	}

	err = huh.NewSelect[string]().
		Title("Choose a template:").
		Options(huh.NewOptions(templateOptions...)...).
		Height(10).
		Value(&answerTemplate).
		WithTheme(ThemeAzion()).
		Run()
	if err != nil {
		return err
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, "%s %s\n", labelStyle.Render("Template:"), answerStyle.Render(answerTemplate))
	
	// Step 4: Project name input (only if not provided via flag)
	if cmd.name == "" && !cmd.f.GlobalFlagAll {
		err = huh.NewInput().
			Title(msg.InitProjectQuestion).
			Value(&projectName).
			Placeholder(thoth.GenerateName()).
			Description("This will be your application identifier on Azion").
			WithTheme(ThemeAzion()).
			Run()
		if err != nil {
			return err
		}
		displayName := projectName
		if displayName == "" {
			displayName = thoth.GenerateName()
		}
		fmt.Fprintf(cmd.f.IOStreams.Out, "%s %s\n\n", labelStyle.Render("Project Name:"), answerStyle.Render(displayName))
	}
	
	// Step 5: Confirmation (only in interactive mode)
	if !cmd.auto {
		name := projectName
		if name == "" {
			name = thoth.GenerateName()
		}
		if cmd.name != "" {
			name = cmd.name
		}
		
		err = huh.NewConfirm().
			Title(fmt.Sprintf("Create project '%s' with template '%s'?", name, answerTemplate)).
			Value(&confirmProceed).
			Affirmative("Yes").
			Negative("No").
			WithTheme(ThemeAzion()).
			Run()
		if err != nil {
			return err
		}
	}

	// Check if user cancelled
	if !confirmProceed && !cmd.auto {
		fmt.Fprintln(cmd.f.IOStreams.Out, "\nOperation cancelled.")
		return nil
	}

	// Store selected template info
	selectedTemplate := templateOptionsMap[answerTemplate]
	cmd.preset = strings.ToLower(selectedTemplate.Preset)

	// Set project name from form input or use default
	if cmd.name == "" {
		if projectName != "" {
			cmd.name = projectName
		} else {
			cmd.name = thoth.GenerateName()
		}
	}

	cmd.pathWorkingDir = path.Join(pathWorkingDirHere, cmd.name)

	// Phase 5: Execute project creation with progress indicators
	cmd.printTitle("Creating your project...")

	dirPath := cmd.dir()

	// Create a temporary directory
	tempDir, err := cmd.mkdirTemp(dirPath.Dir, "tempclonesamples")
	if err != nil {
		return err
	}

	// Defer deletion of the temporary directory
	defer func() {
		if cleanupErr := cmd.removeAll(tempDir); cleanupErr != nil {
			logger.Debug("Failed to cleanup temporary directory", zap.Error(cleanupErr))
		}
	}()

	// Use huh spinner for fetching template
	var cloneErr error
	if !cmd.f.Debug {
		err = huhspinner.New().
			Title("Fetching selected template...").
			Action(func() {
				options := &git.CloneOptions{
					SingleBranch:  true,
					ReferenceName: plumbing.ReferenceName("dev"),
				}
				cloneErr = cmd.git.Clone(options, SAMPLESURL, tempDir)
			}).
			Run()
		
		if err != nil {
			logger.Debug("Error running spinner", zap.Error(err))
			return err
		}
		
		if cloneErr != nil {
			logger.Debug("Error while cloning the repository", zap.Error(cloneErr))
			return cloneErr
		}
	} else {
		// In debug mode, run without spinner
		options := &git.CloneOptions{
			SingleBranch:  true,
			ReferenceName: plumbing.ReferenceName("dev"),
		}
		err = cmd.git.Clone(options, SAMPLESURL, tempDir)
		if err != nil {
			logger.Debug("Error while cloning the repository", zap.Error(err))
			return err
		}
	}

	oldPath := path.Join(tempDir, "templates", templateOptionsMap[answerTemplate].Path)
	newPath := path.Join(pathWorkingDirHere, cmd.name)

	//move contents from temporary directory into final destination
	err = cmd.rename(oldPath, newPath)
	if err != nil {
		logger.Debug("Error move contents directory", zap.Error(err))
		return utils.ErrorMovingFiles
	}

	cmd.printSuccess("Template downloaded")
	cmd.printSuccess("Files extracted")

	cmd.preset = strings.ToLower(templateOptionsMap[answerTemplate].Preset)
	if err = cmd.createTemplateAzion(); err != nil {
		return err
	}

	cmd.printSuccess("Configuration generated")

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

	if !cmd.auto {
		successStyle := GetAzionSuccessStyle() // Orange (#f3652b)
		fmt.Fprintln(cmd.f.IOStreams.Out, "")
		fmt.Fprintln(cmd.f.IOStreams.Out, successStyle.Render("Template successfully configured"))
	} else {
		// In auto mode, use the logger
		logger.FInfoFlags(cmd.f.IOStreams.Out, msg.WebAppInitCmdSuccess, cmd.f.Format, cmd.f.Out)
	}
	
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

	// Phase 6: Show next steps
	if !cmd.auto {
		cmd.showNextSteps()
	}

	// Phase 7: Optional dev server
	if cmd.auto || !cmd.confirmWithHuh(msg.AskLocalDev, false) {
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

	if cmd.auto || !cmd.confirmWithHuh(msg.AskDeploy, false) {
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
		// Use huh.Confirm for better UX
		if !cmd.confirmWithHuh(m, true) {
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

// confirmWithHuh shows a confirmation prompt using huh library
func (cmd *initCmd) confirmWithHuh(message string, defaultYes bool) bool {
	// Skip confirmation if --yes flag is set
	if cmd.f.GlobalFlagAll {
		return true
	}

	var confirm bool = defaultYes
	
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Value(&confirm).
				Affirmative("Yes").
				Negative("No"),
		),
	).WithTheme(ThemeAzion())

	err := cmd.runForm(confirmForm)
	if err != nil {
		// If there's an error (e.g., user cancelled), return false
		return false
	}

	return confirm
}

// printSuccess prints a success message with an orange checkmark
func (cmd *initCmd) printSuccess(message string) {
	if cmd.auto {
		return
	}
	successStyle := GetAzionSuccessStyle() // Orange (#f3652b)
	checkmark := successStyle.Render("✓")
	fmt.Fprintf(cmd.f.IOStreams.Out, "  %s %s\n", checkmark, message)
}

// printTitle prints a title message in bold with color
func (cmd *initCmd) printTitle(message string) {
	if cmd.auto {
		return
	}
	titleStyle := GetAzionTitleStyle() // Purple (#b5b1f4)
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, titleStyle.Render(message))
}
