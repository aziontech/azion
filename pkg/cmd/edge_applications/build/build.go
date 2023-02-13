package build

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/storage_api"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type BuildCmd struct {
	Io                 *iostreams.IOStreams
	WriteFile          func(filename string, data []byte, perm fs.FileMode) error
	CommandRunner      func(cmd string, envvars []string) (string, int, error)
	FileReader         func(path string) ([]byte, error)
	ConfigRelativePath string
	GetWorkDir         func() (string, error)
	EnvLoader          func(path string) ([]string, error)
	Stat               func(path string) (fs.FileInfo, error)
	GetVerId           func(cmd *BuildCmd, appID string) (string, error)
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
		GetVerId:           GetVersionID,
		f:                  f,
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
	applicationId := gjson.Get(string(file), "application.id")

	if typeLang.String() == "cdn" {
		fmt.Fprintf(cmd.Io.Out, "%s\n", msg.EdgeApplicationsBuildCdn)
		return nil
	}

	conf, err := getConfig(cmd)
	if err != nil {
		return err
	}

	envs := make([]string, 0)
	notFound := false

	_, err = cmd.Stat(path + "/azion/webdev.env")
	if err == nil {
		envs, err = cmd.EnvLoader(conf.BuildData.Env)
		if err != nil {
			return msg.ErrReadEnvFile
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if typeLang.String() == "nextjs" || typeLang.String() == "flareact" {
			envs = insertAWSCredentials(cmd)
			notFound = true
		}
	} else {
		return msg.ErrReadEnvFile
	}

	err = checkArgsJson(cmd)
	if err != nil {
		return err
	}

	switch typeLang.String() {
	case "nextjs", "flareact":

		var verId string = ""

		if applicationId.String() != "" && applicationId.String() != "0" {
			verId, err = cmd.GetVerId(cmd, applicationId.String())
			if err != nil {
				return err
			}
		}

		jsonReplaceFunc, err := sjson.Set(string(file), "version-id", verId)
		if err != nil {
			return msg.ErrorFailedUpdateAzionJson
		}

		err = cmd.WriteFile(azionJson, []byte(jsonReplaceFunc), 0644)
		if err != nil {
			return fmt.Errorf(utils.ErrorCreateFile.Error(), azionJson)
		}

		err = runCommand(cmd, conf, envs)
		if err != nil {
			return err
		}
		if notFound {
			errEnv := writeWebdevEnvFile(cmd, path, envs)
			if errEnv != nil {
				return errEnv
			}
		}
	case "javascript":
		err := runCommand(cmd, conf, envs)
		if err != nil {
			return err
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

func runCommand(cmd *BuildCmd, conf *contracts.AzionApplicationConfig, envs []string) error {
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

		output, _, err := cmd.CommandRunner(command, envs)
		if err != nil {
			fmt.Fprintf(cmd.Io.Out, "%s\n", output)
			return msg.ErrFailedToRunBuildCommand
		}

		fmt.Fprintf(cmd.Io.Out, "%s\n", output)

	case "on-error":
		output, exitCode, err := cmd.CommandRunner(command, envs)
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

func insertAWSCredentials(cmd *BuildCmd) []string {

	var access string
	var secret string
	envs := make([]string, 2)

	filled := false

	fmt.Fprintf(cmd.Io.Out, "%s \n", msg.EdgeApplicationsAWSMesaage)

	for !filled {
		fmt.Fprintf(cmd.Io.Out, "%s ", msg.EdgeApplicationsAWSAcess)
		fmt.Fscanln(cmd.Io.In, &access)
		fmt.Fprintf(cmd.Io.Out, "%s ", msg.EdgeApplicationsAWSSecret)
		fmt.Fscanln(cmd.Io.In, &secret)
		fmt.Fprintf(cmd.Io.Out, "\n")
		if len(access) > 0 && len(secret) > 0 {
			filled = true
		}
		envs[0] = "AWS_ACCESS_KEY_ID=" + access
		envs[1] = "AWS_SECRET_ACCESS_KEY=" + secret
	}

	return envs

}

func writeWebdevEnvFile(cmd *BuildCmd, path string, envs []string) error {
	var fileContent string

	for _, env := range envs {
		fileContent += env + "\n"
	}

	err := cmd.WriteFile(path+"/azion/webdev.env", []byte(fileContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func GetVersionID(cmd *BuildCmd, appID string) (string, error) {

	client := api.NewClient(cmd.f.HttpClient, cmd.f.Config.GetString("storage_url"), cmd.f.Config.GetString("token"))

	ctx := context.Background()

	verId, err := client.CreateVersionId(ctx, appID)
	if err != nil {
		return "", fmt.Errorf(msg.ErrorGetVersionId.Error(), err)
	}

	return verId.GetVersionId(), nil
}
