package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	apidom "github.com/aziontech/azion-cli/pkg/api/domains"

	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	apipurge "github.com/aziontech/azion-cli/pkg/api/realtime_purge"
	"github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmd/edge_applications/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

type publishCmd struct {
	Io                    *iostreams.IOStreams
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	CommandRunner         func(cmd string, envvars []string) (string, int, error)
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent   func() (*contracts.AzionApplicationOptions, error)
	GetAzionJsonCdn       func() (*contracts.AzionApplicationCdn, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions) error
	EnvLoader             func(path string) ([]string, error)
	BuildCmd              func(f *cmdutil.Factory) *build.BuildCmd
	Open                  func(name string) (*os.File, error)
	FilepathWalk          func(root string, fn filepath.WalkFunc) error
	f                     *cmdutil.Factory
}

var InstanceId int64

func NewPublishCmd(f *cmdutil.Factory) *publishCmd {
	return &publishCmd{
		Io:         f.IOStreams,
		GetWorkDir: utils.GetWorkingDir,
		FileReader: os.ReadFile,
		CommandRunner: func(cmd string, envvars []string) (string, int, error) {
			return utils.RunCommandWithOutput(envvars, cmd)
		},
		WriteFile:             os.WriteFile,
		EnvLoader:             utils.LoadEnvVarsFromFile,
		BuildCmd:              build.NewBuildCmd,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		GetAzionJsonCdn:       utils.GetAzionJsonCdn,
		Open:                  os.Open,
		FilepathWalk:          filepath.Walk,
		f:                     f,
	}
}

func NewCobraCmd(publish *publishCmd) *cobra.Command {
	publishCmd := &cobra.Command{
		Use:           msg.EdgeApplicationsPublishUsage,
		Short:         msg.EdgeApplicationsPublishShortDescription,
		Long:          msg.EdgeApplicationsPublishLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_applications publish --help
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return publish.run(publish.f)
		},
	}

	publishCmd.Flags().BoolP("help", "h", false, msg.EdgeApplicationsPublishFlagHelp)

	return publishCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewPublishCmd(f))
}

func (cmd *publishCmd) run(f *cmdutil.Factory) error {

	// Run build command
	build := cmd.BuildCmd(f)
	err := build.Run()
	if err != nil {
		return err
	}

	path, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	jsonConf := path + "/azion/azion.json"
	file, err := cmd.FileReader(jsonConf)
	if err != nil {
		return msg.ErrorOpeningAzionFile
	}

	typeLang := gjson.Get(string(file), "type")

	if typeLang.String() == "cdn" {
		err := publishCdn(cmd, f)
		if err != nil {
			return err
		}
		return nil
	}

	versionID := gjson.Get(string(file), "version-id")

	pathStatic := ".vercel/output/static"

	// Get total amount of files to display progress
	totalFiles := 0
	if err = cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalFiles++
		}
		return nil
	}); err != nil {
		return err
	}

	clientUpload := storage.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))

	currentFile := 0
	if err = cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileContent, err := cmd.Open(path)
			if err != nil {
				return err
			}

			fileString := strings.TrimPrefix(path, pathStatic)
			if err = clientUpload.Upload(context.Background(), versionID.String(), fileString, fileContent); err != nil {
				return err
			}

			percentage := float64(currentFile+1) * 100 / float64(totalFiles)
			progress := int(percentage / 10)
			bar := strings.Repeat("#", progress) + strings.Repeat(".", 10-progress)
			fmt.Fprintf(f.IOStreams.Out, "\033[2K\r[%s] %.2f%% %s ", bar, percentage, path)
			currentFile++
		}
		return nil
	}); err != nil {
		return err
	}

	fmt.Fprintf(f.IOStreams.Out, msg.UploadSuccessful)

	conf, err := cmd.GetAzionJsonContent()
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

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		return err
	}

	cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clidom := apidom.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	applicationName := conf.Name
	if conf.Application.Name != "__DEFAULT__" {
		applicationName = conf.Application.Name
	}

	if conf.Application.Id == 0 {
		applicationId, err := cmd.createApplication(cliapp, ctx, conf, applicationName)
		if err != nil {
			return err
		}
		conf.Application.Id = applicationId

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			return err
		}

		//TODO: Review what to do when user updates Function ID directly in azion.json
		err = cmd.updateRulesEngine(cliapp, ctx, conf)
		if err != nil {
			return err
		}
	} else {
		err := cmd.updateApplication(cliapp, ctx, conf, applicationName)
		if err != nil {
			return err
		}
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		return err
	}

	domaiName := conf.Name
	if conf.Domain.Name != "__DEFAULT__" {
		domaiName = conf.Domain.Name
	}

	var domain apidom.DomainResponse

	newDomain := false
	if conf.Domain.Id == 0 {
		domain, err = cmd.createDomain(clidom, ctx, conf, domaiName)
		if err != nil {
			return err
		}
		conf.Domain.Id = domain.GetId()
		newDomain = true
	} else {
		domain, err = cmd.updateDomain(clidom, ctx, conf, domaiName)
		if err != nil {
			return err
		}
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		return err
	}

	domainReturnedName := []string{domain.GetDomainName()}

	if conf.RtPurge.PurgeOnPublish && !newDomain {
		err = cmd.purgeDomains(f, domainReturnedName)
		if err != nil {
			return err
		}
	}

	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishSuccessful)
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputDomainSuccess, "https//"+domainReturnedName[0])
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishPropagation)

	return nil
}

