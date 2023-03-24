package build

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

var buildType = "standard"

type BuildCmd struct {
	Io                 *iostreams.IOStreams
	WriteFile          func(filename string, data []byte, perm fs.FileMode) error
	CommandRunner      func(cmd string, envvars []string) (string, int, error)
	FileReader         func(path string) ([]byte, error)
	ConfigRelativePath string
	GetWorkDir         func() (string, error)
	EnvLoader          func(path string) ([]string, error)
	Stat               func(path string) (fs.FileInfo, error)
	VersionId          func(dir string) (string, error)
	f                  *cmdutil.Factory
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
		Example: heredoc.Doc(`
        $ azioncli edge_applications build
        `),
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
		CommandRunner: func(cmd string, envs []string) (string, int, error) {
			return utils.RunCommandWithOutput(envs, cmd)
		},
		ConfigRelativePath: "/azion/config.json",
		GetWorkDir:         utils.GetWorkingDir,
		EnvLoader:          utils.LoadEnvVarsFromFile,
		WriteFile:          os.WriteFile,
		Stat:               os.Stat,
		f:                  f,
		VersionId:          createVersionID,
	}
}

func NewBuildCmd(f *cmdutil.Factory) *BuildCmd {
	return newBuildCmd(f)
}

func (cmd *BuildCmd) run() error {
	path, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	err = RunBuildCmdLine(cmd, path)
	if err != nil {
		return err
	}

	return nil
}

func RunBuildCmdLine(cmd *BuildCmd, path string) error {
	var err error

	azionJson := path + "/azion/azion.json"
	file, err := cmd.FileReader(azionJson)
	if err != nil {
		return msg.ErrorOpeningAzionFile
	}

	typeLang := gjson.Get(string(file), "type")

	if typeLang.String() == "cdn" {
		fmt.Fprintf(cmd.Io.Out, "%s\n", msg.EdgeApplicationsBuildCdn)
		return nil
	}

	conf, err := getConfig(cmd)
	if err != nil {
		return err
	}

	err = checkArgsJson(cmd)
	if err != nil {
		return err
	}

	switch typeLang.String() {
	case "nextjs", "flareact":

		//pre-build version id. Used to check if there were changes to the project
		verID, err := cmd.VersionId(path)
		if err != nil {
			return err
		}

		confS := conf.BuildData.Default
		confS = strings.Replace(confS, "%s", verID, 1)
		conf.BuildData.Default = confS

		err = runCommand(cmd, conf)
		if err != nil {
			return err
		}

		azJson, err := sjson.Set(string(file), "version-id", verID)
		if err != nil {
			return utils.ErrorWritingAzionJsonFile
		}

		err = cmd.WriteFile(azionJson, []byte(azJson), 0644)
		if err != nil {
			return utils.ErrorWritingAzionJsonFile
		}

	case "javascript":
		err := runCommand(cmd, conf)
		if err != nil {
			return err
		}

	default:
		return utils.ErrorUnsupportedType
	}

	return nil
}

func (cmd *BuildCmd) Run() error {
	buildType = "publish"
	return cmd.run()
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

	return conf, nil

}

func checkArgsJson(cmd *BuildCmd) error {
	workDirPath, err := cmd.GetWorkDir()
	if err != nil {
		return utils.ErrorInternalServerError
	}

	workDirPath += "/azion/args.json"
	_, err = cmd.FileReader(workDirPath)
	if err != nil {
		if err := cmd.WriteFile(workDirPath, []byte("{}"), 0644); err != nil {
			return fmt.Errorf(utils.ErrorCreateFile.Error(), workDirPath)
		}
	}

	return nil
}

func runCommand(cmd *BuildCmd, conf *contracts.AzionApplicationConfig) error {
	var command string = conf.BuildData.Cmd
	if len(conf.BuildData.Cmd) > 0 && len(conf.BuildData.Default) > 0 {
		command += " && "
	}
	command += conf.BuildData.Default

	//if no cmd is specified, we just return nil (no error)
	if command == "" {
		return nil
	}

	fmt.Fprintf(cmd.Io.Out, msg.EdgeApplicationsBuildStart)

	switch conf.BuildData.OutputCtrl {
	case "disable":
		fmt.Fprintf(cmd.Io.Out, msg.EdgeApplicationsBuildRunningCmd)
		fmt.Fprintf(cmd.Io.Out, "$ %s\n", command)

		output, _, err := cmd.CommandRunner(command, []string{})
		if err != nil {
			fmt.Fprintf(cmd.Io.Out, "%s\n", output)
			return msg.ErrFailedToRunBuildCommand
		}

		fmt.Fprintf(cmd.Io.Out, "%s\n", output)

	case "on-error":
		output, exitCode, err := cmd.CommandRunner(command, []string{})
		if exitCode != 0 {
			fmt.Fprintf(cmd.Io.Out, "%s\n", output)
			return msg.ErrFailedToRunBuildCommand
		}
		if err != nil {
			return err
		}

	default:
		return msg.EdgeApplicationsOutputErr
	}

	fmt.Fprintf(cmd.Io.Out, msg.EdgeApplicationsBuildSuccessful)

	return nil
}

func createVersionID(dir string) (string, error) {
	t := time.Now()
	timeFormatted := t.Format("20060102150405")
	return timeFormatted, nil
}
