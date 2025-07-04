package build

import (
	"io"
	"io/fs"
	"os"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/build"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/command"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type BuildCmd struct {
	Io                    *iostreams.IOStreams
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	CommandRunnerStream   func(out io.Writer, cmd string, envvars []string) error
	CommandRunInteractive func(f *cmdutil.Factory, comm string) error
	CommandRunner         func(f *cmdutil.Factory, comm string, envVars []string) (string, error)
	FileReader            func(path string) ([]byte, error)
	GetAzionJsonContent   func(pathConf string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confPath string) error
	EnvLoader             func(path string) ([]string, error)
	Stat                  func(path string) (fs.FileInfo, error)
	GetWorkDir            func() (string, error)
	f                     *cmdutil.Factory
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewBuildCmd(f))
}

func NewCobraCmd(build *BuildCmd) *cobra.Command {
	fields := &contracts.BuildInfo{}
	buildCmd := &cobra.Command{
		Use:           msg.BuildUsage,
		Short:         msg.BuildShortDescription,
		Long:          msg.BuildLongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example:       heredoc.Doc("\n$ azion build\n"),
		RunE: func(cmd *cobra.Command, args []string) error {
			msgs := []string{}
			return build.run(fields, &msgs)
		},
	}

	buildCmd.Flags().BoolP("help", "h", false, msg.BuildFlagHelp)
	buildCmd.Flags().StringVar(&fields.Preset, "preset", "", msg.FlagTemplate)
	buildCmd.Flags().StringVar(&fields.Entry, "entry", "", msg.FlagEntry)
	buildCmd.Flags().StringVar(&fields.NodePolyfills, "use-node-polyfills", "", msg.FlagPolyfill)
	buildCmd.Flags().StringVar(&fields.OwnWorker, "use-own-worker", "", msg.FlagWorker)
	buildCmd.Flags().StringVar(&fields.ProjectPath, "config-dir", "azion", msg.ProjectConfFlag)
	buildCmd.Flags().BoolVar(&fields.SkipFramework, "skip-framework-build", false, msg.SkipFrameworkBuild)

	return buildCmd
}

func NewBuildCmd(f *cmdutil.Factory) *BuildCmd {
	return &BuildCmd{
		Io:         f.IOStreams,
		FileReader: os.ReadFile,
		CommandRunnerStream: func(out io.Writer, cmd string, envs []string) error {
			return command.RunCommandStreamOutput(f.IOStreams.Out, envs, cmd)
		},
		CommandRunInteractive: func(f *cmdutil.Factory, comm string) error {
			return command.CommandRunInteractive(f, comm)
		},
		CommandRunner: func(f *cmdutil.Factory, comm string, envVars []string) (string, error) {
			return command.CommandRunInteractiveWithOutput(f, comm, envVars)
		},
		EnvLoader:             utils.LoadEnvVarsFromFile,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		WriteFile:             os.WriteFile,
		Stat:                  os.Stat,
		GetWorkDir:            utils.GetWorkingDir,
		f:                     f,
	}
}

func (b *BuildCmd) ExternalRun(fields *contracts.BuildInfo, confPath string, msgs *[]string, skipFramework bool) error {
	fields.ProjectPath = confPath
	fields.SkipFramework = skipFramework
	return b.run(fields, msgs)
}
