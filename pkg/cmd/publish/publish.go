package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/MakeNowJust/heredoc"
	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	errmsg "github.com/aziontech/azion-cli/pkg/cmd/edge_functions/error_messages"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
)

type publishInfo struct {
	yesOption bool
	noOption  bool
}

type publishCmd struct {
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
	f             *cmdutil.Factory
}

func newPublishCmd(f *cmdutil.Factory) *publishCmd {
	return &publishCmd{
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
		f:             f,
	}
}

func newCobraCmd(publish *publishCmd) *cobra.Command {
	options := &contracts.AzionApplicationOptions{}
	info := &publishInfo{}
	cobraCmd := &cobra.Command{
		Use:           "publish [flags]",
		Short:         "Use Azion templates along with your Web applications",
		Long:          `Use Azion templates along with your Web applications`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Annotations: map[string]string{
			"Category": "Publish",
		},
		Example: heredoc.Doc(`
        $ azioncli publish
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return publish.run(publish.f, info, options)
		},
	}

	cobraCmd.Flags().BoolVarP(&info.yesOption, "yes", "y", false, "Force yes to all user input")
	cobraCmd.Flags().BoolVarP(&info.noOption, "no", "n", false, "Force no to all user input")

	return cobraCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return newCobraCmd(newPublishCmd(f))
}

func (cmd *publishCmd) run(f *cmdutil.Factory, info *publishInfo, options *contracts.AzionApplicationOptions) error {
	if info.yesOption && info.noOption {
		return ErrorYesAndNoOptions
	}

	err := cmd.runPublishPreCmdLine()
	if err != nil {
		return err
	}

	conf, err := utils.GetAzionJsonContent()
	if err != nil {
		return err
	}

	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	if conf.Function.Id == 0 {
		//Create New function
		PublishId, err := cmd.fillCreateRequestFromConf(client, ctx, conf)
		if err != nil {
			return err
		}

		conf.Function.Id = PublishId
	} else {
		//Update existing function
		_, err := cmd.fillUpdateRequestFromConf(client, ctx, conf.Function.Id, conf)
		if err != nil {
			return err
		}
	}

	cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	applicationName := conf.Name
	if conf.Application.Name != "__DEFAULT__" {
		applicationName = conf.Application.Name
	}

	if conf.Application.Id == 0 {
		reqApp := apiapp.CreateRequest{}
		reqApp.SetName(applicationName)
		reqApp.SetDeliveryProtocol("http,https")
		application, err := cliapp.Create(ctx, &reqApp)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.f.IOStreams.Out, "Created Edge Application with ID %d\n", application.GetId())
		reqUpApp := apiapp.UpdateRequest{}
		reqUpApp.SetEdgeFunctions(true)
		reqUpApp.Id = strconv.FormatInt(application.GetId(), 10)
		application, err = cliapp.Update(ctx, &reqUpApp)
		if err != nil {
			return err
		}
		reqIns := apiapp.CreateInstanceRequest{}
		reqIns.SetEdgeFunctionId(conf.Function.Id)
		reqIns.SetName(conf.Name)
		reqIns.ApplicationId = application.GetId()
		_, err = cliapp.CreateInstance(ctx, &reqIns)
		if err != nil {
			return err
		}
		conf.Application.Id = application.GetId()
		conf.Application.Name = application.GetName()
	} else {
		reqApp := apiapp.UpdateRequest{}
		reqApp.SetName(applicationName)
		reqApp.Id = strconv.FormatInt(conf.Application.Id, 10)
		application, err := cliapp.Update(ctx, &reqApp)
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.f.IOStreams.Out, "Updated Edge Application with ID %d\n", application.GetId())
		reqIns := apiapp.UpdateInstanceRequest{}
		reqIns.SetName(conf.Name)
		reqIns.SetEdgeFunctionId(conf.Function.Id)
		conf.Application.Name = application.GetName()
	}

	err = utils.WriteAzionJsonContent(conf)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *publishCmd) fillCreateRequestFromConf(client *api.Client, ctx context.Context, conf *contracts.AzionJsonData) (int64, error) {
	reqCre := api.CreateRequest{}

	//Read code to upload
	code, err := cmd.fileReader(conf.Function.File)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsg.ErrorCodeFlag, err)
	}

	reqCre.SetCode(string(code))
	reqCre.SetActive(conf.Function.Active)
	if conf.Function.Name == "__DEFAULT__" {
		reqCre.SetName(conf.Name)
	} else {
		reqCre.SetName(conf.Function.Name)
	}

	//Read args
	marshalledArgs, err := cmd.fileReader(conf.Function.Args)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		return 0, fmt.Errorf("%s: %w", errmsg.ErrorParseArgs, err)
	}

	reqCre.SetJsonArgs(args)
	response, err := client.Create(ctx, &reqCre)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsg.ErrorCreateFunction, err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, "Created Edge Function with ID %d\n", response.GetId())
	return response.GetId(), nil
}

func (cmd *publishCmd) fillUpdateRequestFromConf(client *api.Client, ctx context.Context, idReq int64, conf *contracts.AzionJsonData) (int64, error) {
	reqUpd := api.UpdateRequest{}

	//Read code to upload
	code, err := cmd.fileReader(conf.Function.File)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsg.ErrorCodeFlag, err)
	}

	reqUpd.SetCode(string(code))
	reqUpd.SetActive(conf.Function.Active)
	if conf.Function.Name == "__DEFAULT__" {
		reqUpd.SetName(conf.Name)
	} else {
		reqUpd.SetName(conf.Function.Name)
	}

	//Read args
	marshalledArgs, err := cmd.fileReader(conf.Function.Args)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		return 0, fmt.Errorf("%s: %w", errmsg.ErrorParseArgs, err)
	}

	reqUpd.Id = idReq
	reqUpd.SetJsonArgs(args)
	response, err := client.Update(ctx, &reqUpd)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", errmsg.ErrorUpdateFunction, err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, "Updated Edge Function with ID %d\n", idReq)
	return response.GetId(), nil
}

func (cmd *publishCmd) runPublishPreCmdLine() error {
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

	if conf.PublishData.Cmd == "" {
		fmt.Fprintf(cmd.io.Out, "Publish pre command not specified. No action will be taken\n")
		return nil
	}

	envs, err := cmd.envLoader(conf.PublishData.Env)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.io.Out, "Running publish pre command:\n\n")
	fmt.Fprintf(cmd.io.Out, "$ %s\n", conf.PublishData.Cmd)

	output, exitCode, err := cmd.commandRunner(conf.PublishData.Cmd, envs)

	fmt.Fprintf(cmd.io.Out, "%s\n", output)
	fmt.Fprintf(cmd.io.Out, "\nCommand exited with code %d\n", exitCode)

	if err != nil {
		return utils.ErrorRunningCommand
	}

	return nil
}
