package init

import (
	"fmt"
	"path"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"go.uber.org/zap"
)

var jsonTemplate = `{
	"scope": "global",
	"preset": "%s"
  }`

func (cmd *initCmd) askForInput(msg string, defaultIn string) (string, error) {
	var userInput string
	prompt := &survey.Input{
		Message: msg,
		Default: defaultIn,
		Help:    "This will be your application identifier on Azion",
	}

	// Prompt the user for input
	err := cmd.askOne(prompt, &userInput, survey.WithKeepFilter(true))
	if err != nil {
		return "", err
	}

	return userInput, nil
}

func (cmd *initCmd) selectVulcanTemplates(vul *vulcanPkg.VulcanPkg) error {
	// checking if vulcan major is correct
	vulcanVer, err := cmd.commandRunnerOutput(cmd.f, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	err = vul.CheckVulcanMajor(vulcanVer, cmd.f, vul)
	if err != nil {
		return err
	}

	cmdVulcanInit := "store init"
	if len(cmd.preset) > 0 {
		formatted := fmt.Sprintf(jsonTemplate, cmd.preset)
		cmdVulcanInit = fmt.Sprintf("%s --config '%s'", cmdVulcanInit, formatted)
	}

	// cmdVulcanBuild := "build"
	// if len(cmd.preset) > 0 {
	// 	cmdVulcanBuild = fmt.Sprintf("%s --preset '%s' --only-generate-config --skip-framework-build", cmdVulcanBuild, cmd.preset)
	// }

	command := vul.Command("", cmdVulcanInit, cmd.f)

	err = cmd.commandRunInteractive(cmd.f, command)
	if err != nil {
		return err
	}

	preset, err := cmd.getVulcanInfo()
	if err != nil {
		return err
	}

	if preset == strings.ToLower("vite") {
		preset = "vue"
	}

	cmd.preset = strings.ToLower(preset)
	return nil
}

func (cmd *initCmd) depsInstall() error {
	command := fmt.Sprintf("%s install", cmd.packageManager)
	err := cmd.commandRunInteractive(cmd.f, command)
	if err != nil {
		logger.Debug("Error while running command with simultaneous output", zap.Error(err))
		return msg.ErrorDeps
	}
	return nil
}

func (cmd *initCmd) getVulcanInfo() (string, error) {

	fileContent, err := cmd.fileReader(path.Join(cmd.pathWorkingDir, "info.json"))
	if err != nil {
		logger.Debug("Error reading template info", zap.Error(err))
		return "", err
	}

	var infoJson map[string]string
	err = cmd.unmarshal(fileContent, &infoJson)
	if err != nil {
		logger.Debug("Error unmarshalling template info", zap.Error(err))
		return "", err
	}

	logger.Debug("Information about the template:", zap.Any("preset", infoJson["preset"]))
	return infoJson["preset"], nil
}

// showWelcome displays the welcome message
func (cmd *initCmd) showWelcome() {
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "╭─────────────────────────────────────────────────────────╮")
	fmt.Fprintln(cmd.f.IOStreams.Out, "│  🚀 Welcome to Azion Web Platform                      │")
	fmt.Fprintln(cmd.f.IOStreams.Out, "│  Let's create your web application                     │")
	fmt.Fprintln(cmd.f.IOStreams.Out, "╰─────────────────────────────────────────────────────────╯")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
}

// showSummaryAndConfirm displays the configuration summary and asks for confirmation
func (cmd *initCmd) showSummaryAndConfirm(templateName string, template Item) bool {
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "╭─────────────────────────────────────────────────────────╮")
	fmt.Fprintln(cmd.f.IOStreams.Out, "│  📋 Project Summary                                      │")
	fmt.Fprintln(cmd.f.IOStreams.Out, "├─────────────────────────────────────────────────────────┤")
	fmt.Fprintf(cmd.f.IOStreams.Out, "│  Name:        %-42s │\n", cmd.name)
	fmt.Fprintf(cmd.f.IOStreams.Out, "│  Template:    %-42s │\n", templateName)
	fmt.Fprintf(cmd.f.IOStreams.Out, "│  Location:    %-42s │\n", cmd.pathWorkingDir)
	if template.Message != "" {
		// Wrap long messages
		messageLines := wrapText(template.Message, 42)
		for i, line := range messageLines {
			if i == 0 {
				fmt.Fprintf(cmd.f.IOStreams.Out, "│  Description: %-42s │\n", line)
			} else {
				fmt.Fprintf(cmd.f.IOStreams.Out, "│                %-42s │\n", line)
			}
		}
	}
	fmt.Fprintln(cmd.f.IOStreams.Out, "╰─────────────────────────────────────────────────────────╯")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")

	return cmd.confirmProceed()
}

// confirmProceed asks the user to confirm proceeding with creation
func (cmd *initCmd) confirmProceed() bool {
	prompt := &survey.Confirm{
		Message: "Proceed with creation?",
		Default: true,
	}
	var proceed bool
	err := cmd.askOne(prompt, &proceed)
	if err != nil {
		return false
	}
	return proceed
}

// showNextSteps displays the next steps after project creation
func (cmd *initCmd) showNextSteps() {
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "╭─────────────────────────────────────────────────────────╮")
	fmt.Fprintln(cmd.f.IOStreams.Out, "│  🎉 Success! Your project is ready                      │")
	fmt.Fprintln(cmd.f.IOStreams.Out, "╰─────────────────────────────────────────────────────────╯")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintf(cmd.f.IOStreams.Out, "📁 Project created at: %s\n", cmd.pathWorkingDir)
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "🚀 Next steps:")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "  1. Navigate to your project:")
	fmt.Fprintf(cmd.f.IOStreams.Out, "     $ cd %s\n", cmd.name)
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "  2. Start development server:")
	fmt.Fprintln(cmd.f.IOStreams.Out, "     $ azion dev")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "  3. Deploy to Azion Edge:")
	fmt.Fprintln(cmd.f.IOStreams.Out, "     $ azion deploy")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "📚 Learn more:")
	fmt.Fprintln(cmd.f.IOStreams.Out, "  • Documentation: https://docs.azion.com")
	fmt.Fprintln(cmd.f.IOStreams.Out, "  • Examples: https://github.com/aziontech/examples")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
}

// wrapText wraps text to a specified width
func wrapText(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	currentLine := words[0]
	for _, word := range words[1:] {
		if len(currentLine)+1+len(word) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
