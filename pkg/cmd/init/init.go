package init

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/init"
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

var (
	TemplateBranch = "main"
	TemplateMajor  = "0"
)

type InitInfo struct {
	Name           string
	Template       string
	Mode           string
	PathWorkingDir string
	YesOption      bool
	NoOption       bool
}

type InitCmd struct {
	Io            *iostreams.IOStreams
	GetWorkDir    func() (string, error)
	FileReader    func(path string) ([]byte, error)
	LookPath      func(bin string) (string, error)
	IsDirEmpty    func(dirpath string) (bool, error)
	CleanDir      func(dirpath string) error
	WriteFile     func(filename string, data []byte, perm fs.FileMode) error
	OpenFile      func(name string) (*os.File, error)
	RemoveAll     func(path string) error
	Rename        func(oldpath string, newpath string) error
	CreateTempDir func(dir string, pattern string) (string, error)
	EnvLoader     func(path string) ([]string, error)
	Stat          func(path string) (fs.FileInfo, error)
	Mkdir         func(path string, perm os.FileMode) error
	GitPlainClone func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	CommandRunner func(cmd string, envvars []string) (string, int, error)
}

func NewInitCmd(f *cmdutil.Factory) *InitCmd {
	return &InitCmd{
		Io:            f.IOStreams,
		GetWorkDir:    utils.GetWorkingDir,
		FileReader:    os.ReadFile,
		LookPath:      exec.LookPath,
		IsDirEmpty:    utils.IsDirEmpty,
		CleanDir:      utils.CleanDirectory,
		WriteFile:     os.WriteFile,
		OpenFile:      os.Open,
		RemoveAll:     os.RemoveAll,
		Rename:        os.Rename,
		CreateTempDir: os.MkdirTemp,
		EnvLoader:     utils.LoadEnvVarsFromFile,
		Stat:          os.Stat,
		Mkdir:         os.MkdirAll,
		GitPlainClone: git.PlainClone,
		CommandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
	}
}

func NewCobraCmd(init *InitCmd) *cobra.Command {
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
		$ azioncli edge_applications init --name "thisisatest" --type static
		$ azion init --name "thisisatest" --template hexo --mode deliver
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return init.run(info, options, cmd)
		},
	}

	cobraCmd.Flags().StringVar(&info.Name, "name", "", msg.EdgeApplicationsInitFlagName)
	cobraCmd.Flags().StringVar(&info.Template, "template", "", msg.EdgeApplicationsInitFlagTemplate)
	cobraCmd.Flags().StringVar(&info.Mode, "mode", "", msg.EdgeApplicationsInitFlagMode)
	cobraCmd.Flags().BoolVarP(&info.YesOption, "yes", "y", false, msg.EdgeApplicationsInitFlagYes)
	cobraCmd.Flags().BoolVarP(&info.NoOption, "no", "n", false, msg.EdgeApplicationsInitFlagNo)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewInitCmd(f))
}

func (cmd *InitCmd) run(info *InitInfo, options *contracts.AzionApplicationOptions, c *cobra.Command) error {
	logger.Debug("Running init subcommand from edge_applications command tree")
	if info.YesOption && info.NoOption {
		return msg.ErrorYesAndNoOptions
	}

	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}
	info.PathWorkingDir = path

	if !c.Flags().Changed("template") {
		err = cmd.selectVulcanTemplates(info)
		if err != nil {
			return err
		}
	}

	switch info.Template {
	case "simple":
		return initSimple(cmd, path, info, c)
	case "static":
		return initStatic(cmd, info, options, c)
	}

	if (!c.Flags().Changed("mode") || !c.Flags().Changed("template")) && info.Template != "nextjs" {
		return msg.ErrorModeNotSent
	}

	shouldFetchTemplates, err := shouldFetch(cmd, info)
	if err != nil {
		return err
	}

	if shouldFetchTemplates {
		if info.YesOption {
			info.Name = thoth.GenerateName()
		} else {
			if !c.Flags().Changed("name") {
				projName, err := askForInput(msg.InitProjectQuestion, thoth.GenerateName())
				if err != nil {
					return err
				}

				info.Name = projName
			}
		}

		if err = cmd.createTemplateAzion(info); err != nil {
			return err
		}

		logger.FInfo(cmd.Io.Out, msg.WebAppInitCmdSuccess)
		logger.FInfo(cmd.Io.Out, fmt.Sprintf(msg.EdgeApplicationsInitSuccessful, info.Name))
	}

	err = InitNextjs(info, cmd)
	if err != nil {
		return err
	}

	return nil
}
