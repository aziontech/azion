package publish

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/edge_applications"
	apidom "github.com/aziontech/azion-cli/pkg/api/domains"
	"github.com/aziontech/azion-cli/pkg/logger"
	sdk "github.com/aziontech/azionapi-go-sdk/edgeapplications"
	"go.uber.org/zap"

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
	"github.com/tidwall/sjson"
	"github.com/zRedShift/mimemagic"
)

type PublishCmd struct {
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
	F                     *cmdutil.Factory
	createVersionID       func() string
}

var InstanceId int64
var Path string

var DEFAULTORIGIN [1]string = [1]string{"www.example.com"}

func NewPublishCmd(f *cmdutil.Factory) *PublishCmd {
	return &PublishCmd{
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
		F:                     f,
		createVersionID:       utils.CreateVersionID,
	}
}

func NewCobraCmd(publish *PublishCmd) *cobra.Command {
	publishCmd := &cobra.Command{
		Use:           msg.EdgeApplicationsPublishUsage,
		Short:         msg.EdgeApplicationsPublishShortDescription,
		Long:          msg.EdgeApplicationsPublishLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azioncli edge_applications publish --help
        $ azioncli edge_applications publish --path dist/static
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return publish.run(publish.F)
		},
	}
	publishCmd.Flags().BoolP("help", "h", false, msg.EdgeApplicationsPublishFlagHelp)
	publishCmd.Flags().StringVar(&Path, "path", "public", msg.EdgeApplicationPublishPathFlag)
	return publishCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewPublishCmd(f))
}

