package build

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type buildCmd struct {
	io        *iostreams.IOStreams
	writeFile func(filename string, data []byte, perm fs.FileMode) error
	// Return output, exit code and any errors
	commandRunner      func(cmd string, envvars []string) (string, int, error)
	fileReader         func(path string) ([]byte, error)
	configRelativePath string
	getWorkDir         func() (string, error)
	envLoader          func(path string) ([]string, error)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	command := newBuildCmd(f)
	buildCmd := &cobra.Command{
		Use:           msg.WebappBuildUsage,
		Short:         msg.WebappBuildShortDescription,
		Long:          msg.WebappBuildLongDescription,
		SilenceErrors: true,
		SilenceUsage:  true,
		Example: heredoc.Doc(`
        $ azioncli webapp build
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return command.run()
		},
	}

	buildCmd.Flags().BoolP("help", "h", false, msg.WebappBuildFlagHelp)

	return buildCmd
}

func newBuildCmd(f *cmdutil.Factory) *buildCmd {
	return &buildCmd{
		io:         f.IOStreams,
		fileReader: os.ReadFile,
		commandRunner: func(cmd string, envs []string) (string, int, error) {
			return utils.RunCommandWithOutput(envs, cmd)
		},
		configRelativePath: "/azion/config.json",
		getWorkDir:         utils.GetWorkingDir,
		envLoader:          utils.LoadEnvVarsFromFile,
		writeFile:          os.WriteFile,
	}
}

func NewBuildCmd(f *cmdutil.Factory) *buildCmd {
	return newBuildCmd(f)
}

func (cmd *buildCmd) readConfig() (*contracts.AzionApplicationConfig, error) {
	path, err := cmd.getWorkDir()
	if err != nil {
		return nil, err
	}

	file, err := cmd.fileReader(path + cmd.configRelativePath)
	if err != nil {
		return nil, msg.ErrOpeningConfigFile
	}

	conf := &contracts.AzionApplicationConfig{}

	if err := json.Unmarshal(file, &conf); err != nil {
		return nil, msg.ErrUnmarshalConfigFile
	}

	return conf, nil
}

func (cmd *buildCmd) run() error {
	conf, err := cmd.readConfig()
	if err != nil {
		return err
	}

	envs, err := cmd.envLoader(conf.BuildData.Env)
	if err != nil {
		return msg.ErrReadEnvFile
	}

	if conf.BuildData.Cmd == "" {
		fmt.Fprintf(cmd.io.Out, msg.WebappBuildCmdNotSpecified)
		return nil
	}

	workDirPath, err := cmd.getWorkDir()

	workDirPath += "/args.json"
	_, err = cmd.fileReader(workDirPath)
	if err != nil {
		cmd.writeFile(workDirPath, []byte("{}"), 0644)
	}

	cmdRunner := "npx --yes --package=webpack@5.72.0 --package=webpack-cli@4.9.2 -- webpack --config ./azion/webpack.config.js -o ${OUTPUT_DIR} --mode production || exit $? ;;"
	fmt.Fprintf(cmd.io.Out, msg.WebappBuildRunningCmd)
	fmt.Fprintf(cmd.io.Out, "$ %s\n", cmdRunner)

	out, exitCode, err := cmd.commandRunner(cmdRunner, envs)

	fmt.Fprintf(cmd.io.Out, "%s\n", out)
	fmt.Fprintf(cmd.io.Out, msg.WebappOutput, exitCode)

	if err != nil {
		return msg.ErrFailedToRunCommand
	}

	return nil
}

func (cmd *buildCmd) Run() error {
	return cmd.run()
}
