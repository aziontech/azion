package build

import (
	"io"
	"io/fs"
	"os"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/build"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var Preset string
var Mode string

type BuildCmd struct {
	Io                    *iostreams.IOStreams
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	CommandRunnerStream   func(out io.Writer, cmd string, envvars []string) error
	CommandRunInteractive func(f *cmdutil.Factory, comm string) error
	FileReader            func(path string) ([]byte, error)
	ConfigRelativePath    string
	GetAzionJsonContent   func() (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions) error
	EnvLoader             func(path string) ([]string, error)
	Stat                  func(path string) (fs.FileInfo, error)
	VersionID             func() string
	f                     *cmdutil.Factory
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewBuildCmd(f))
}

func NewCobraCmd(build *BuildCmd) *cobra.Command {
	buildCmd := &cobra.Command{
		Use:           msg.BuildUsage,
		Short:         msg.BuildShortDescription,
		Long:          msg.BuildLongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example:       heredoc.Doc("\n$ azion build\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.run()
		},
	}

	buildCmd.Flags().BoolP("help", "h", false, msg.BuildFlagHelp)
	buildCmd.Flags().StringVar(&Preset, "preset", "", msg.FlagTemplate)
	buildCmd.Flags().StringVar(&Mode, "mode", "", msg.FlagMode)

	return buildCmd
}

func NewBuildCmd(f *cmdutil.Factory) *BuildCmd {
	return &BuildCmd{
		Io:         f.IOStreams,
		FileReader: os.ReadFile,
		CommandRunnerStream: func(out io.Writer, cmd string, envs []string) error {
			return utils.RunCommandStreamOutput(f.IOStreams.Out, envs, cmd)
		},
		CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
			return utils.CommandRunInteractive(f, comm)
		},
		ConfigRelativePath:    "/azion/config.json",
		EnvLoader:             utils.LoadEnvVarsFromFile,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		WriteFile:             os.WriteFile,
		Stat:                  os.Stat,
		f:                     f,
		VersionID:             createVersionID,
	}
}

func (cmd *BuildCmd) Run() error {
	return cmd.run()
}
