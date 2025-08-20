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
	GetAzionJsonContent   func(pathConf string) (*contracts.AzionApplicationOptionsV3, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptionsV3, confPath string) error
	EnvLoader             func(path string) ([]string, error)
	Stat                  func(path string) (fs.FileInfo, error)
	GetWorkDir            func() (string, error)
	f                     *cmdutil.Factory
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewBuildCmd(f))
}

func NewCobraCmd(build *BuildCmd) *cobra.Command {
	fields := &contracts.BuildInfoV3{}
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
	buildCmd.Flags().BoolVar(&fields.IsFirewall, "firewall", false, msg.IsFirewall)
	buildCmd.Flags().StringVar(&fields.ProjectPath, "config-dir", "azion", msg.ProjectConfFlag)

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
		GetAzionJsonContent:   utils.GetAzionJsonContentV3,
		WriteAzionJsonContent: utils.WriteAzionJsonContentV3,
		WriteFile:             os.WriteFile,
		Stat:                  os.Stat,
		GetWorkDir:            utils.GetWorkingDir,
		f:                     f,
	}
}

func (b *BuildCmd) ExternalRun(fields *contracts.BuildInfoV3, confPath string, msgs *[]string) error {
	fields.ProjectPath = confPath
	return b.run(fields, msgs)
}
