package init

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/tidwall/sjson"
)

type initInfo struct {
	name           string
	typeLang       string
	pathWorkingDir string
	yesOption      bool
	noOption       bool
}

const (
	REPO string = "https://github.com/aziontech/azioncli-template.git"
)

type initCmd struct {
	io            *iostreams.IOStreams
	getWorkDir    func() (string, error)
	fileReader    func(path string) ([]byte, error)
	commandRunner func(cmd string, envvars []string) (string, int, error)
	lookPath      func(bin string) (string, error)
	isDirEmpty    func(dirpath string) (bool, error)
	cleanDir      func(dirpath string) error
	writeFile     func(filename string, data []byte, perm fs.FileMode) error
	removeAll     func(path string) error
	rename        func(oldpath string, newpath string) error
	createTempDir func(dir string, pattern string) (string, error)
	envLoader     func(path string) ([]string, error)
	stat          func(path string) (fs.FileInfo, error)
}

func newInitCmd(f *cmdutil.Factory) *initCmd {
	return &initCmd{
		io:         f.IOStreams,
		getWorkDir: utils.GetWorkingDir,
		fileReader: os.ReadFile,
		commandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
		lookPath:      exec.LookPath,
		isDirEmpty:    utils.IsDirEmpty,
		cleanDir:      utils.CleanDirectory,
		writeFile:     os.WriteFile,
		removeAll:     os.RemoveAll,
		rename:        os.Rename,
		createTempDir: ioutil.TempDir,
		envLoader:     utils.LoadEnvVarsFromFile,
		stat:          os.Stat,
	}
}

func newCobraCmd(init *initCmd) *cobra.Command {
	options := &contracts.AzionApplicationOptions{}
	info := &initInfo{}
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
	cobraCmd.Flags().StringVar(&info.name, "name", "", msg.WebappInitFlagName)
	_ = cobraCmd.MarkFlagRequired("name")
	cobraCmd.Flags().StringVar(&info.typeLang, "type", "", msg.WebappInitFlagType)
	_ = cobraCmd.MarkFlagRequired("type")
	cobraCmd.Flags().BoolVarP(&info.yesOption, "yes", "y", false, msg.WebappInitFlagYes)
	cobraCmd.Flags().BoolVarP(&info.noOption, "no", "n", false, msg.WebappInitFlagNo)

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return newCobraCmd(newInitCmd(f))
}

func (cmd *initCmd) run(info *initInfo, options *contracts.AzionApplicationOptions) error {
	if info.yesOption && info.noOption {
		return msg.ErrorYesAndNoOptions
	}

	//gets the test function (if it could not find it, it means it is currently not supported)
	testFunc, ok := makeTestFuncMap(cmd.stat)[info.typeLang]
	if !ok {
		return utils.ErrorUnsupportedType
	}

	path, err := cmd.getWorkDir()
	if err != nil {
		return err
	}

	info.pathWorkingDir = path

	options.Test = testFunc
	workingDir := info.pathWorkingDir
	if err := options.Test(workingDir); err != nil {
		return err
	}

	var response string
	shouldFetchTemplates := true

	if empty, _ := cmd.isDirEmpty("./azion"); !empty {
		if info.noOption || info.yesOption {
			shouldFetchTemplates = yesNoFlagToResponse(info)
		} else {
			fmt.Fprintf(cmd.io.Out, "%s: ", msg.WebAppInitContentOverridden)
			fmt.Fscanln(cmd.io.In, &response)
			shouldFetchTemplates, err = utils.ResponseToBool(response)
			if err != nil {
				return err
			}
		}

		if shouldFetchTemplates {
			err = cmd.cleanDir("./azion")
			if err != nil {
				return err
			}
		}
	}

	if shouldFetchTemplates {
		if err := cmd.fetchTemplates(info); err != nil {
			return err
		}

		if err := cmd.organizeJsonFile(options, info); err != nil {
			return err
		}

		fmt.Fprintf(cmd.io.Out, "%s\n", msg.WebAppInitCmdSuccess)
	}

	err = cmd.runInitCmdLine(workingDir)
	if err != nil {
		return err
	}
	return nil
}