func (cmd *PublishCmd) run(f *cmdutil.Factory) error {
	logger.Debug("Running publish subcommand from edge_applications command tree")

	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Error("GetWorkDir return error", zap.Error(err))
		return err
	}

	jsonConf := path + "/azion/azion.json"
	file, err := cmd.FileReader(jsonConf)
	if err != nil {
		logger.Error("FileReader return error", zap.Error(err))
		return msg.ErrorOpeningAzionFile
	}

	typeLang := gjson.Get(string(file), "type")

	if typeLang.String() == "cdn" {
		err := publishCdn(cmd, f)
		if err != nil {
			logger.Error("publishCdn return error", zap.Error(err))
			return err
		}
		return nil
	}

	if typeLang.String() == "static" {
		err = publishStatic(cmd, f)
		if err != nil {
			logger.Error("publishStatic return error", zap.Error(err))
			return err
		}
		return nil
	}

	// Run build command
	build := cmd.BuildCmd(f)
	err = build.Run()
	if err != nil {
		logger.Error("build.Run return error", zap.Error(err))
		return err
	}

	file, err = cmd.FileReader(jsonConf)
	if err != nil {
		logger.Error("FileReader return error", zap.Error(err))
		return msg.ErrorOpeningAzionFile
	}

	conf, err := cmd.GetAzionJsonContent()
	if err != nil {
		return err
	}

	versionIDG := gjson.Get(string(file), "version-id")
	var versionID = versionIDG.String()
	if versionIDG.String() == "" {
		envPath := path + "/.edge/.env"
		fileEnv, err := cmd.FileReader(envPath)
		if err != nil {
			return msg.ErrorEnvFileVulcan
		}
		verIdSlice := strings.Split(string(fileEnv), "=")
		versionID = verIdSlice[1]
	}

	var pathStatic string = ".edge/statics"
	conf.Function.File = ".edge/worker.js"

	// legacy type - will be removed once Framework Adapter is fully substituted by Vulcan
	if typeLang.String() == "nextjs" {
		pathStatic = ".vercel/output/static"
		conf.Function.File = "./out/worker.js"
	}

	// Get total amount of files to display progress
	totalFiles := 0
	if err = cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("FilepathWalk return error", zap.Error(err))
			return err
		}
		if !info.IsDir() {
			totalFiles++
		}
		return nil
	}); err != nil {
		logger.Error("FilepathWalk return error", zap.Error(err))
		return err
	}

	clientUpload := storage.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))

	logger.FInfo(f.IOStreams.Out, msg.UploadStart)

	currentFile := 0
	if err = cmd.FilepathWalk(pathStatic, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileContent, err := cmd.Open(path)
			if err != nil {
				logger.Error("cmd.Open return error", zap.Error(err))
				return err
			}

			fileString := strings.TrimPrefix(path, pathStatic)
			mimeType, err := mimemagic.MatchFilePath(path, -1)
			if err != nil {
				logger.Error("MatchFilePath return error", zap.Error(err))
				return err
			}

			if err = clientUpload.Upload(context.Background(), versionID, fileString, mimeType.MediaType(), fileContent); err != nil {
				logger.Error("clientUpload return error", zap.Error(err))
				return err
			}

			percentage := float64(currentFile+1) * 100 / float64(totalFiles)
			progress := int(percentage / 10)
			bar := strings.Repeat("#", progress) + strings.Repeat(".", 10-progress)
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf("\033[2K\r[%v] %v %v", bar, percentage, path))
			currentFile++
		}
		return nil
	}); err != nil {
		logger.Error("FilepathWalk return error", zap.Error(err))
		return err
	}

	logger.FInfo(f.IOStreams.Out, msg.UploadSuccessful)

	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	if conf.Function.Id == 0 {
		//Create New function
		PublishId, err := cmd.fillCreateRequestFromConf(client, ctx, conf)
		if err != nil {
			logger.Error("fillCreateRequestFromConf return error", zap.Error(err))
			return err
		}

		conf.Function.Id = PublishId
	} else {
		//Update existing function
		_, err := cmd.fillUpdateRequestFromConf(client, ctx, conf.Function.Id, conf)
		if err != nil {
			logger.Error("fillUpdateRequestFromConf return error", zap.Error(err))
			return err
		}
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Error("WriteAzionJsonContent return error", zap.Error(err))
		return err
	}

	cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clidom := apidom.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	applicationName := conf.Name
	if conf.Application.Name != "__DEFAULT__" {
		applicationName = conf.Application.Name
	}

	if conf.Application.Id == 0 {
		applicationId, _, err := cmd.createApplication(cliapp, ctx, conf, applicationName)
		if err != nil {
			logger.Error("createApplication return error", zap.Error(err))
			return err
		}
		conf.Application.Id = applicationId

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			logger.Error("WriteAzionJsonContent return error", zap.Error(err))
			return err
		}

		//TODO: Review what to do when user updates Function ID directly in azion.json
		err = cmd.updateRulesEngine(cliapp, ctx, conf)
		if err != nil {
			logger.Error("updateRulesEngine return error", zap.Error(err))
			return err
		}
	} else {
		err := cmd.updateApplication(cliapp, ctx, conf, applicationName)
		if err != nil {
			logger.Error("updateApplication return error", zap.Error(err))
			return err
		}
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Error("WriteAzionJsonContent return error", zap.Error(err))
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
			logger.Error("createDomain return error", zap.Error(err))
			return err
		}
		conf.Domain.Id = domain.GetId()
		newDomain = true

		//after everything was create, we now create the cache and rules required
		reqOrigin := apiapp.CreateOriginsRequest{}
		var addresses []string
		if len(conf.Origin.Address) > 0 {
			address := prepareAddresses(conf.Origin.Address)
			addresses = conf.Origin.Address
			reqOrigin.SetAddresses(address)
		} else {
			addresses := prepareAddresses(DEFAULTORIGIN[:])
			reqOrigin.SetAddresses(addresses)
		}
		reqOrigin.SetName(conf.Name)
		reqOrigin.SetHostHeader("${host}")
		origin, err := cliapp.CreateOrigins(ctx, conf.Application.Id, &reqOrigin)
		if err != nil {
			logger.Error("CreateOrigin return error", zap.Error(err))
			return err
		}
		conf.Origin.Id = origin.GetOriginId()
		conf.Origin.Address = addresses
		conf.Origin.Name = origin.GetName()
		reqCache := apiapp.CreateCacheSettingsRequest{}
		reqCache.SetName(conf.Name)
		cache, err := cliapp.CreateCacheSettingsNextApplication(ctx, &reqCache, conf.Application.Id)
		if err != nil {
			logger.Error("CreateCacheSettingsNextApplication return error", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsCacheSettingsSuccessful)
		err = cliapp.CreateRulesEngineNextApplication(ctx, conf.Application.Id, cache.GetId(), typeLang.String())
		if err != nil {
			logger.Error("CreateRulesEngineNextApplication return error", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsRulesEngineSuccessful)

	} else {
		domain, err = cmd.updateDomain(clidom, ctx, conf, domaiName)
		if err != nil {
			logger.Error("updateDomain return error", zap.Error(err))
			return err
		}
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Error("WriteAzionJsonContent return error", zap.Error(err))
		return err
	}

	domainReturnedName := []string{domain.GetDomainName()}

	if conf.RtPurge.PurgeOnPublish && !newDomain {
		err = cmd.purgeDomains(f, domainReturnedName)
		if err != nil {
			logger.Error("purgeDomains return error", zap.Error(err))
			return err
		}
	}

	logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsPublishSuccessful)
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputDomainSuccess, "https://"+domainReturnedName[0]))
	logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationPublishDomainHint)
	logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsPublishPropagation)

	return nil
}

