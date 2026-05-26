package install

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/install"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/output"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type installCmd struct {
	f      *cmdutil.Factory
	skills bool

	// Dependency injection for testability
	userHomeDir   func() string
	executableDir func() (string, error)
	readDir       func(name string) ([]os.DirEntry, error)
	mkdirAll      func(path string, perm os.FileMode) error
	removeAll     func(path string) error
	copyDir       func(src, dst string) error
	stat          func(name string) (os.FileInfo, error)
	getWorkingDir func() (string, error)
}

func NewInstallCmd(f *cmdutil.Factory) *installCmd {
	return &installCmd{
		f:             f,
		userHomeDir:   getUserHomeDir,
		executableDir: getExecutableDir,
		readDir:       os.ReadDir,
		mkdirAll:      os.MkdirAll,
		removeAll:     os.RemoveAll,
		copyDir:       utils.CopyDirectory,
		stat:          os.Stat,
		getWorkingDir: utils.GetWorkingDir,
	}
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := NewInstallCmd(f)

	cobraCmd := &cobra.Command{
		Use:           msg.Usage,
		Short:         msg.ShortDescription,
		Long:          msg.LongDescription,
		Example:       heredoc.Doc(msg.Example),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          cmd.Run,
	}

	cobraCmd.Flags().BoolVar(&cmd.skills, "skills", false, msg.FlagSkills)
	cobraCmd.Flags().BoolP("help", "h", false, msg.FlagHelp)

	return cobraCmd
}

func (cmd *installCmd) Run(cobraCmd *cobra.Command, _ []string) error {
	if !cmd.skills {
		return cobraCmd.Help()
	}

	return cmd.installSkills()
}

// installSkills performs the skills installation
func (cmd *installCmd) installSkills() error {
	out := cmd.f.IOStreams.Out

	// Step 1: Resolve home directory
	logger.FInfo(out, msg.MsgResolveHome)
	homeDir := cmd.userHomeDir()
	if homeDir == "" {
		logger.Debug("Failed to resolve home directory")
		return msg.ErrorResolveHomeDir
	}

	// Step 2: Find source skills directory
	logger.FInfo(out, msg.MsgValidateSource)
	sourceDir, err := cmd.findSkillsSourceDir()
	if err != nil {
		logger.Debug("Failed to find skills source directory", zap.Error(err))
		return msg.ErrorSourceDirNotFound
	}

	// Step 3: Read available skills
	entries, err := cmd.readDir(sourceDir)
	if err != nil {
		logger.Debug("Failed to read source directory", zap.Error(err))
		return msg.ErrorReadSourceDir
	}

	// Filter to only directories (skills)
	var skillDirs []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "." && entry.Name() != ".." {
			skillDirs = append(skillDirs, entry)
		}
	}

	if len(skillDirs) == 0 {
		logger.FInfo(out, msg.MsgNoSkills)
		return nil
	}

	// Step 4: Create target base directory
	targetBase := filepath.Join(homeDir, ".claude", "skills")
	logger.Debug("Creating target directory", zap.String("path", targetBase))

	if err := cmd.mkdirAll(targetBase, 0755); err != nil {
		logger.Debug("Failed to create target directory", zap.Error(err))
		return fmt.Errorf(msg.ErrorCreateTargetDir.Error(), targetBase)
	}

	// Step 5: Install each skill
	installedCount := 0
	for _, entry := range skillDirs {
		skillName := entry.Name()
		sourceSkillPath := filepath.Join(sourceDir, skillName)
		targetSkillPath := filepath.Join(targetBase, skillName)

		// Check if skill already exists (clean overwrite)
		if info, err := cmd.stat(targetSkillPath); err == nil && info.IsDir() {
			logger.FInfo(out, fmt.Sprintf(msg.MsgRemoveExisting, skillName))
			if err := cmd.removeAll(targetSkillPath); err != nil {
				logger.Debug("Failed to remove existing skill", zap.Error(err))
				return fmt.Errorf(msg.ErrorRemoveExisting.Error(), targetSkillPath)
			}
		}

		// Copy skill
		logger.FInfo(out, fmt.Sprintf(msg.MsgCopySkill, skillName))
		if err := cmd.copyDir(sourceSkillPath, targetSkillPath); err != nil {
			logger.Debug("Failed to copy skill", zap.String("skill", skillName), zap.Error(err))
			return fmt.Errorf(msg.ErrorCopySkill.Error(), skillName, err)
		}

		installedCount++
	}

	// Step 6: Success output
	successMsg := fmt.Sprintf(msg.MsgDone, installedCount)
	logger.FInfo(out, successMsg)

	outputOut := output.GeneralOutput{
		Msg:   successMsg,
		Out:   out,
		Flags: cmd.f.Flags,
	}
	return output.Print(&outputOut)
}

// findSkillsSourceDir locates the bundled skills directory
func (cmd *installCmd) findSkillsSourceDir() (string, error) {
	// Strategy 1: Relative to executable (for production builds)
	// The skills directory is expected to be alongside the binary
	execDir, err := cmd.executableDir()
	if err == nil {
		candidate := filepath.Join(execDir, "..", "skills")
		if info, err := cmd.stat(candidate); err == nil && info.IsDir() {
			return filepath.Clean(candidate), nil
		}
	}

	// Strategy 2: Development path (relative to working directory)
	// This allows running `go run . install --skills` during development
	cwd, err := cmd.getWorkingDir()
	if err == nil {
		candidate := filepath.Join(cwd, "skills")
		if info, err := cmd.stat(candidate); err == nil && info.IsDir() {
			return candidate, nil
		}
	}

	return "", msg.ErrorSourceDirNotFound
}

// getUserHomeDir returns the user's home directory in a cross-platform way
func getUserHomeDir() string {
	env := "HOME"
	switch runtime.GOOS {
	case "windows":
		env = "USERPROFILE"
	case "plan9":
		env = "home"
	}
	if v := os.Getenv(env); v != "" {
		return v
	}
	switch runtime.GOOS {
	case "android":
		return "/sdcard"
	case "ios":
		return "/"
	}
	return ""
}

// getExecutableDir returns the directory of the current executable
func getExecutableDir() (string, error) {
	exec, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exec), nil
}