func (cmd *initCmd) fetchTemplates(info *initInfo) error {
	//create temporary directory to clone template into
	dir, err := cmd.createTempDir(info.pathWorkingDir, ".template")
	if err != nil {
		return utils.ErrorInternalServerError
	}
	defer func() {
		_ = cmd.removeAll(dir)
	}()

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL: REPO,
	})
	if err != nil {
		return utils.ErrorFetchingTemplates
	}

	azionDir := info.pathWorkingDir + "/azion"

	//move contents from temporary directory into final destination
	err = cmd.rename(dir+"/webdev/"+info.typeLang, azionDir)
	if err != nil {
		return utils.ErrorMovingFiles
	}

	return nil
}

func (cmd *initCmd) runInitCmdLine(workingDir string) error {
	path, err := cmd.getWorkDir()
	if err != nil {
		return err
	}
	jsonConf := path + "/azion/config.json"
	file, err := cmd.fileReader(jsonConf)
	if err != nil {
		fmt.Println(jsonConf)
		return msg.ErrorOpeningConfigFile
	}

	conf := &contracts.AzionApplicationConfig{}
	err = json.Unmarshal(file, &conf)
	fmt.Println("err :", err)
	if err != nil {
		return msg.ErrorUnmarshalConfigFile
	}

	if conf.InitData.Cmd == "" {
		fmt.Fprintf(cmd.io.Out, msg.WebappInitCmdNotSpecified)
		return nil
	}

	envs, err := cmd.envLoader(conf.InitData.Env)
	if err != nil {
		return err
	}

	_, err = cmd.lookPath("npm")
	if err != nil {
		return errors.New("npm not found")
	}

	mkdir := fmt.Sprintf("mkdir -p %s/worker", workingDir)
	npmCleanWebPack := "npm install --yes --save-dev clean-webpack-plugin"
	npmWebPackCli := "npm install --yes --save-dev webpack-cli@4.9.2"
	cmdRunner := fmt.Sprintf("%s && %s && %s", mkdir, npmCleanWebPack, npmWebPackCli)
	fmt.Fprintf(cmd.io.Out, msg.WebappInitRunningCmd)
	fmt.Fprintf(cmd.io.Out, "$ %s\n", cmdRunner)

	output, exitCode, err := cmd.commandRunner(cmdRunner, envs)

	packageJsonPath := workingDir + "/package.json"
	packageJson, err := cmd.fileReader(packageJsonPath)
	if err != nil {
		return errors.New("failed on read file")
	}

	packJsonReplaceBuild, err := sjson.Set(string(packageJson), "scripts.build", "azioncli webapp build")
	if err != nil {
		return errors.New("failed replace scripts.build")
	}

	packJsonReplaceDeploy, err := sjson.Set(string(packJsonReplaceBuild), "scripts.deploy", "azioncli webapp publish")
	if err != nil {
		return errors.New("failed replace scripts.deploy")
	}

	cmd.writeFile(packageJsonPath, []byte(packJsonReplaceDeploy), 0644)

	fmt.Fprintf(cmd.io.Out, "%s\n", output)
	fmt.Fprintf(cmd.io.Out, msg.WebappOutput, exitCode)

	if err != nil {
		return utils.ErrorRunningCommand
	}

	return nil
}

func (cmd *initCmd) organizeJsonFile(options *contracts.AzionApplicationOptions, info *initInfo) error {
	file, err := cmd.fileReader(info.pathWorkingDir + "/azion/azion.json")
	if err != nil {
		return msg.ErrorOpeningAzionFile
	}
	err = json.Unmarshal(file, &options)
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}
	options.Name = info.name

	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}

	err = cmd.writeFile(info.pathWorkingDir+"/azion/azion.json", data, 0644)
	if err != nil {
		return utils.ErrorInternalServerError
	}
	return nil
}

func yesNoFlagToResponse(info *initInfo) bool {
	if info.yesOption {
		return info.yesOption
	}

	return false
}
