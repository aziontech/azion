package init

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type initInfo struct {
	name     string
	typeLang string
}

const (
	GIT   string = "git"
	CLONE string = "clone"
	REPO  string = "https://github.com/aziontech/azioncli-template.git"
)

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	options := &contracts.AzionApplicationOptions{}
	info := &initInfo{}
	initCmd := &cobra.Command{
		Use:           "init [flags]",
		Short:         "Use Azion templates along with your Web applications",
		Long:          `Use Azion templates along with your Web applications`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli init --name "thisisatest" --type javascript
        `),
		RunE: func(cmd *cobra.Command, args []string) error {

			testFunc, ok := testFuncByType[info.typeLang]
			if !ok {
				return utils.ErrorUnsupportedType
			}

			// if not javascript, we currently do nothing
			if testFunc == nil {
				return nil
			}

			options.Test = testFunc
			if err := options.Test(); err != nil {
				return err
			}

			//checks if user has GIT binary installed
			_, err := exec.LookPath(GIT)
			if err != nil {
				return utils.ErrorMissingGitBinary
			}

			var response string
			//checks if azion directory exists and is not empty
			if _, err := os.Stat("./azion"); !errors.Is(err, os.ErrNotExist) {
				if empty, _ := utils.IsDirEmpty("./azion"); !empty {
					fmt.Fprintf(f.IOStreams.Out, "%s: ", msgContentOverridden)
					fmt.Fscanln(f.IOStreams.In, &response)
					switch strings.ToLower(response) {
					case "no":
						fmt.Fprintf(f.IOStreams.Out, "%s\n", msgCmdStopped)
						return nil

					case "yes":
						break

					default:
						return utils.ErrorInvalidOption
					}
				}

				err = utils.CleanDirectory("./azion")
				if err != nil {
					return err
				}

			}

			if err := fetchTemplates(info); err != nil {
				return err
			}

			if err := organizeJsonFile(options, info); err != nil {
				return err
			}

			fmt.Fprintf(f.IOStreams.Out, "%s\n", msgCmdSuccess)

			return nil
		},
	}

	initCmd.Flags().StringVar(&info.name, "name", "", "Your JAMstack Application's name")
	_ = initCmd.MarkFlagRequired("name")
	initCmd.Flags().StringVar(&info.typeLang, "type", "", "Your JAMstack Application's type (javascript | nextjs | flareact)")
	_ = initCmd.MarkFlagRequired("type")

	return initCmd

}

func fetchTemplates(info *initInfo) error {

	//create temporary directory to clone template into
	dir, err := ioutil.TempDir("./", ".template")
	if err != nil {
		return utils.ErrorInternalServerError
	}
	defer os.RemoveAll(dir)

	command := exec.Command(GIT, CLONE, REPO, dir)
	err = command.Run()
	if err != nil {
		return utils.ErrorFetchingTemplates
	}

	//move contents from temporary directory into final destination
	err = os.Rename(dir+"/webdev/"+info.typeLang, "./azion")
	if err != nil {
		fmt.Println(err)
		return utils.ErrorMovingFiles
	}

	return nil
}

func organizeJsonFile(options *contracts.AzionApplicationOptions, info *initInfo) error {
	file, err := os.ReadFile("./azion/azion.json")
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
	err = ioutil.WriteFile("./azion/azion.json", data, 0644)
	if err != nil {
		return utils.ErrorInternalServerError
	}
	return nil
}
