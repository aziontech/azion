package init

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"github.com/tidwall/sjson"
)

type InitInfo struct {
	Name           string
	TypeLang       string
	PathWorkingDir string
	YesOption      bool
	NoOption       bool
}

const (
	REPO string = "https://github.com/aziontech/azioncli-template.git"
)

type InitCmd struct {
	Io            *iostreams.IOStreams
	GetWorkDir    func() (string, error)
	FileReader    func(path string) ([]byte, error)
	CommandRunner func(cmd string, envvars []string) (string, int, error)
	LookPath      func(bin string) (string, error)
	IsDirEmpty    func(dirpath string) (bool, error)
	CleanDir      func(dirpath string) error
	WriteFile     func(filename string, data []byte, perm fs.FileMode) error
	RemoveAll     func(path string) error
	Rename        func(oldpath string, newpath string) error
	CreateTempDir func(dir string, pattern string) (string, error)
	EnvLoader     func(path string) ([]string, error)
	Stat          func(path string) (fs.FileInfo, error)
}

func newInitCmd(f *cmdutil.Factory) *InitCmd {
	return &InitCmd{
		Io:         f.IOStreams,
		GetWorkDir: utils.GetWorkingDir,
		FileReader: os.ReadFile,
		CommandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
		LookPath:      exec.LookPath,
		IsDirEmpty:    utils.IsDirEmpty,
		CleanDir:      utils.CleanDirectory,
		WriteFile:     os.WriteFile,
		RemoveAll:     os.RemoveAll,
		Rename:        os.Rename,
		CreateTempDir: os.MkdirTemp,
		EnvLoader:     utils.LoadEnvVarsFromFile,
		Stat:          os.Stat,
	}
}

func newCobraCmd(init *InitCmd) *cobra.Command {
	options := &contracts.AzionApplicationOptions{}
	info := &InitInfo{}
	cobraCmd := &cobra.Command{
		Use:           msg.WebappInitUsage,
		Short:         msg.WebappInitShortDescription,
		Long:          msg.WebappInitLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
		$ azioncli init --help
		$ azioncli webapp init --name "thisisatest" --type javascript
		$ azioncli webapp init --name "thisisatest" --type flareact
		$ azioncli webapp init --name "thisisatest" --type nextjs
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return init.run(info, options)
		},
	}
	cobraCmd.Flags().StringVar(&info.Name, "name", "", msg.WebappInitFlagName)
	_ = cobraCmd.MarkFlagRequired("name")
	cobraCmd.Flags().StringVar(&info.TypeLang, "type", "", msg.WebappInitFlagType)
	_ = cobraCmd.MarkFlagRequired("type")
	cobraCmd.Flags().BoolVarP(&info.YesOption, "yes", "y", false, msg.WebappInitFlagYes)
	cobraCmd.Flags().BoolVarP(&info.NoOption, "no", "n", false, msg.WebappInitFlagNo)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return newCobraCmd(newInitCmd(f))
}

func (cmd *InitCmd) run(info *InitInfo, options *contracts.AzionApplicationOptions) error {
	if info.YesOption && info.NoOption {
		return msg.ErrorYesAndNoOptions
	}

	//gets the test function (if it could not find it, it means it is currently not supported)
	testFunc, ok := makeTestFuncMap(cmd.Stat)[info.TypeLang]
	if !ok {
		return utils.ErrorUnsupportedType
	}

	path, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	info.PathWorkingDir = path

	options.Test = testFunc
	if err := options.Test(info.PathWorkingDir); err != nil {
		return err
	}

	var response string
	shouldFetchTemplates := true

	if empty, _ := cmd.IsDirEmpty("./azion"); !empty {
		if info.NoOption || info.YesOption {
			shouldFetchTemplates = yesNoFlagToResponse(info)
		} else {
			fmt.Fprintf(cmd.Io.Out, "%s: ", msg.WebAppInitContentOverridden)
			fmt.Fscanln(cmd.Io.In, &response)
			shouldFetchTemplates, err = utils.ResponseToBool(response)
			if err != nil {
				return err
			}
		}

		if shouldFetchTemplates {
			err = cmd.CleanDir("./azion")
			if err != nil {
				return err
			}
		}
	}

	if shouldFetchTemplates {
		if err := cmd.fetchTemplates(info); err != nil {
			return err
		}

		if err := UpdateScript(info, cmd); err != nil {
			return err
		}

		if err := cmd.organizeJsonFile(options, info); err != nil {
			return err
		}

		fmt.Fprintf(cmd.Io.Out, "%s\n", msg.WebAppInitCmdSuccess)
	}

	err = cmd.runInitCmdLine(info)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *InitCmd) fetchTemplates(info *InitInfo) error {
	//create temporary directory to clone template into
	dir, err := cmd.CreateTempDir(info.PathWorkingDir, ".template")
	if err != nil {
		return utils.ErrorInternalServerError
	}
	defer func() {
		_ = cmd.RemoveAll(dir)
	}()

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: REPO,
	})
	if err != nil {
		return utils.ErrorFetchingTemplates
	}

	azionDir := info.PathWorkingDir + "/azion"

	//move contents from temporary directory into final destination
	err = cmd.Rename(dir+"/webdev/"+info.TypeLang, azionDir)
	if err != nil {
		return utils.ErrorMovingFiles
	}

	return nil
}

