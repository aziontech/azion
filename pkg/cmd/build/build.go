package build

import (
	"io"
	"io/fs"
	"os"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/build"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type BuildCmd struct {
	Io                  *iostreams.IOStreams
	WriteFile           func(filename string, data []byte, perm fs.FileMode) error
	CommandRunnerStream func(out io.Writer, cmd string, envvars []string) error
	FileReader          func(path string) ([]byte, error)
	ConfigRelativePath  string
	GetWorkDir          func() (string, error)
	EnvLoader           func(path string) ([]string, error)
	Stat                func(path string) (fs.FileInfo, error)
	VersionID           func(dir string) string
	f                   *cmdutil.Factory
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(newBuildCmd(f))
}

func NewCobraCmd(build *BuildCmd) *cobra.Command {
	buildCmd := &cobra.Command{
		Use:           msg.EdgeApplicationsBuildUsage,
		Short:         msg.EdgeApplicationsBuildShortDescription,
		Long:          msg.EdgeApplicationsBuildLongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example:       heredoc.Doc("\n$ azion build\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.run()
		},
	}

	buildCmd.Flags().BoolP("help", "h", false, msg.EdgeApplicationsBuildFlagHelp)

	return buildCmd
}

func newBuildCmd(f *cmdutil.Factory) *BuildCmd {
	return &BuildCmd{
		Io:         f.IOStreams,
		FileReader: os.ReadFile,
		CommandRunnerStream: func(out io.Writer, cmd string, envs []string) error {
			return utils.RunCommandStreamOutput(f.IOStreams.Out, envs, cmd)
		},
		ConfigRelativePath: "/azion/config.json",
		GetWorkDir:         utils.GetWorkingDir,
		EnvLoader:          utils.LoadEnvVarsFromFile,
		WriteFile:          os.WriteFile,
		Stat:               os.Stat,
		f:                  f,
		VersionID:          createVersionID,
	}
}

func (cmd *BuildCmd) Run() error {
	return cmd.run()
}
