package build

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/tidwall/gjson"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type BuildCmd struct {
	Io        *iostreams.IOStreams
	WriteFile func(filename string, data []byte, perm fs.FileMode) error
	// Return output, exit code and any errors
	CommandRunner      func(cmd string, envvars []string) (string, int, error)
	FileReader         func(path string) ([]byte, error)
	ConfigRelativePath string
	GetWorkDir         func() (string, error)
	EnvLoader          func(path string) ([]string, error)
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

func newBuildCmd(f *cmdutil.Factory) *BuildCmd {
	return &BuildCmd{
		Io:         f.IOStreams,
		FileReader: os.ReadFile,
		CommandRunner: func(cmd string, envs []string) (string, int, error) {
			return utils.RunCommandWithOutput(envs, cmd)
		},
		ConfigRelativePath: "/azion/config.json",
		GetWorkDir:         utils.GetWorkingDir,
		EnvLoader:          utils.LoadEnvVarsFromFile,
		WriteFile:          os.WriteFile,
	}
}

func NewBuildCmd(f *cmdutil.Factory) *BuildCmd {
	return newBuildCmd(f)
}

// var runBuild = RunBuildCmdLine
func (cmd *BuildCmd) run() error {
	path, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	jsonConf := path + "/azion/azion.json"
	file, err := cmd.FileReader(jsonConf)
	if err != nil {
		return msg.ErrorOpeningAzionFile
	}

	typeLang := gjson.Get(string(file), "type")
	err = RunBuildCmdLine(cmd, typeLang.String())
	if err != nil {
		return err
	}

	return nil
}

func RunBuildCmdLine(cmd *BuildCmd, typeLang string) error {
	var output string
	var exitCode int
	var err error

	switch typeLang {
	case "javascript":
		output, exitCode, err = BuildJavascript(cmd)
		if err != nil {
			return err
		}
	case "nextjs":
		output, exitCode, err = BuildNextjs(cmd)
		if err != nil {
			return err
		}
	case "flareact":
		output, exitCode, err = BuildFlareact(cmd)
		if err != nil {
			return err
		}
	default:
		output = ""
		exitCode = 0
		err = utils.ErrorUnsupportedType
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Io.Out, "%s\n", output)
	fmt.Fprintf(cmd.Io.Out, msg.WebappOutput, exitCode)

	return nil
}

func (cmd *BuildCmd) Run() error {
	return cmd.run()
}

func BuildJavascript(cmd *BuildCmd) (string, int, error) {
	conf, err := getConfig(cmd)
	if err != nil {
		return "", 0, err
	}

	envs, err := cmd.EnvLoader(conf.InitData.Env)
	if err != nil {
		return "", 0, msg.ErrReadEnvFile
	}

	workDirPath, err := cmd.GetWorkDir()
	if err != nil {
		return "", 0, utils.ErrorInternalServerError
	}

	workDirPath += "/args.json"
	_, err = cmd.FileReader(workDirPath)
	if err != nil {
		if err := cmd.WriteFile(workDirPath, []byte("{}"), 0644); err != nil {
			return "", 0, fmt.Errorf(utils.ErrorCreateFile.Error(), workDirPath)
		}
	}

	fmt.Fprintf(cmd.Io.Out, msg.WebappBuildRunningCmd)
	fmt.Fprintf(cmd.Io.Out, "$ %s\n", conf.BuildData.Cmd)

	output, exitCode, err := cmd.CommandRunner(conf.BuildData.Cmd, envs)
	if err != nil {
		return "", exitCode, err
	}

	return output, exitCode, nil
}

func BuildNextjs(cmd *BuildCmd) (string, int, error) {
	return "", 0, nil
}

func BuildFlareact(cmd *BuildCmd) (string, int, error) {
	return "", 0, nil
}

func getConfig(cmd *BuildCmd) (conf *contracts.AzionApplicationConfig, err error) {
	path, err := utils.GetWorkingDir()
	if err != nil {
		return conf, err
	}

	jsonConf := path + "/azion/config.json"
	file, err := cmd.FileReader(jsonConf)
	if err != nil {
		return conf, msg.ErrorOpeningConfigFile
	}

	conf = &contracts.AzionApplicationConfig{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return conf, msg.ErrorUnmarshalConfigFile
	}

	if conf.BuildData.Cmd == "" {
		return conf, msg.ErrorWebappBuildCmdNotSpecified
	}

	return conf, nil
}
