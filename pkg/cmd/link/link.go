package link

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/link"
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

type LinkInfo struct {
	Name           string
	Template       string
	Mode           string
	PathWorkingDir string
	GlobalFlagAll  bool
}

type LinkCmd struct {
	Io              *iostreams.IOStreams
	GetWorkDir      func() (string, error)
	FileReader      func(path string) ([]byte, error)
	LookPath        func(bin string) (string, error)
	IsDirEmpty      func(dirpath string) (bool, error)
	CleanDir        func(dirpath string) error
	WriteFile       func(filename string, data []byte, perm fs.FileMode) error
	OpenFile        func(name string) (*os.File, error)
	RemoveAll       func(path string) error
	Rename          func(oldpath string, newpath string) error
	CreateTempDir   func(dir string, pattern string) (string, error)
	EnvLoader       func(path string) ([]string, error)
	Stat            func(path string) (fs.FileInfo, error)
	Mkdir           func(path string, perm os.FileMode) error
	GitPlainClone   func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
	CommandRunner   func(cmd string, envvars []string) (string, int, error)
	ShouldConfigure func(cmd *LinkCmd, info *LinkInfo) (bool, error)
}

func NewLinkCmd(f *cmdutil.Factory) *LinkCmd {
	return &LinkCmd{
		Io:              f.IOStreams,
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
		CommandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
	}
}

func NewCobraCmd(init *LinkCmd, f *cmdutil.Factory) *cobra.Command {
	options := &contracts.AzionApplicationOptions{}
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
		$ azion link --name "thisisatest" --template nextjs
		$ azion link --name "thisisatest" --template static
		$ azion link --name "thisisatest" --template hexo --mode deliver
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			info.GlobalFlagAll = f.GlobalFlagAll
			return init.run(info, options, cmd)
		},
	}

	cobraCmd.Flags().StringVar(&info.Name, "name", "", msg.EdgeApplicationsLinkFlagName)
	cobraCmd.Flags().StringVar(&info.Template, "template", "", msg.EdgeApplicationsLinkFlagTemplate)
	cobraCmd.Flags().StringVar(&info.Mode, "mode", "", msg.EdgeApplicationsLinkFlagMode)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewLinkCmd(f), f)
}

func (cmd *LinkCmd) run(info *LinkInfo, options *contracts.AzionApplicationOptions, c *cobra.Command) error {
	logger.Debug("Running link subcommand from edge_applications command tree")

	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return err
	}
	info.PathWorkingDir = path

	shouldLink, err := cmd.ShouldConfigure(cmd, info)
	if err != nil {
		return err
	}
	if !shouldLink {
		return nil
	}

	switch info.Template {
	case "simple":
		return initSimple(cmd, path, info, c)
	case "static":
		return initStatic(cmd, info, options, c)
	}

	shouldFetchTemplates, err := shouldFetch(cmd, info)
	if err != nil {
		return err
	}

	if shouldFetchTemplates {
		if info.GlobalFlagAll {
			info.Name = thoth.GenerateName()
		} else {
			if !c.Flags().Changed("name") {
				projName, err := askForInput(msg.LinkProjectQuestion, thoth.GenerateName())
				if err != nil {
					return err
				}

				info.Name = projName
			}
		}

		if !c.Flags().Changed("template") || !c.Flags().Changed("mode") {
			err = cmd.selectVulcanTemplates(info)
			if err != nil {
				return err
			}
		}

		if err = cmd.createTemplateAzion(info); err != nil {
			return err
		}

		logger.FInfo(cmd.Io.Out, msg.WebAppLinkCmdSuccess)
		logger.FInfo(cmd.Io.Out, fmt.Sprintf(msg.EdgeApplicationsLinkSuccessful, info.Name))
	}

	err = InitNextjs(info, cmd)
	if err != nil {
		return err
	}

	return nil
}
