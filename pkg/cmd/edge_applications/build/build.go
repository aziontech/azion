package build

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

const VERSIONID_FORMAT string = "20060102150405"

type BuildCmd struct {
	Io                  *iostreams.IOStreams
	WriteFile           func(filename string, data []byte, perm fs.FileMode) error
	CommandRunner       func(cmd string, envvars []string) (string, int, error)
	CommandRunnerStream func(out io.Writer, cmd string, envvars []string) error
	FileReader          func(path string) ([]byte, error)
	ConfigRelativePath  string
	GetWorkDir          func() (string, error)
	EnvLoader           func(path string) ([]string, error)
	Stat                func(path string) (fs.FileInfo, error)
	VersionId           func(dir string) string
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
		CommandRunnerStream: func(out io.Writer, cmd string, envs []string) error {
			return utils.RunCommandStreamOutput(f.IOStreams.Out, envs, cmd)
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
	logger.Debug("Running build subcommand from edge_applications command tree")
	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
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
		logger.Debug("Error while reading azion.json file", zap.Error(err))
		return msg.ErrorOpeningAzionFile
	}

	typeLang := gjson.Get(string(file), "type")
	mode := gjson.Get(string(file), "mode")

	if typeLang.String() == "simple" {
		logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildSimple)
		return nil
	}

	err = checkArgsJson(cmd)
	if err != nil {
		return err
	}

	if typeLang.String() != "nextjs" {
		return callVulcan(cmd, typeLang.String(), mode.String())
	}

	conf, err := getConfig(cmd)
	if err != nil {
		return err
	}

	switch typeLang.String() {
	case "nextjs":

		//pre-build version id. Used to check if there were changes to the project
		verID := cmd.VersionId(path)

		confS := conf.BuildData.Default
		confS = strings.Replace(confS, "%s", verID, 1)
		conf.BuildData.Default = confS

		err = runCommand(cmd, conf)
		if err != nil {
			return err
		}

		azJson, err := sjson.Set(string(file), "version-id", verID)
		if err != nil {
			logger.Debug("Error while writing version-id to azion.json file: ", zap.Error(err))
			return utils.ErrorWritingAzionJsonFile
		}

		err = cmd.WriteFile(azionJson, []byte(azJson), 0644)
		if err != nil {
			logger.Debug("Error while writing azion.json file", zap.Error(err))
			return utils.ErrorWritingAzionJsonFile
		}

	default:
		return utils.ErrorUnsupportedType
	}

	return nil
}

func (cmd *BuildCmd) Run() error {
	return cmd.run()
}

func getConfig(cmd *BuildCmd) (conf *contracts.AzionApplicationConfig, err error) {
	path, err := utils.GetWorkingDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return conf, err
	}

	jsonConf := path + cmd.ConfigRelativePath
	file, err := cmd.FileReader(jsonConf)
	if err != nil {
		logger.Debug("Error while reading config.json file", zap.Error(err))
		return conf, msg.ErrorOpeningConfigFile
	}

	conf = &contracts.AzionApplicationConfig{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		logger.Debug("Error while unmarshling config.json file", zap.Error(err))
		return conf, msg.ErrorUnmarshalConfigFile
	}

	return conf, nil

}

func checkArgsJson(cmd *BuildCmd) error {
	workDirPath, err := cmd.GetWorkDir()
	if err != nil {
		logger.Debug("Error while getting working directory", zap.Error(err))
		return utils.ErrorInternalServerError
	}

	workDirPath += "/azion/args.json"
	_, err = cmd.FileReader(workDirPath)
	if err != nil {
		if err := cmd.WriteFile(workDirPath, []byte("{}"), 0644); err != nil {
			logger.Debug("Error while trying to create args.json file", zap.Error(err))
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

	logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildStart)

	switch conf.BuildData.OutputCtrl {
	case "disable":
		logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildRunningCmd)
		logger.FInfo(cmd.Io.Out, fmt.Sprintf("$ %s\n", command))

		err := cmd.CommandRunnerStream(cmd.Io.Out, command, []string{})
		if err != nil {
			logger.Debug("Error while running command with simultaneous output", zap.Error(err))
			return msg.ErrFailedToRunBuildCommand
		}

	case "on-error":
		output, exitCode, err := cmd.CommandRunner(command, []string{})
		if exitCode != 0 {
			logger.FInfo(cmd.Io.Out, fmt.Sprintf("%s\n", output))
			return msg.ErrFailedToRunBuildCommand
		}
		if err != nil {
			logger.Debug("Error while running command", zap.Error(err))
			return err
		}

	default:
		return msg.EdgeApplicationsOutputErr
	}

	logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsBuildSuccessful)

	return nil
}

func createVersionID(dir string) string {
	t := time.Now()
	timeFormatted := t.Format(VERSIONID_FORMAT)
	return timeFormatted
}

func callVulcan(cmd *BuildCmd, typeLang, mode string) error {
	command := "npx --yes edge-functions@1.0.0 build --preset " + typeLang + " --mode " + mode
	_, _, err := cmd.CommandRunner(command, []string{})
	if err != nil {
		return fmt.Errorf(msg.ErrorVulcanExecute.Error(), err.Error())
	}
	return nil
}
