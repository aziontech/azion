package init

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/webapp"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/storage/memory"
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

var (
	TemplateBranch = "dev"
	TemplateMajor  = "0"
)

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
	OpenFile      func(name string) (*os.File, error)
	RemoveAll     func(path string) error
	Rename        func(oldpath string, newpath string) error
	CreateTempDir func(dir string, pattern string) (string, error)
	EnvLoader     func(path string) ([]string, error)
	Stat          func(path string) (fs.FileInfo, error)
	Mkdir         func(path string, perm os.FileMode) error
	GitPlainClone func(path string, isBare bool, o *git.CloneOptions) (*git.Repository, error)
}

func NewInitCmd(f *cmdutil.Factory) *InitCmd {
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
		OpenFile:      os.Open,
		RemoveAll:     os.RemoveAll,
		Rename:        os.Rename,
		CreateTempDir: os.MkdirTemp,
		EnvLoader:     utils.LoadEnvVarsFromFile,
		Stat:          os.Stat,
		Mkdir:         os.MkdirAll,
		GitPlainClone: git.PlainClone,
	}
}

func NewCobraCmd(init *InitCmd) *cobra.Command {
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
	return NewCobraCmd(NewInitCmd(f))
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
	if err = options.Test(info.PathWorkingDir); err != nil {
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

		if err = UpdateScript(info, cmd, path); err != nil {
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

	fmt.Fprintf(cmd.Io.Out, "%s\n", msg.WebappInitSuccessful)

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

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{URL: REPO})
	if err != nil {
		return utils.ErrorFetchingTemplates
	}

	tags, err := r.Tags()
	if err != nil {
		return msg.ErrorGetAllTags
	}

	tag, err := sortTag(tags, TemplateMajor, TemplateBranch)
	if err != nil {
		return msg.ErrorIterateOverGit
	}

	_, err = cmd.GitPlainClone(dir, false, &git.CloneOptions{
		URL:           REPO,
		ReferenceName: plumbing.ReferenceName(tag),
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

func sortTag(tags storer.ReferenceIter, major, branch string) (string, error) {
	var tagCurrent int = 0
	var tagCurrentStr string
	var tagWithMajorOk int = 0
	var tagWithMajorOKStr string
	var err error
	err = tags.ForEach(func(t *plumbing.Reference) error {
		tagFormat := formatTag(string(t.Name()))
		tagFormat = checkBranch(tagFormat, branch)
		if tagFormat != "" {
			if strings.Split(tagFormat, "")[0] == major {
				var numberTag int
				numberTag, err = strconv.Atoi(tagFormat)
				if numberTag > tagWithMajorOk {
					tagWithMajorOk = numberTag
					tagWithMajorOKStr = string(t.Name())
				}
			} else {
				var numberTag int
				numberTag, err = strconv.Atoi(tagFormat)
				if numberTag > tagCurrent {
					tagCurrent = numberTag
					tagCurrentStr = string(t.Name())
				}
			}
		}
		return err
	})

	if tagWithMajorOKStr != "" {
		return tagWithMajorOKStr, err
	}

	return tagCurrentStr, err
}

// formatTag slice tag by '/' taking index 2 where the version is, transforming it into a list taking only the numbers
func formatTag(tag string) string {
	var t string
	for _, v := range strings.Split(strings.Split(tag, "/")[2], "") {
		if _, err := strconv.Atoi(v); err == nil {
			t += v
		}
	}
	return t
}

func checkBranch(num, branch string) string {
	if branch == "dev" {
		if len(num) == 4 {
			return num
		}
	} else if len(num) == 3 {
		return num
	}
	return ""
}

func (cmd *InitCmd) runInitCmdLine(info *InitInfo) error {
	var err error

	_, err = cmd.LookPath("npm")
	if err != nil {
		return msg.ErrorNpmNotInstalled
	}

	conf, err := getConfig(cmd, info.PathWorkingDir)
	if err != nil {
		return err
	}

	envs, err := cmd.EnvLoader(conf.InitData.Env)
	if err != nil {
		return msg.ErrReadEnvFile
	}

	switch info.TypeLang {
	case "javascript":
		err = InitJavascript(info, cmd, conf, envs)
		if err != nil {
			return err
		}
	case "nextjs":
		err = InitNextjs(info, cmd, conf, envs)
		if err != nil {
			return err
		}
	case "flareact":
		err = InitFlareact(info, cmd, conf, envs)
		if err != nil {
			return err
		}
	default:
		return utils.ErrorUnsupportedType
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

func InitJavascript(info *InitInfo, cmd *InitCmd, conf *contracts.AzionApplicationConfig, envs []string) error {
	pathWorker := info.PathWorkingDir + "/worker"
	if err := cmd.Mkdir(pathWorker, os.ModePerm); err != nil {
		return msg.ErrorFailedCreatingWorkerDirectory
	}

	err := runCommand(cmd, conf, envs)
	if err != nil {
		return err
	}

	return nil
}

func InitNextjs(info *InitInfo, cmd *InitCmd, conf *contracts.AzionApplicationConfig, envs []string) error {
	pathWorker := info.PathWorkingDir + "/worker"
	if err := cmd.Mkdir(pathWorker, os.ModePerm); err != nil {
		return msg.ErrorFailedCreatingWorkerDirectory
	}

	err := runCommand(cmd, conf, envs)
	if err != nil {
		return err
	}
	showInstructions(cmd)
	return nil
}

func showInstructions(cmd *InitCmd) {
	fmt.Fprintf(cmd.Io.Out, `    [ General Instructions ]
    - Requirements:
        - Tools: npm
        - AWS Credentials (./azion/webdev.env): AWS_ACCESS_KEY_ID AWS_SECRET_ACCESS_KEY
        - Customize the path to static content - AWS S3 storage (.azion/kv.json)

    [ Usage ]
    - Build Command: npm run build
    - Publish Command: npm run deploy
    [ Notes ]
        - Node 16x or higher`)
}

func InitFlareact(info *InitInfo, cmd *InitCmd, conf *contracts.AzionApplicationConfig, envs []string) error {
	pathWorker := info.PathWorkingDir + "/worker"
	if err := cmd.Mkdir(pathWorker, os.ModePerm); err != nil {
		return msg.ErrorFailedCreatingWorkerDirectory
	}

	err := runCommand(cmd, conf, envs)
	if err != nil {
		return err
	}

	if err = cmd.Mkdir(info.PathWorkingDir+"/public", os.ModePerm); err != nil {
		return msg.ErrorFailedCreatingPublicDirectory
	}

	return nil
}

func getConfig(cmd *InitCmd, path string) (conf *contracts.AzionApplicationConfig, err error) {
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

func UpdateScript(info *InitInfo, cmd *InitCmd, path string) error {
	packageJsonPath := path + "/package.json"
	packageJson, err := cmd.FileReader(packageJsonPath)
	if err != nil {
		return msg.ErrorPackageJsonNotFound
	}

	packJsonReplaceBuild, err := sjson.Set(string(packageJson), "scripts.build", "azioncli webapp build")
	if err != nil {
		return msg.FailedUpdatingScriptsBuildField
	}

	packJsonReplaceDeploy, err := sjson.Set(packJsonReplaceBuild, "scripts.deploy", "azioncli webapp publish")
	if err != nil {
		return msg.FailedUpdatingScriptsDeployField
	}

	err = cmd.WriteFile(packageJsonPath, []byte(packJsonReplaceDeploy), 0644)
	if err != nil {
		return fmt.Errorf(utils.ErrorCreateFile.Error(), packageJsonPath)
	}

	return nil
}

func runCommand(cmd *InitCmd, conf *contracts.AzionApplicationConfig, envs []string) error {

	//if no cmd is specified, we just return nil (no error)
	if conf.InitData.Cmd == "" {
		return nil
	}

	switch conf.InitData.OutputCtrl {
	case "disable":
		fmt.Fprintf(cmd.Io.Out, msg.WebappInitRunningCmd)
		fmt.Fprintf(cmd.Io.Out, "$ %s\n", conf.InitData.Cmd)

		output, _, err := cmd.CommandRunner(conf.InitData.Cmd, envs)
		if err != nil {
			fmt.Fprintf(cmd.Io.Out, "%s\n", output)
			return msg.ErrFailedToRunInitCommand
		}

		fmt.Fprintf(cmd.Io.Out, "%s\n", output)

	case "on-error":
		output, exitCode, err := cmd.CommandRunner(conf.InitData.Cmd, envs)
		if exitCode != 0 {
			fmt.Fprintf(cmd.Io.Out, "%s\n", output)
			return msg.ErrFailedToRunInitCommand
		}
		if err != nil {
			return err
		}

	default:
		return msg.WebappOutputErr
	}

	return nil
}
