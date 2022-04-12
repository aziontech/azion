package build

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type buildCmd struct {
	io *iostreams.IOStreams
	// Return output, exit code and any errors
	commandRunner      func(cmd string, envvars []string) (string, int, error)
	fileReader         func(path string) ([]byte, error)
	configRelativePath string
	getWorkDir         func() (string, error)
	envLoader          func(path string) ([]string, error)
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	build := &buildCmd{
		io:         f.IOStreams,
		fileReader: os.ReadFile,
		commandRunner: func(cmd string, envs []string) (string, int, error) {
			return utils.RunCommandWithOutput(envs, cmd)
		},
		configRelativePath: "/azion/config.json",
		getWorkDir:         utils.GetWorkingDir,
		envLoader:          utils.LoadEnvVarsFromFile,
	}

	cmd := &cobra.Command{
		Use:           "build [flags]",
		Short:         "Build your Web application",
		Long:          "Build your Web application",
		SilenceErrors: true,
		SilenceUsage:  true,
		Annotations: map[string]string{
			"Category": "Build",
		},
		Example: heredoc.Doc(`
        $ azioncli build
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.runCmd()
		},
	}
	return cmd
}

func (b *buildCmd) readConfig() (*contracts.AzionApplicationConfig, error) {
	path, err := b.getWorkDir()
	if err != nil {
		return nil, err
	}

	file, err := b.fileReader(path + b.configRelativePath)
	if err != nil {
		return nil, ErrOpeningConfigFile
	}

	conf := &contracts.AzionApplicationConfig{}

	if err := json.Unmarshal(file, &conf); err != nil {
		return nil, ErrUnmarshalConfigFile
	}

	return conf, nil
}

func (b *buildCmd) runCmd() error {
	conf, err := b.readConfig()
	if err != nil {
		return err
	}

	envs, err := b.envLoader(conf.BuildData.Env)
	if err != nil {
		return ErrReadEnvFile
	}

	fmt.Fprintf(b.io.Out, "Running build command\n\n")
	fmt.Fprintf(b.io.Out, "$ %s\n", conf.BuildData.Cmd)

	out, exitCode, err := b.commandRunner(conf.BuildData.Cmd, envs)

	fmt.Fprintf(b.io.Out, "%s\n", string(out))
	fmt.Fprintf(b.io.Out, "\nCommand exited with exit code %d\n", exitCode)

	if err != nil {
		return err
	}

	return nil
}