func (cmd *PublishCmd) purgeDomains(f *cmdutil.Factory, domainNames []string) error {
	ctx := context.Background()
	clipurge := apipurge.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	err := clipurge.Purge(ctx, domainNames)
	if err != nil {
		logger.Error("clipurge.Purge return error", zap.Error(err))
		return err
	}

	fmt.Fprintln(cmd.F.IOStreams.Out, msg.EdgeApplicationsPublishOutputCachePurge)
	return nil
}

func (cmd *PublishCmd) fillCreateRequestFromConf(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqCre := api.CreateRequest{}

	//Read code to upload
	code, err := cmd.FileReader(conf.Function.File)
	if err != nil {
		logger.Error("FileReader return error", zap.Error(err))
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
		logger.Error("FileReader return error", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Error("Unmarshal return error", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqCre.SetJsonArgs(args)
	response, err := client.Create(ctx, &reqCre)
	if err != nil {
		logger.Error("client.Create return error", zap.Error(err))
		return 0, fmt.Errorf(msg.ErrorCreateFunction.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeFunctionCreate, response.GetName(), response.GetId()))
	return response.GetId(), nil
}

func (cmd *PublishCmd) fillUpdateRequestFromConf(client *api.Client, ctx context.Context, idReq int64, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqUpd := api.UpdateRequest{}

	//Read code to upload
	code, err := cmd.FileReader(conf.Function.File)
	if err != nil {
		logger.Error("FileReader return error", zap.Error(err))
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
		logger.Error("FileReader return error", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Error("Unmarshal return error", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqUpd.Id = idReq
	reqUpd.SetJsonArgs(args)
	response, err := client.Update(ctx, &reqUpd)
	if err != nil {
		logger.Error("Update return error", zap.Error(err))
		return 0, fmt.Errorf(msg.ErrorUpdateFunction.Error(), err)
	}

	logger.Info(response.GetCode())
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeFunctionUpdate, response.GetName(), idReq))
	return response.GetId(), nil
}

func (cmd *PublishCmd) runPublishPreCmdLine() error {
	conf, err := getConfig(cmd)
	if err != nil {
		logger.Error("getConfig return error", zap.Error(err))
		return err
	}

	envs, err := cmd.EnvLoader(conf.PublishData.Env)
	if err != nil {
		logger.Error("EnvLoader return error", zap.Error(err))
		return msg.ErrReadEnvFile
	}

	err = runCommand(cmd, conf, envs)
	if err != nil {
		logger.Error("runCommand return error", zap.Error(err))
		return err
	}

	return nil
}

func (cmd *PublishCmd) createApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, name string) (int64, int64, error) {
	reqApp := apiapp.CreateRequest{}
	reqApp.SetName(name)
	reqApp.SetDeliveryProtocol("http,https")
	application, err := client.Create(ctx, &reqApp)
	if err != nil {
		logger.Error("Create return error", zap.Error(err))
		return 0, 0, fmt.Errorf(msg.ErrorCreateApplication.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeApplicationCreate, application.GetName(), application.GetId()))
	reqUpApp := apiapp.UpdateRequest{}
	reqUpApp.SetEdgeFunctions(true)
	reqUpApp.SetApplicationAcceleration(true)
	reqUpApp.Id = application.GetId()
	application, err = client.Update(ctx, &reqUpApp)
	if err != nil {
		logger.Error("Update return error", zap.Error(err))
		return 0, 0, fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	reqIns := apiapp.CreateInstanceRequest{}
	reqIns.SetEdgeFunctionId(conf.Function.Id)
	reqIns.SetName(conf.Name)
	reqIns.ApplicationId = application.GetId()
	instance, err := client.CreateInstancePublish(ctx, &reqIns)
	if err != nil {
		logger.Error("CreateInstancePublish return error", zap.Error(err))
		return 0, 0, fmt.Errorf(msg.ErrorCreateInstance.Error(), err)
	}
	InstanceId = instance.GetId()
	return application.GetId(), instance.GetId(), nil
}

func (cmd *PublishCmd) createApplicationCdn(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationCdn, name string) (int64, error) {
	reqApp := apiapp.CreateRequest{}
	reqApp.SetName(name)
	reqApp.SetDeliveryProtocol("http,https")
	application, err := client.Create(ctx, &reqApp)
	if err != nil {
		logger.Error("Create return error", zap.Error(err))
		return 0, fmt.Errorf(msg.ErrorCreateApplication.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeApplicationCreate, application.GetName(), application.GetId()))
	return application.GetId(), nil
}

func (cmd *PublishCmd) updateApplication(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, name string) error {
	reqApp := apiapp.UpdateRequest{}
	reqApp.SetName(name)
	reqApp.Id = conf.Application.Id
	application, err := client.Update(ctx, &reqApp)
	if err != nil {
		logger.Error("Update return error", zap.Error(err))
		return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeApplicationUpdate, application.GetName(), application.GetId()))
	return nil
}

func (cmd *PublishCmd) updateApplicationCdn(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationCdn, name string) error {
	reqApp := apiapp.UpdateRequest{}
	reqApp.SetName(name)
	reqApp.Id = conf.Application.Id
	application, err := client.Update(ctx, &reqApp)
	if err != nil {
		logger.Error("Update return error", zap.Error(err))
		return fmt.Errorf(msg.ErrorUpdateApplication.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeApplicationUpdate, application.GetName(), application.GetId()))
	return nil
}

func (cmd *PublishCmd) createDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, name string) (apidom.DomainResponse, error) {
	reqDom := apidom.CreateRequest{}
	reqDom.SetName(name)
	reqDom.SetCnames([]string{})
	reqDom.SetCnameAccessOnly(false)
	reqDom.SetIsActive(true)
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	domain, err := client.Create(ctx, &reqDom)
	if err != nil {
		logger.Error("Create return error", zap.Error(err))
		return nil, fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputDomainCreate, name, domain.GetId()))
	return domain, nil
}

func (cmd *PublishCmd) createDomainCdn(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationCdn, name string) (apidom.DomainResponse, error) {
	reqDom := apidom.CreateRequest{}
	reqDom.SetName(name)
	reqDom.SetCnames([]string{})
	reqDom.SetCnameAccessOnly(false)
	reqDom.SetIsActive(true)
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	domain, err := client.Create(ctx, &reqDom)
	if err != nil {
		logger.Error("Create return error", zap.Error(err))
		return nil, fmt.Errorf(msg.ErrorCreateDomain.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputDomainCreate, name, domain.GetId()))
	return domain, nil
}

func (cmd *PublishCmd) updateDomain(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationOptions, name string) (apidom.DomainResponse, error) {
	reqDom := apidom.UpdateRequest{}
	reqDom.SetName(name)
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	reqDom.Id = conf.Domain.Id
	domain, err := client.Update(ctx, &reqDom)
	if err != nil {
		logger.Error("Update return error", zap.Error(err))
		return nil, fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputDomainUpdate, name, domain.GetId()))
	return domain, nil
}

