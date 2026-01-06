package init

import (
	"fmt"
	"path"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/logger"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"go.uber.org/zap"
)

func (cmd *initCmd) selectVulcanTemplates(vul *vulcanPkg.VulcanPkg) error {
	logger.Debug("Running bundler store init")
	// checking if vulcan major is correct
	vulcanVer, err := cmd.commandRunnerOutput(cmd.f, "npm show edge-functions version", []string{})
	if err != nil {
		return err
	}

	err = vul.CheckVulcanMajor(vulcanVer, cmd.f, vul)
	if err != nil {
		return err
	}

	// TODO: use later
	cmdVulcanBuild := "build"
	if len(cmd.preset) > 0 {
		cmdVulcanBuild = fmt.Sprintf("%s --preset '%s' --only-generate-config", cmdVulcanBuild, cmd.preset)
	}

	command := vul.Command("", cmdVulcanBuild, cmd.f)
	logger.Debug("Running the following command", zap.Any("Command", command))

	err = cmd.commandRunInteractive(cmd.f, command)
	if err != nil {
		return err
	}

	preset, err := cmd.getVulcanInfo()
	if err != nil {
		return err
	}

	cmd.preset = strings.ToLower(preset)

	if cmd.preset == "vite" {
		cmd.preset = "vue"
	}

	return nil
}

func (cmd *initCmd) depsInstall() error {
	command := fmt.Sprintf("%s install", cmd.packageManager)

	// Note: We don't use a spinner here because npm/yarn/pnpm install
	// produces useful output that users want to see in real-time
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

	var infoJson Item
	err = cmd.unmarshal(fileContent, &infoJson)
	if err != nil {
		logger.Debug("Error unmarshalling template info", zap.Error(err))
		return "", err
	}

	logger.Debug("Information about the template:", zap.Any("preset", infoJson.Preset))
	return infoJson.Preset, nil
}

// showWelcome displays the welcome message
func (cmd *initCmd) showWelcome() {
	titleStyle := GetAzionTitleStyle() // Purple for welcome title
	
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "╭─────────────────────────────────────────────────────────╮")
	fmt.Fprintf(cmd.f.IOStreams.Out, "│  %s                       │\n", titleStyle.Render("🚀 Welcome to Azion Web Platform"))
	fmt.Fprintln(cmd.f.IOStreams.Out, "│  Let's create your web application                      │")
	fmt.Fprintln(cmd.f.IOStreams.Out, "╰─────────────────────────────────────────────────────────╯")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
}

// showNextSteps displays the next steps after project creation
func (cmd *initCmd) showNextSteps() {
	// Styles using Azion theme colors
	successStyle := GetAzionSuccessStyle()  // Orange for success message
	labelStyle := GetAzionLabelStyle()      // Purple for labels/headings
	pathStyle := GetAzionAnswerStyle()      // Orange for paths/commands
	
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "╭─────────────────────────────────────────────────────────╮")
	fmt.Fprintf(cmd.f.IOStreams.Out, "│  %s                      │\n", successStyle.Render("🎉 Success! Your project is ready"))
	fmt.Fprintln(cmd.f.IOStreams.Out, "╰─────────────────────────────────────────────────────────╯")
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintf(cmd.f.IOStreams.Out, "%s %s\n", labelStyle.Render("📁 Project created at:"), pathStyle.Render(cmd.pathWorkingDir))
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, labelStyle.Render("🚀 Next steps:"))
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "  1. Navigate to your project:")
	fmt.Fprintf(cmd.f.IOStreams.Out, "     %s\n", pathStyle.Render("$ cd "+cmd.name))
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "  2. Start development server:")
	fmt.Fprintf(cmd.f.IOStreams.Out, "     %s\n", pathStyle.Render("$ azion dev"))
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, "  3. Deploy to Azion Edge:")
	fmt.Fprintf(cmd.f.IOStreams.Out, "     %s\n", pathStyle.Render("$ azion deploy"))
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
	fmt.Fprintln(cmd.f.IOStreams.Out, labelStyle.Render("📚 Learn more:"))
	fmt.Fprintf(cmd.f.IOStreams.Out, "  • Documentation: %s\n", pathStyle.Render("https://docs.azion.com"))
	fmt.Fprintf(cmd.f.IOStreams.Out, "  • Examples: %s\n", pathStyle.Render("https://github.com/aziontech/examples"))
	fmt.Fprintln(cmd.f.IOStreams.Out, "")
}

