package init

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type initInfo struct {
	name           string
	typeLang       string
	pathWorkingDir string
	yesOption      bool
	noOption       bool
}

const (
	GIT   string = "git"
	CLONE string = "clone"
	REPO  string = "https://github.com/aziontech/azioncli-template.git"
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
	}
}

func newCobraCmd(init *initCmd) *cobra.Command {
	options := &contracts.AzionApplicationOptions{}
	info := &initInfo{}
	cobraCmd := &cobra.Command{
		Use:           "init [flags]",
		Short:         "Use Azion templates along with your Web applications",
		Long:          `Use Azion templates along with your Web applications`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Annotations: map[string]string{
			"Category": "Build",
		},
		Example: heredoc.Doc(`
        $ azioncli init --name "thisisatest" --type javascript
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return init.run(info, options)
		},
	}
	cobraCmd.Flags().StringVar(&info.name, "name", "", "Your Web Application's name")
	_ = cobraCmd.MarkFlagRequired("name")
	cobraCmd.Flags().StringVar(&info.typeLang, "type", "", "Your Web Application's type <javascript>")
	_ = cobraCmd.MarkFlagRequired("type")
	cobraCmd.Flags().BoolVarP(&info.yesOption, "yes", "y", false, "Force yes to all user input")
	cobraCmd.Flags().BoolVarP(&info.noOption, "no", "n", false, "Force no to all user input")

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return newCobraCmd(newInitCmd(f))
}

func (cmd *initCmd) run(info *initInfo, options *contracts.AzionApplicationOptions) error {
	if info.yesOption && info.noOption {
		return ErrorYesAndNoOptions
	}

	//gets the test function (if it could not find it, it means it is currently not supported)
	testFunc, ok := testFuncByType[info.typeLang]
	if !ok {
		return utils.ErrorUnsupportedType
	}

	path, err := cmd.getWorkDir()
	if err != nil {
		return err
	}

	info.pathWorkingDir = path

	options.Test = testFunc
	if err := options.Test(info.pathWorkingDir); err != nil {
		return err
	}

	//checks if user has GIT binary installed
	_, err = cmd.lookPath(GIT)
	if err != nil {
		return utils.ErrorMissingGitBinary
	}

	var response string
	shouldFetchTemplates := true

	if empty, _ := cmd.isDirEmpty("./azion"); !empty {
		if info.noOption || info.yesOption {
			shouldFetchTemplates = yesNoFlagToResponse(info)
		} else {
			fmt.Fprintf(cmd.io.Out, "%s: ", msgContentOverridden)
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

		fmt.Fprintf(cmd.io.Out, "%s\n", msgCmdSuccess)
	}

	err = cmd.runInitCmdLine()
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

	_, _, err = cmd.commandRunner(strings.Join([]string{GIT, CLONE, REPO, dir}, " "), nil)
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

func (cmd *initCmd) runInitCmdLine() error {
	path, err := cmd.getWorkDir()
	if err != nil {
		return err
	}
	jsonConf := path + "/azion/config.json"
	file, err := cmd.fileReader(jsonConf)
	if err != nil {
		fmt.Println(jsonConf)
		return ErrorOpeningConfigFile
	}

	conf := &contracts.AzionApplicationConfig{}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return ErrorUnmarshalConfigFile
	}

	envs, err := cmd.envLoader(conf.InitData.Env)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.io.Out, "Running init command\n\n")
	fmt.Fprintf(cmd.io.Out, "$ %s\n", conf.InitData.Cmd)

	output, exitCode, err := cmd.commandRunner(conf.InitData.Cmd, envs)

	fmt.Fprintf(cmd.io.Out, "%s\n", output)
	fmt.Fprintf(cmd.io.Out, "\nCommand exited with code %d\n", exitCode)

	if err != nil {
		return utils.ErrorRunningCommand
	}

	return nil
}

func (cmd *initCmd) organizeJsonFile(options *contracts.AzionApplicationOptions, info *initInfo) error {
	file, err := cmd.fileReader(info.pathWorkingDir + "/azion/azion.json")
	if err != nil {
		return ErrorOpeningAzionFile
	}
	err = json.Unmarshal(file, &options)
	if err != nil {
		return ErrorUnmarshalAzionFile
	}
	options.Name = info.name

	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return ErrorUnmarshalAzionFile
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