func (cmd *publishCmd) purgeDomains(f *cmdutil.Factory, domainNames []string) error {
	ctx := context.Background()
	clipurge := apipurge.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	err := clipurge.Purge(ctx, domainNames)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputCachePurge)
	return nil
}

func (cmd *publishCmd) fillCreateRequestFromConf(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqCre := api.CreateRequest{}

	//Read code to upload
	code, err := cmd.FileReader(conf.Function.File)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
	}

	reqCre.SetCode(string(code))
	reqCre.SetActive(true)
	if conf.Function.Name == "__DEFAULT__" {
		reqCre.SetName(conf.Name)
	} else {
		reqCre.SetName(conf.Function.Name)
	}

	//Read args
	marshalledArgs, err := cmd.FileReader(conf.Function.Args)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqCre.SetJsonArgs(args)
	response, err := client.Create(ctx, &reqCre)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorCreateFunction.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputEdgeFunctionCreate, response.GetName(), response.GetId())
	return response.GetId(), nil
}

func (cmd *publishCmd) fillUpdateRequestFromConf(client *api.Client, ctx context.Context, idReq int64, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqUpd := api.UpdateRequest{}

	//Read code to upload
	code, err := cmd.FileReader(conf.Function.File)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.ErrorCodeFlag, err)
	}

	reqUpd.SetCode(string(code))
	reqUpd.SetActive(true)
	if conf.Function.Name == "__DEFAULT__" {
		reqUpd.SetName(conf.Name)
	} else {
		reqUpd.SetName(conf.Function.Name)
	}

	//Read args
	marshalledArgs, err := cmd.FileReader(conf.Function.Args)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqUpd.Id = idReq
	reqUpd.SetJsonArgs(args)
	response, err := client.Update(ctx, &reqUpd)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorUpdateFunction.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputEdgeFunctionUpdate, response.GetName(), idReq)
	return response.GetId(), nil
}

func (cmd *publishCmd) runPublishPreCmdLine() error {
	conf, err := getConfig(cmd)
	if err != nil {
		return err
	}

	envs, err := cmd.EnvLoader(conf.PublishData.Env)
	if err != nil {
		return msg.ErrReadEnvFile
	}

	err = runCommand(cmd, conf, envs)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *publishCmd) createApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, name string) (int64, error) {
	reqApp := apiapp.CreateRequest{}
	reqApp.SetName(name)
	reqApp.SetDeliveryProtocol("http,https")
	application, err := client.Create(ctx, &reqApp)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorCreateApplication.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputEdgeApplicationCreate, application.GetName(), application.GetId())
	reqUpApp := apiapp.UpdateRequest{}
	reqUpApp.SetEdgeFunctions(true)
	reqUpApp.Id = application.GetId()
	application, err = client.Update(ctx, &reqUpApp)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	reqIns := apiapp.CreateInstanceRequest{}
	reqIns.SetEdgeFunctionId(conf.Function.Id)
	reqIns.SetName(conf.Name)
	reqIns.ApplicationId = application.GetId()
	instance, err := client.CreateInstance(ctx, &reqIns)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorCreateInstance.Error(), err)
	}
	InstanceId = instance.GetId()
	return application.GetId(), nil
}

func (cmd *publishCmd) createApplicationCdn(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationCdn, name string) (int64, error) {
	reqApp := apiapp.CreateRequest{}
	reqApp.SetName(name)
	reqApp.SetDeliveryProtocol("http,https")
	application, err := client.Create(ctx, &reqApp)
	if err != nil {
		return 0, fmt.Errorf(msg.ErrorCreateApplication.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputEdgeApplicationCreate, application.GetName(), application.GetId())
	return application.GetId(), nil
}

func (cmd *publishCmd) updateApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, name string) error {
	reqApp := apiapp.UpdateRequest{}
	reqApp.SetName(name)
	reqApp.Id = conf.Application.Id
	application, err := client.Update(ctx, &reqApp)
	if err != nil {
		return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputEdgeApplicationUpdate, application.GetName(), application.GetId())
	return nil
}

func (cmd *publishCmd) updateApplicationCdn(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationCdn, name string) error {
	reqApp := apiapp.UpdateRequest{}
	reqApp.SetName(name)
	reqApp.Id = conf.Application.Id
	application, err := client.Update(ctx, &reqApp)
	if err != nil {
		return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputEdgeApplicationUpdate, application.GetName(), application.GetId())
	return nil
}

