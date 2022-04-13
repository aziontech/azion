package init

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
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

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	config := &contracts.AzionApplicationConfig{}
	options := &contracts.AzionApplicationOptions{}
	info := &initInfo{}
	initCmd := &cobra.Command{
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

			if info.yesOption && info.noOption {
				return ErrorYesAndNoOptions
			}

			//gets the test function (if it could not find it, it means it is currently not supported)
			testFunc, ok := testFuncByType[info.typeLang]
			if !ok {
				return utils.ErrorUnsupportedType
			}

			path, err := utils.GetWorkingDir()
			if err != nil {
				return err
			}

			info.pathWorkingDir = path

			options.Test = testFunc
			if err := options.Test(info.pathWorkingDir); err != nil {
				return err
			}

			//checks if user has GIT binary installed
			_, err = exec.LookPath(GIT)
			if err != nil {
				return utils.ErrorMissingGitBinary
			}

			var response string
			shouldFetchTemplates := true
			//checks if azion directory exists and is not empty
			if _, err := os.Stat("./azion"); !errors.Is(err, os.ErrNotExist) {
				if empty, _ := utils.IsDirEmpty("./azion"); !empty {
					if info.noOption || info.yesOption {
						shouldFetchTemplates = yesNoFlagToResponse(info)
					} else {
						fmt.Fprintf(f.IOStreams.Out, "%s: ", msgContentOverridden)
						fmt.Fscanln(f.IOStreams.In, &response)
						shouldFetchTemplates, err = utils.ResponseToBool(response)
						if err != nil {
							return err
						}
					}

					if shouldFetchTemplates {
						err = utils.CleanDirectory("./azion")
						if err != nil {
							return err
						}
					}
				}

			}

			if shouldFetchTemplates {
				if err := fetchTemplates(info); err != nil {
					return err
				}

				if err := organizeJsonFile(options, info); err != nil {
					return err
				}

				fmt.Fprintf(f.IOStreams.Out, "%s\n", msgCmdSuccess)
			}

			err = runInitCmdLine(config)
			if err != nil {
				return err
			}

			return nil
		},
	}

	initCmd.Flags().StringVar(&info.name, "name", "", "Your Web Application's name")
	_ = initCmd.MarkFlagRequired("name")
	initCmd.Flags().StringVar(&info.typeLang, "type", "", "Your Web Application's type <javascript>")
	_ = initCmd.MarkFlagRequired("type")
	initCmd.Flags().BoolVarP(&info.yesOption, "yes", "y", false, "Force yes to all user input")
	initCmd.Flags().BoolVarP(&info.noOption, "no", "n", false, "Force no to all user input")

	return initCmd

}

func fetchTemplates(info *initInfo) error {

	//create temporary directory to clone template into
	dir, err := ioutil.TempDir(info.pathWorkingDir, ".template")
	if err != nil {
		return utils.ErrorInternalServerError
	}
	defer os.RemoveAll(dir)

	command := exec.Command(GIT, CLONE, REPO, dir)
	err = command.Run()
	if err != nil {
		return utils.ErrorFetchingTemplates
	}

	azionDir := info.pathWorkingDir + "/azion"

	//move contents from temporary directory into final destination
	err = os.Rename(dir+"/webdev/"+info.typeLang, azionDir)
	if err != nil {
		return utils.ErrorMovingFiles
	}

	return nil
}

func runInitCmdLine(conf *contracts.AzionApplicationConfig) error {
	path, err := utils.GetWorkingDir()
	if err != nil {
		return err
	}
	jsonConf := path + "/azion/config.json"
	file, err := os.ReadFile(jsonConf)
	if err != nil {
		fmt.Println(jsonConf)
		return ErrorOpeningConfigFile
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		return ErrorUnmarshalConfigFile
	}

	envs, err := utils.LoadEnvVarsFromFile(conf.InitData.Env)
	if err != nil {
		return utils.ErrorRunningCommand
	}

	err = utils.RunCommand(envs, conf.InitData.Cmd)
	if err != nil {
		return utils.ErrorRunningCommand
	}

	return nil
}

func organizeJsonFile(options *contracts.AzionApplicationOptions, info *initInfo) error {
	file, err := os.ReadFile(info.pathWorkingDir + "/azion/azion.json")
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
	err = ioutil.WriteFile(info.pathWorkingDir+"/azion/azion.json", data, 0644)
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