func (cmd *PublishCmd) updateDomainCdn(client *apidom.Client, ctx context.Context, conf *contracts.AzionApplicationCdn, name string) (apidom.DomainResponse, error) {
	reqDom := apidom.UpdateRequest{}
	reqDom.SetName(name)
	reqDom.SetEdgeApplicationId(conf.Application.Id)
	reqDom.Id = conf.Domain.Id
	domain, err := client.Update(ctx, &reqDom)
	if err != nil {
		logger.Error("Update return error", zap.Error(err))
		return nil, fmt.Errorf(msg.ErrorUpdateDomain.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputDomainUpdate, name, domain.GetId()))
	return domain, nil
}

func (cmd *PublishCmd) updateRulesEngine(client *apiapp.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) error {
	reqRules := apiapp.UpdateRulesEngineRequest{}
	reqRules.IdApplication = conf.Application.Id

	_, err := client.UpdateRulesEnginePublish(ctx, &reqRules, InstanceId)
	if err != nil {
		logger.Error("UpdateRulesEnginePublish return error", zap.Error(err))
		return err
	}

	return nil
}

func runCommand(cmd *PublishCmd, conf *contracts.AzionApplicationConfig, envs []string) error {
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
		logger.FInfo(cmd.Io.Out, msg.EdgeApplicationsPublishRunningCmd)
		logger.FInfo(cmd.Io.Out, fmt.Sprintf("$ %s\n", command))

		output, _, err := cmd.CommandRunner(command, envs)
		if err != nil {
			logger.FInfo(cmd.Io.Out, fmt.Sprintf("%s\n", output))
			return msg.ErrFailedToRunPublishCommand
		}

		logger.FInfo(cmd.Io.Out, fmt.Sprintf("%s\n", output))

	case "on-error":
		output, exitCode, err := cmd.CommandRunner(command, envs)
		if exitCode != 0 {
			logger.FInfo(cmd.Io.Out, fmt.Sprintf("%s\n", output))
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

func getConfig(cmd *PublishCmd) (conf *contracts.AzionApplicationConfig, err error) {
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

func publishCdn(cmd *PublishCmd, f *cmdutil.Factory) error {

	conf, err := cmd.GetAzionJsonCdn()
	if err != nil {
		logger.Error("GetAzionJsonCdn return error", zap.Error(err))
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
			logger.Error("createApplicationCdn return error", zap.Error(err))
			return err
		}
		conf.Application.Id = applicationId

	} else {
		err := cmd.updateApplicationCdn(cliapp, ctx, conf, applicationName)
		if err != nil {
			logger.Error("updateApplicationCdn return error", zap.Error(err))
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
			logger.Error("createDomainCdn return error", zap.Error(err))
			return err
		}
		conf.Domain.Id = domain.GetId()
	} else {
		_, err = cmd.updateDomainCdn(clidom, ctx, conf, domainName)
		if err != nil {
			logger.Error("updateDomainCdn return error", zap.Error(err))
			return err
		}
	}

	workingDir, err := cmd.GetWorkDir()
	if err != nil {
		logger.Error("GetWorkDir return error", zap.Error(err))
		return err
	}

	azionCdnFile := workingDir + "/azion/azion.json"

	data, err := json.MarshalIndent(conf, "", "  ")
	if err != nil {
		logger.Error("MarshalIndent return error", zap.Error(err))
		return msg.ErrorUnmarshalAzionFile
	}

	err = cmd.WriteFile(azionCdnFile, data, 0644)
	if err != nil {
		logger.Error("WriteFile return error", zap.Error(err))
		return err
	}

	logger.FInfo(cmd.Io.Out, fmt.Sprintf("%s\n", msg.EdgeApplicationsCdnPublishSuccessful))
	logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsPublishPropagation)

	return nil
}

func prepareAddresses(addrs []string) (addresses []sdk.CreateOriginsRequestAddresses) {
	var addr sdk.CreateOriginsRequestAddresses
	for _, v := range addrs {
		addr.Address = v
		addresses = append(addresses, addr)
	}
	return
}

func publishStatic(cmd *PublishCmd, f *cmdutil.Factory) error {
	path, err := cmd.GetWorkDir()
	if err != nil {
		logger.Error("GetWorkDir return error", zap.Error(err))
		return err
	}

	azionJson := path + "/azion/azion.json"
	file, err := cmd.FileReader(azionJson)
	if err != nil {
		logger.Error("FileReader return error", zap.Error(err))
		return msg.ErrorOpeningAzionFile
	}

	azJson, err := sjson.Set(string(file), "version-id", cmd.createVersionID())
	if err != nil {
		logger.Error("sjson.Set return error", zap.Error(err))
		return utils.ErrorWritingAzionJsonFile
	}

	err = cmd.WriteFile(azionJson, []byte(azJson), 0644)
	if err != nil {
		logger.Error("WriteFile return error", zap.Error(err))
		return utils.ErrorWritingAzionJsonFile
	}

	conf, err := cmd.GetAzionJsonContent()
	if err != nil {
		logger.Error("GetAzionJsonContent return error", zap.Error(err))
		return err
	}

	// upload the page static
	// Get total amount of files to display progress
	totalFiles := 0
	if err = cmd.FilepathWalk(Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalFiles++
		}
		return nil
	}); err != nil {
		logger.Error("FilepathWalk return error", zap.Error(err))
		return err
	}

	clientUpload := storage.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))

	logger.FInfo(f.IOStreams.Out, msg.UploadStart)

	versionID := conf.VersionID
	currentFile := 0
	if err = cmd.FilepathWalk(Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Error("FilepathWalk return error", zap.Error(err))
			return err
		}
		if !info.IsDir() {
			fileContent, err := cmd.Open(path)
			if err != nil {
				logger.Error("FilepathWalk return error", zap.Error(err))
				return err
			}

			fileString := strings.TrimPrefix(path, Path)
			mimeType, err := mimemagic.MatchFilePath(path, -1)
			if err != nil {
				logger.Error("MatchFilePath return error", zap.Error(err))
				return err
			}

			if err = clientUpload.Upload(context.Background(), versionID, fileString, mimeType.MediaType(), fileContent); err != nil {
				logger.Error("Upload return error", zap.Error(err))
				return err
			}

			percentage := float64(currentFile+1) * 100 / float64(totalFiles)
			progress := int(percentage / 10)
			bar := strings.Repeat("#", progress) + strings.Repeat(".", 10-progress)
			logger.FInfo(f.IOStreams.Out, fmt.Sprintf("\033[2K\r[%v] %v %v", bar, percentage, path))
			currentFile++
		}
		return nil
	}); err != nil {
		logger.Error("FilepathWalk return error", zap.Error(err))
		return err
	}

	logger.FInfo(f.IOStreams.Out, msg.UploadSuccessful)

	// create function
	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	if conf.Function.Id == 0 {
		//Create New function
		PublishId, err := cmd.CreateFunction(client, ctx, conf)
		if err != nil {
			logger.Error("CreateFunction return error", zap.Error(err))
			return err
		}
		conf.Function.Id = PublishId
	} else {
		//Update existing function
		_, err := cmd.UpdateFunction(client, ctx, conf.Function.Id, conf)
		if err != nil {
			logger.Error("UpdateFunction return error", zap.Error(err))
			return err
		}
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Error("WriteAzionJsonContent return error", zap.Error(err))
		return err
	}

	clientApplication := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clientDomain := apidom.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))

	applicationName := conf.Name
	if conf.Application.Name != "__DEFAULT__" {
		applicationName = conf.Application.Name
	}

	// create application
	if conf.Application.Id == 0 {
		applicationID, instanceID, err := cmd.createApplication(clientApplication, ctx, conf, applicationName)
		if err != nil {
			logger.Error("createApplication return error", zap.Error(err))
			return err
		}
		conf.Application.Id = applicationID
		InstanceId = instanceID

		err = cmd.WriteAzionJsonContent(conf)
		if err != nil {
			logger.Error("WriteAzionJsonContent return error", zap.Error(err))
			return err
		}

		//TODO: Review what to do when user updates Function ID directly in azion.json
		err = cmd.updateRulesEngine(clientApplication, ctx, conf)
		if err != nil {
			logger.Error("updateRulesEngine return error", zap.Error(err))
			return err
		}
	} else {
		err := cmd.updateApplication(clientApplication, ctx, conf, applicationName)
		if err != nil {
			logger.Error("updateApplication return error", zap.Error(err))
			return err
		}
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Error("WriteAzionJsonContent return error", zap.Error(err))
		return err
	}

	// create domain
	domaiName := conf.Name
	if conf.Domain.Name != "__DEFAULT__" {
		domaiName = conf.Domain.Name
	}

	var domain apidom.DomainResponse

	if conf.Domain.Id == 0 {
		domain, err = cmd.createDomain(clientDomain, ctx, conf, domaiName)
		if err != nil {
			logger.Error("createDomain return error", zap.Error(err))
			return err
		}
		conf.Domain.Id = domain.GetId()
	} else {
		domain, err = cmd.updateDomain(clientDomain, ctx, conf, domaiName)
		if err != nil {
			logger.Error("updateDomain return error", zap.Error(err))
			return err
		}
	}

	if conf.Origin.Id == 0 {
		//after everything was created, we now create the cache and rules required
		reqOrigin := apiapp.CreateOriginsRequest{}
		var addresses []string
		if len(conf.Origin.Address) > 0 {
			address := prepareAddresses(conf.Origin.Address)
			addresses = conf.Origin.Address
			reqOrigin.SetAddresses(address)
		} else {
			addresses := prepareAddresses(DEFAULTORIGIN[:])
			reqOrigin.SetAddresses(addresses)
		}
		reqOrigin.SetName(conf.Name)
		reqOrigin.SetHostHeader("${host}")
		origin, err := clientApplication.CreateOrigins(ctx, conf.Application.Id, &reqOrigin)
		if err != nil {
			logger.Error("CreateOrigins return error", zap.Error(err))
			return err
		}
		conf.Origin.Id = origin.GetOriginId()
		conf.Origin.Address = addresses
		conf.Origin.Name = origin.GetName()
		reqCache := apiapp.CreateCacheSettingsRequest{}
		reqCache.SetName(conf.Name)
		cache, err := clientApplication.CreateCacheSettingsNextApplication(ctx, &reqCache, conf.Application.Id)
		if err != nil {
			logger.Error("CreateCacheSettingsNextApplication return error", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsCacheSettingsSuccessful)
		err = clientApplication.CreateRulesEngineNextApplication(ctx, conf.Application.Id, cache.GetId(), "static")
		if err != nil {
			logger.Error("CreateRulesEngineNextApplication return error", zap.Error(err))
			return err
		}
		logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsRulesEngineSuccessful)
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Error("WriteAzionJsonContent return error", zap.Error(err))
		return err
	}

	domainReturnedName := []string{domain.GetDomainName()}

	logger.FInfo(cmd.F.IOStreams.Out, msg.EdgeApplicationsPublishSuccessful)
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputDomainSuccess, "https://"+domainReturnedName[0]))

	return nil
}

