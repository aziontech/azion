package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/storage"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeployCmd struct {
	Io                    *iostreams.IOStreams
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent   func(pathConfig string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confConf string) error
	EnvLoader             func(path string) ([]string, error)
	BuildCmd              func(f *cmdutil.Factory) *build.BuildCmd
	Open                  func(name string) (*os.File, error)
	FilepathWalk          func(root string, fn filepath.WalkFunc) error
	F                     *cmdutil.Factory
	Unmarshal             func(data []byte, v interface{}) error
	Interpreter           func() *manifestInt.ManifestInterpreter
	VersionID             func() string
}

var (
	Path        string
	Auto        bool
	NoPrompt    bool
	SkipBuild   bool
	ProjectConf string
	Sync        bool
	Env         string
)

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
		Unmarshal:             json.Unmarshal,
		F:                     f,
		Interpreter:           manifestInt.NewManifestInterpreter,
		VersionID:             utils.Timestamp,
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
       $ azion deploy --auto
       `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy.Run(deploy.F)
		},
	}
	deployCmd.Flags().BoolP("help", "h", false, msg.DeployFlagHelp)
	deployCmd.Flags().StringVar(&Path, "path", "", msg.EdgeApplicationDeployPathFlag)
	deployCmd.Flags().BoolVar(&Auto, "auto", false, msg.DeployFlagAuto)
	deployCmd.Flags().BoolVar(&NoPrompt, "no-prompt", false, msg.DeployFlagNoPrompt)
	deployCmd.Flags().BoolVar(&SkipBuild, "skip-build", false, msg.DeployFlagSkipBuild)
	deployCmd.Flags().StringVar(&ProjectConf, "config-dir", "azion", msg.EdgeApplicationDeployProjectConfFlag)
	deployCmd.Flags().BoolVar(&Sync, "sync", false, msg.EdgeApplicationDeploySync)
	deployCmd.Flags().StringVar(&Env, "env", ".edge/.env", msg.EnvFlag)
	return deployCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeployCmd(f))
}

func (cmd *DeployCmd) ExternalRun(f *cmdutil.Factory, configPath string) error {
	ProjectConf = configPath
	return cmd.Run(f)
}

func (cmd *DeployCmd) Run(f *cmdutil.Factory) error {
	msgs := []string{}
	logger.FInfoFlags(cmd.F.IOStreams.Out, "Running deploy command\n", cmd.F.Format, cmd.F.Out)
	msgs = append(msgs, "Running deploy command")
	ctx := context.Background()

	settings, err := token.ReadSettings()
	if err != nil {
		return err
	}

	conf, err := cmd.GetAzionJsonContent(ProjectConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	//create credentials if they are not found on settings file
	if settings.S3AccessKey == "" || settings.S3SecreKey == "" {
		nameBucket := fmt.Sprintf("%s-%s", conf.Name, cmd.VersionID())
		storageClient := storage.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
		err := storageClient.CreateBucket(ctx, storage.RequestBucket{BucketCreate: sdk.BucketCreate{Name: nameBucket, EdgeAccess: sdk.READ_WRITE}})
		if err != nil {
			return err
		}

		// Get the current time
		now := time.Now()

		// Add one year to the current time
		oneYearLater := now.AddDate(1, 0, 0)

		request := new(storage.RequestCredentials)
		request.Name = &nameBucket
		request.Capabilities = []string{"listAllBucketNames", "listBuckets", "listFiles", "readFiles", "writeFiles", "deleteFiles"}
		request.Bucket = &nameBucket
		request.ExpirationDate = &oneYearLater

		creds, err := storageClient.CreateCredentials(ctx, *request)
		if err != nil {
			return err
		}
		settings.S3AccessKey = creds.Data.GetAccessKey()
		settings.S3SecreKey = creds.Data.GetSecretKey()
		settings.S3Bucket = nameBucket

		err = token.WriteSettings(settings)
		if err != nil {
			return err
		}
	}

	localDir, err := cmd.GetWorkDir()
	if err != nil {
		return err
	}

	conf.Prefix = cmd.VersionID()

	err = cmd.uploadFiles(f, conf, &msgs, localDir)
	if err != nil {
		return err
	}

	id, err := callScript("azion3c6bad99a7f30a2491e5423e227050aa72a", settings.S3AccessKey, settings.S3SecreKey, conf.Prefix, settings.S3Bucket)
	if err != nil {
		return err
	}
	fmt.Println(id)

	fmt.Println("5")

	err = openBrowser(f, fmt.Sprintf("https://stage-console.azion.com/create/deploy/%s", id))
	if err != nil {
		return err
	}

	fmt.Println("6")

	return nil
}