func (cmd *publishCmd) createDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, name string) (apidom.DomainResponse, error) {
	reqDom := apidom.CreateRequest{}
	reqDom.SetName(name)
	reqDom.SetCnames([]string{})
	reqDom.SetCnameAccessOnly(false)
	reqDom.SetIsActive(true)
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	domain, err := client.Create(ctx, &reqDom)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputDomainCreate, name, domain.GetId())
	return domain, nil
}

func (cmd *publishCmd) createDomainCdn(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationCdn, name string) (apidom.DomainResponse, error) {
	reqDom := apidom.CreateRequest{}
	reqDom.SetName(name)
	reqDom.SetCnames([]string{})
	reqDom.SetCnameAccessOnly(false)
	reqDom.SetIsActive(true)
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	domain, err := client.Create(ctx, &reqDom)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputDomainCreate, name, domain.GetId())
	return domain, nil
}

func (cmd *publishCmd) updateDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, name string) (apidom.DomainResponse, error) {
	reqDom := apidom.UpdateRequest{}
	reqDom.SetName(name)
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	reqDom.Id = conf.Domain.Id
	domain, err := client.Update(ctx, &reqDom)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputDomainUpdate, name, domain.GetId())
	return domain, nil
}

func (cmd *publishCmd) updateDomainCdn(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationCdn, name string) (apidom.DomainResponse, error) {
	reqDom := apidom.UpdateRequest{}
	reqDom.SetName(name)
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	reqDom.Id = conf.Domain.Id
	domain, err := client.Update(ctx, &reqDom)
	if err != nil {
		return nil, fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
	}
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishOutputDomainUpdate, name, domain.GetId())
	return domain, nil
}

func (cmd *publishCmd) updateRulesEngine(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {

	reqRules := apiapp.UpdateRulesEngineRequest{}
	reqRules.IdApplication = conf.Application.Id

	_, err := client.UpdateRulesEngine(ctx, &reqRules, InstanceId)
	if err != nil {
		return err
	}

	return nil

}

func runCommand(cmd *publishCmd, conf *contracts.AzionApplicationConfig, envs []string) error {
	var command string = conf.PublishData.Cmd
	if len(conf.PublishData.Cmd) > 0 && len(conf.PublishData.Default) > 0 {
		command += " && "
	}
	command += conf.PublishData.Default

	//if no cmd is specified, we just return nil (no error)
	if command == "" {
		return nil
	}

	switch conf.PublishData.OutputCtrl {
	case "disable":
		fmt.Fprintf(cmd.Io.Out, msg.EdgeApplicationsPublishRunningCmd)
		fmt.Fprintf(cmd.Io.Out, "$ %s\n", command)

		output, _, err := cmd.CommandRunner(command, envs)
		if err != nil {
			fmt.Fprintf(cmd.Io.Out, "%s\n", output)
			return msg.ErrFailedToRunPublishCommand
		}

		fmt.Fprintf(cmd.Io.Out, "%s\n", output)

	case "on-error":
		output, exitCode, err := cmd.CommandRunner(command, envs)
		if exitCode != 0 {
			fmt.Fprintf(cmd.Io.Out, "%s\n", output)
			return msg.ErrFailedToRunPublishCommand
		}
		if err != nil {
			return err
		}

	default:
		return msg.EdgeApplicationsOutputErr
	}

	return nil
}

func getConfig(cmd *publishCmd) (conf *contracts.AzionApplicationConfig, err error) {
	path, err := cmd.GetWorkDir()
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

func publishCdn(cmd *publishCmd, f *cmdutil.Factory) error {

	conf, err := cmd.GetAzionJsonCdn()
	if err != nil {
		return err
	}

	ctx := context.Background()
	cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clidom := apidom.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	applicationName := conf.Name
	if conf.Application.Name != "__DEFAULT__" {
		applicationName = conf.Application.Name
	}

	if conf.Application.Id == 0 {
		applicationId, err := cmd.createApplicationCdn(cliapp, ctx, conf, applicationName)
		if err != nil {
			return err
		}
		conf.Application.Id = applicationId

	} else {
		err := cmd.updateApplicationCdn(cliapp, ctx, conf, applicationName)
		if err != nil {
			return err
		}
	}

	domainName := conf.Name
	if conf.Domain.Name != "__DEFAULT__" {
		domainName = conf.Domain.Name
	}

	var domain apidom.DomainResponse

	if conf.Domain.Id == 0 {
		domain, err = cmd.createDomainCdn(clidom, ctx, conf, domainName)
		if err != nil {
			return err
		}
		conf.Domain.Id = domain.GetId()
	} else {
		_, err = cmd.updateDomainCdn(clidom, ctx, conf, domainName)
		if err != nil {
			return err
		}
	}

	workingDir, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	azionCdnFile := workingDir + "/azion/azion.json"

	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}

	err = cmd.WriteFile(azionCdnFile, data, 0644)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Io.Out, "%s\n", msg.EdgeApplicationsCdnPublishSuccessful)
	fmt.Fprintf(cmd.f.IOStreams.Out, msg.EdgeApplicationsPublishPropagation)

	return nil
}