func (cmd *PublishCmd) CreateFunction(client *api.Client, ctx context.Context, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqCre := api.CreateRequest{}

	conf.Function.File = "./azion/function.js"

	jsByte, err := os.ReadFile(conf.Function.File)
	if err != nil {
		logger.Error("ReadFile return error", zap.Error(err))
		return 0, utils.ErrorReadingFile
	}

	tmpl, err := template.New("jsTemplate").Parse(string(jsByte))
	if err != nil {
		logger.Error("template.New return error", zap.Error(err))
		return 0, utils.ErrorParsingModel
	}

	data := struct {
		VersionId string
	}{
		VersionId: conf.VersionID,
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		logger.Error("tmpl.Execute return error", zap.Error(err))
		return 0, utils.ErrorExecTemplate
	}

	reqCre.SetCode(result.String())
	reqCre.SetActive(true)
	if conf.Function.Name == "__DEFAULT__" {
		reqCre.SetName(conf.Name)
	} else {
		reqCre.SetName(conf.Function.Name)
	}
	args := make(map[string]interface{})
	reqCre.SetJsonArgs(args)
	response, err := client.Create(ctx, &reqCre)
	if err != nil {
		logger.Error("client.Create return error", zap.Error(err))
		return 0, fmt.Errorf(msg.ErrorCreateFunction.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeFunctionCreate, response.GetName(), response.GetId()))
	return response.GetId(), nil
}

