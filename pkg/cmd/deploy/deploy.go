package deploy

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	apidom "github.com/aziontech/azion-cli/pkg/api/domain"
	apiapp "github.com/aziontech/azion-cli/pkg/api/edge_applications"
	api "github.com/aziontech/azion-cli/pkg/api/edge_functions"
	apiori "github.com/aziontech/azion-cli/pkg/api/origin"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeployCmd struct {
	Io                    *iostreams.IOStreams
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent   func() (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions) error
	EnvLoader             func(path string) ([]string, error)
	BuildCmd              func(f *cmdutil.Factory) *build.BuildCmd
	Open                  func(name string) (*os.File, error)
	FilepathWalk          func(root string, fn filepath.WalkFunc) error
	F                     *cmdutil.Factory
}

var InstanceId int64
var Path string

var DEFAULTORIGIN [1]string = [1]string{"www.example.com"}

func NewDeployCmd(f *cmdutil.Factory) *DeployCmd {
	return &DeployCmd{
		Io:                    f.IOStreams,
		GetWorkDir:            utils.GetWorkingDir,
		FileReader:            os.ReadFile,
		WriteFile:             os.WriteFile,
		EnvLoader:             utils.LoadEnvVarsFromFile,
		BuildCmd:              build.NewBuildCmd,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
		Open:                  os.Open,
		FilepathWalk:          filepath.Walk,
		F:                     f,
	}
}

func NewCobraCmd(deploy *DeployCmd) *cobra.Command {
	deployCmd := &cobra.Command{
		Use:           msg.DeployUsage,
		Short:         msg.DeployShortDescription,
		Long:          msg.DeployLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
        $ azion deploy --help
        $ azion deploy --path dist/storage
        `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy.Run(deploy.F)
		},
	}
	deployCmd.Flags().BoolP("help", "h", false, msg.DeployFlagHelp)
	deployCmd.Flags().StringVar(&Path, "path", "", msg.EdgeApplicationDeployPathFlag)
	return deployCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeployCmd(f))
}

func (cmd *DeployCmd) Run(f *cmdutil.Factory) error {
	logger.Debug("Running deploy command")

	// Run build command
	build := cmd.BuildCmd(f)
	err := build.Run()
	if err != nil {
		logger.Debug("Error while running build command called by deploy command", zap.Error(err))
		return err
	}

	conf, err := cmd.GetAzionJsonContent()
	if err != nil {
		return err
	}

	var pathStatic string
	conf.Function.File = ".edge/worker.js"

	switch conf.Template {
	// legacy type - will be removed once Framework Adapter is fully substituted by Vulcan
	case "nextjs":
		pathStatic = ".vercel/output/static"
		conf.Function.File = "./out/worker.js"
	case "static":
		pathStatic = "dist"
	default:
		pathStatic = ".edge/storage"
	}

	if Path != "" {
		modified := strings.Replace(Path, "./", "", -1)
		pathStatic = modified
	}

	client := api.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	cliapp := apiapp.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	clidom := apidom.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	cliori := apiori.NewClient(f.HttpClient, f.Config.GetString("api_url"), f.Config.GetString("token"))
	ctx := context.Background()

	err = cmd.uploadFiles(f, pathStatic, conf.VersionID)
	if err != nil {
		return err
	}

	err = cmd.doFunction(client, ctx, conf)
	if err != nil {
		return err
	}
	err = cmd.doApplication(cliapp, ctx, conf)
	if err != nil {
		return err
	}
	domainName, err := cmd.doDomain(clidom, ctx, conf)
	if err != nil {
		return err
	}
	err = cmd.doOrigin(cliapp, cliori, ctx, conf)
	if err != nil {
		return err
	}

	err = cmd.WriteAzionJsonContent(conf)
	if err != nil {
		logger.Debug("Error while writing azion.json file", zap.Error(err))
		return err
	}

	logger.FInfo(cmd.F.IOStreams.Out, msg.DeploySuccessful)
	logger.FInfo(cmd.F.IOStreams.Out, fmt.Sprintf(msg.DeployOutputDomainSuccess, "https://"+domainName))
	logger.FInfo(cmd.F.IOStreams.Out, msg.DeployPropagation)

	return nil
}