func (cmd *InitCmd) runInitCmdLine(info *InitInfo) error {
	var output string
	var exitCode int
	var err error

	switch info.TypeLang {
	case "javascript":
		output, exitCode, err = InitJavascript(info, cmd)
		if err != nil {
			return errors.New("failed initialization err: " + err.Error())
		}
	case "nextjs":
		output, exitCode, err = InitNextjs(info, cmd)
		if err != nil {
			return errors.New("failed initialization err: " + err.Error())
		}
	case "flareact":
		output, exitCode, err = InitFlareact(info, cmd)
		if err != nil {
			return errors.New("failed initialization err: " + err.Error())
		}
	default:
		output = ""
		exitCode = 0
		err = errors.New("setp invalid")
	}

	fmt.Fprintf(cmd.Io.Out, "%s\n", output)
	fmt.Fprintf(cmd.Io.Out, msg.WebappOutput, exitCode)

	if err != nil {
		return utils.ErrorRunningCommand
	}

	return nil
}

func (cmd *InitCmd) organizeJsonFile(options *contracts.AzionApplicationOptions, info *InitInfo) error {
	file, err := cmd.FileReader(info.PathWorkingDir + "/azion/azion.json")
	if err != nil {
		return msg.ErrorOpeningAzionFile
	}
	err = json.Unmarshal(file, &options)
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}
	options.Name = info.Name
	options.Type = info.TypeLang

	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}

	err = cmd.WriteFile(info.PathWorkingDir+"/azion/azion.json", data, 0644)
	if err != nil {
		return utils.ErrorInternalServerError
	}
	return nil
}

func yesNoFlagToResponse(info *InitInfo) bool {
	if info.YesOption {
		return info.YesOption
	}

	return false
}

func InitJavascript(info *InitInfo, cmd *InitCmd) (string, int, error) {
	_, err := cmd.LookPath("npm")
	if err != nil {
		return "", 0, msg.ErrorMissingNpm
	}

	conf, err := getConfig()
	if err != nil {
		return "", 0, err
	}

	envs, err := cmd.EnvLoader(conf.InitData.Env)
	if err != nil {
		return "", 0, msg.ErrReadEnvFile
	}

	pathWorker := info.PathWorkingDir + "/worker"
	if err = os.MkdirAll(pathWorker, os.ModePerm); err != nil {
		return "", 0, utils.ErrorCreateDir
	}

	fmt.Fprintf(cmd.Io.Out, msg.WebappInitRunningCmd)
	fmt.Fprintf(cmd.Io.Out, "$ %s\n", conf.InitData.Cmd)

	output, exitCode, err := cmd.CommandRunner(conf.InitData.Cmd, envs)
	fmt.Println("output: ", output)

	return output, exitCode, err
}

func InitNextjs(info *InitInfo, cmd *InitCmd) (string, int, error) {
	return "", 0, nil
}

func InitFlareact(info *InitInfo, cmd *InitCmd) (string, int, error) {
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
	if conf.InitData.Cmd == "" {
		return conf, msg.ErrorWebappInitCmdNotSpecified
	}
	return conf, nil
}

func UpdateScript(info *InitInfo, cmd *InitCmd) error {
	packageJsonPath := info.PathWorkingDir + "/package.json"
	packageJson, err := cmd.FileReader(packageJsonPath)
	if err != nil {
		return msg.ErrorPackageJsonNotFound
	}

	packJsonReplaceBuild, err := sjson.Set(string(packageJson), "scripts.build", "azioncli webapp build")
	if err != nil {
		return msg.ErrorUpdateBuildScript
	}

	packJsonReplaceDeploy, err := sjson.Set(packJsonReplaceBuild, "scripts.deploy", "azioncli webapp publish")
	if err != nil {
		return msg.ErrorUpdateDeployScript
	}

	err = cmd.WriteFile(packageJsonPath, []byte(packJsonReplaceDeploy), 0644)
	if err != nil {
		return fmt.Errorf(utils.ErrorCreateFile.Error(), packageJsonPath)
	}

	return nil
}