func (cmd *PublishCmd) UpdateFunction(client *api.Client, ctx context.Context, idReq int64, conf *contracts.AzionApplicationOptions) (int64, error) {
	reqUpd := api.UpdateRequest{}

	conf.Function.File = "./azion/function.js"

	jsByte, err := os.ReadFile(conf.Function.File)
	if err != nil {
		logger.Error("ReadFile return error", zap.Error(err))
		return 0, utils.ErrorReadingFile
	}

	tmpl, err := template.New("jsTemplate").Parse(string(jsByte))
	if err != nil {
		logger.Error("template.New return error", zap.Error(err))
		return 0, utils.ErrorParsingModel
	}

	data := struct {
		VersionId string
	}{
		VersionId: conf.VersionID,
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		logger.Error("tmpl.Execute return error", zap.Error(err))
		return 0, utils.ErrorExecTemplate
	}

	reqUpd.SetCode(result.String())
	reqUpd.SetActive(true)
	if conf.Function.Name == "__DEFAULT__" {
		reqUpd.SetName(conf.Name)
	} else {
		reqUpd.SetName(conf.Function.Name)
	}

	//Read args
	marshalledArgs, err := cmd.FileReader(conf.Function.Args)
	if err != nil {
		logger.Error("FileReader return error", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorArgsFlag, err)
	}
	args := make(map[string]interface{})
	if err := json.Unmarshal(marshalledArgs, &args); err != nil {
		logger.Error("Unmarshal return error", zap.Error(err))
		return 0, fmt.Errorf("%s: %w", msg.ErrorParseArgs, err)
	}

	reqUpd.Id = idReq
	reqUpd.SetJsonArgs(args)
	response, err := client.Update(ctx, &reqUpd)
	if err != nil {
		logger.Error("client.Update return error", zap.Error(err))
		return 0, fmt.Errorf(msg.ErrorUpdateFunction.Error(), err)
	}
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.EdgeApplicationsPublishOutputEdgeFunctionUpdate, response.GetName(), idReq))
	return response.GetId(), nil
}
