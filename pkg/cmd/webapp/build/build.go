package build

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/fs"
	"os"

	"errors"
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

func (cmd *BuildCmd) readConfig() (*contracts.AzionApplicationConfig, error) {
	path, err := cmd.GetWorkDir()
	if err != nil {
		return nil, err
	}

	file, err := cmd.FileReader(path + cmd.ConfigRelativePath)
	if err != nil {
		return nil, msg.ErrOpeningConfigFile
	}

	conf := &contracts.AzionApplicationConfig{}

	if err := json.Unmarshal(file, &conf); err != nil {
		return nil, msg.ErrUnmarshalConfigFile
	}

	return conf, nil
}

func (cmd *BuildCmd) run() error {
	path, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	jsonConf := path + "/azion/azion.json"
	file, err := cmd.FileReader(jsonConf)
	if err != nil {
		return err
	}

	typeLang := gjson.Get(string(file), "type")

	err = cmd.runInitCmdLine(typeLang.String())
	if err != nil {
		return err
	}

	return nil
}

func (cmd *BuildCmd) runInitCmdLine(typeLang string) error {
	var output string
	var exitCode int
	var err error

	switch typeLang {
	case "javascript":
		output, exitCode, err = BuildJavascript(cmd)
		if err != nil {
			return errors.New("failed Building err: " + err.Error())
		}
	case "nextjs":
		output, exitCode, err = BuildNextjs(cmd)
		if err != nil {
			return errors.New("failed Building err: " + err.Error())
		}
	case "flareact":
		output, exitCode, err = BuildFlareact(cmd)
		if err != nil {
			return errors.New("failed Building err: " + err.Error())
		}
	default:
		output = ""
		exitCode = 0
		err = errors.New("setp invalid")
	}

	fmt.Fprintf(cmd.Io.Out, "%s\n", output)
	fmt.Fprintf(cmd.Io.Out, msg.WebappOutput, exitCode)

	if err != nil {
		return err
	}

	return nil
}

func (cmd *BuildCmd) Run() error {
	return cmd.run()
}

func BuildJavascript(cmd *BuildCmd) (string, int, error) {
	conf, err := getConfig()
	if err != nil {
		fmt.Println("getConfig :", err.Error())
		return "", 0, err
	}

	envs, err := cmd.EnvLoader(conf.InitData.Env)
	if err != nil {
		fmt.Println("EnvLoader :", err.Error())
		return "", 0, errors.New("failed load envs err: " + err.Error())
	}

	workDirPath, err := cmd.GetWorkDir()
	if err != nil {
		fmt.Println("GetWorkDir :", err.Error())
		return "", 0, errors.New("failed getworkdir err: " + err.Error())
	}

	workDirPath += "/args.json"
	_, err = cmd.FileReader(workDirPath)
	if err != nil {
		fmt.Println("FileReader :", err.Error())
		cmd.WriteFile(workDirPath, []byte("{}"), 0644)
	}

	fmt.Fprintf(cmd.Io.Out, msg.WebappBuildRunningCmd)
	fmt.Fprintf(cmd.Io.Out, "$ %s\n", conf.BuildData.Cmd)

	cmdRunner := "npx --yes --package=webpack@5.72.0 --package=webpack-cli@4.9.2 -- webpack --config ./azion/webpack.config.js -o ./worker --mode production || exit $? ;" // conf.BuildData.Cmd
	output, exitCode, err := cmd.CommandRunner(cmdRunner, envs)
	if err != nil {
		fmt.Println("CommandRunner :", err.Error())
		return "", exitCode, errors.New("failed commandRunner err: " + err.Error())
	}

	return output, exitCode, nil
}

func BuildNextjs(cmd *BuildCmd) (string, int, error) {
	return "", 0, nil
}

func BuildFlareact(cmd *BuildCmd) (string, int, error) {
	return "", 0, nil
}

func getConfig() (conf *contracts.AzionApplicationConfig, err error) {
	path, err := utils.GetWorkingDir()
	if err != nil {
		return conf, err
	}

	jsonConf := path + "/azion/config.json"
	file, err := os.ReadFile(jsonConf)
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