// matchesCategory determines if a template matches the selected category
func (cmd *initCmd) matchesCategory(category string, item Item, presetName string) bool {
	preset := strings.ToLower(item.Preset)
	name := strings.ToLower(item.Name)

	switch category {
	case "Simple Hello World":
		// Simple templates - basic starters
		return strings.Contains(name, "hello") ||
			strings.Contains(name, "basic") ||
			strings.Contains(name, "starter") ||
			preset == "javascript" && strings.Contains(name, "simple")

	case "JavaScript":
		// JavaScript-specific templates
		return preset == "javascript" &&
			!strings.Contains(name, "typescript") &&
			!cmd.isFrameworkTemplate(preset, name)

	case "TypeScript":
		// TypeScript-specific templates
		return preset == "typescript" ||
			strings.Contains(name, "typescript") ||
			strings.Contains(preset, "typescript")

	case "Frameworks":
		// Framework templates (React, Next, Astro, Vue, Angular, etc.)
		return cmd.isFrameworkTemplate(preset, name)

	default:
		return true
	}
}

// isFrameworkTemplate checks if a template is a framework
func (cmd *initCmd) isFrameworkTemplate(preset, name string) bool {
	frameworks := []string{
		"react", "next", "nextjs", "astro", "vue", "nuxt",
		"angular", "svelte", "solid", "qwik", "remix",
		"gatsby", "eleventy", "hexo", "hugo", "jekyll",
		"vite", "vuepress", "vitepress", "docusaurus",
	}

	presetLower := strings.ToLower(preset)
	nameLower := strings.ToLower(name)

	for _, fw := range frameworks {
		if strings.Contains(presetLower, fw) || strings.Contains(nameLower, fw) {
			return true
		}
	}
	return false
}

// normalizeFrameworkName normalizes framework names (e.g., NextJs and OpenNext both become "Next.js")
func (cmd *initCmd) normalizeFrameworkName(preset string) string {
	presetLower := strings.ToLower(preset)

	// Merge NextJs and OpenNext into "Next.js"
	if strings.Contains(presetLower, "next") {
		return "Next.js"
	}
	if strings.Contains(presetLower, "react") {
		return "React"
	}
	if strings.Contains(presetLower, "astro") {
		return "Astro"
	}
	if strings.Contains(presetLower, "vue") {
		return "Vue"
	}
	if strings.Contains(presetLower, "angular") {
		return "Angular"
	}
	if strings.Contains(presetLower, "svelte") {
		return "Svelte"
	}
	if strings.Contains(presetLower, "solid") {
		return "Solid"
	}
	if strings.Contains(presetLower, "qwik") {
		return "Qwik"
	}
	if strings.Contains(presetLower, "remix") {
		return "Remix"
	}
	if strings.Contains(presetLower, "gatsby") {
		return "Gatsby"
	}
	if strings.Contains(presetLower, "eleventy") {
		return "Eleventy"
	}
	if strings.Contains(presetLower, "hexo") {
		return "Hexo"
	}
	if strings.Contains(presetLower, "hugo") {
		return "Hugo"
	}
	if strings.Contains(presetLower, "jekyll") {
		return "Jekyll"
	}
	if strings.Contains(presetLower, "vite") {
		return "Vite"
	}

	// Capitalize first letter for unknown frameworks
	if len(preset) > 0 {
		return strings.ToUpper(preset[:1]) + strings.ToLower(preset[1:])
	}
	return preset
}

// matchesFramework checks if a template matches the selected framework
func (cmd *initCmd) matchesFramework(selectedFramework string, item Item) bool {
	// Normalize the item's preset to compare
	itemFramework := cmd.normalizeFrameworkName(item.Preset)

	return itemFramework == selectedFramework
}
