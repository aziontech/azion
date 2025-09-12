package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	dryrun "github.com/aziontech/azion-cli/pkg/dry_run"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/pkg/v3api/storage"
	deployRemote "github.com/aziontech/azion-cli/pkg/v3commands/deploy_remote"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v4-go-sdk-dev/storage-api"
	"github.com/briandowns/spinner"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type DeployCmd struct {
	Io                    *iostreams.IOStreams
	GetWorkDir            func() (string, error)
	FileReader            func(path string) ([]byte, error)
	WriteFile             func(filename string, data []byte, perm fs.FileMode) error
	GetAzionJsonContent   func(pathConfig string) (*contracts.AzionApplicationOptionsV3, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptionsV3, confConf string) error
	BuildCmd              func(f *cmdutil.Factory) *build.BuildCmd
	Open                  func(name string) (*os.File, error)
	FilepathWalk          func(root string, fn filepath.WalkFunc) error
	F                     *cmdutil.Factory
	Unmarshal             func(data []byte, v interface{}) error
	Interpreter           func() *manifestInt.ManifestInterpreter
	VersionID             func() string
	CallScript            func(token string, id string, secret string, prefix string, name string, confDir string, cmd *DeployCmd) (string, error)
	OpenBrowser           func(f *cmdutil.Factory, urlConsoleDeploy string, cmd *DeployCmd) error
	CaptureLogs           func(execId string, token string, cmd *DeployCmd) error
	CheckToken            func(f *cmdutil.Factory) error
	ReadSettings          func() (token.Settings, error)
	UploadFiles           func(f *cmdutil.Factory, conf *contracts.AzionApplicationOptionsV3, msgs *[]string, pathStatic, bucket string, cmd *DeployCmd, settings token.Settings) error
	OpenBrowserFunc       func(input string) error
}

var (
	Path        string
	Auto        bool
	NoPrompt    bool
	SkipBuild   bool
	ProjectConf string
	Sync        bool
	DryRun      bool
	Local       bool
	Env         string
	Logs        = contracts.Logs{}
	Result      = contracts.Results{}
	DeployURL   = "https://console.azion.com"
	ScriptID    = "17ac912d-5ce9-4806-9fa7-480779e43f58"
)

func NewDeployCmd(f *cmdutil.Factory) *DeployCmd {
	return &DeployCmd{
		Io:                    f.IOStreams,
		GetWorkDir:            utils.GetWorkingDir,
		FileReader:            os.ReadFile,
		WriteFile:             os.WriteFile,
		BuildCmd:              build.NewBuildCmd,
		GetAzionJsonContent:   utils.GetAzionJsonContentV3,
		WriteAzionJsonContent: utils.WriteAzionJsonContentV3,
		Open:                  os.Open,
		FilepathWalk:          filepath.Walk,
		Unmarshal:             json.Unmarshal,
		F:                     f,
		Interpreter:           manifestInt.NewManifestInterpreter,
		VersionID:             utils.Timestamp,
		CallScript:            callScript,
		OpenBrowser:           openBrowser,
		CaptureLogs:           captureLogs,
		UploadFiles:           uploadFiles,
		CheckToken:            checkToken,
		OpenBrowserFunc:       open.Run,
		ReadSettings:          token.ReadSettings,
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
	deployCmd.Flags().BoolVar(&DryRun, "dry-run", false, msg.EdgeApplicationDeployDryrun)
	deployCmd.Flags().BoolVar(&Local, "local", false, msg.EdgeApplicationDeployLocal)
	deployCmd.Flags().StringVar(&Env, "env", ".edge/.env", msg.EnvFlag)
	return deployCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeployCmd(f))
}

func (cmd *DeployCmd) ExternalRun(f *cmdutil.Factory, configPath string, local bool, sync bool) error {
	ProjectConf = configPath
	Local = local
	Sync = sync
	return cmd.Run(f)
}

func (cmd *DeployCmd) Run(f *cmdutil.Factory) error {

	if DryRun {
		dryStructure := dryrun.NewDryrunCmd(f)
		pathWorkingDir, err := cmd.GetWorkDir()
		if err != nil {
			return err
		}
		return dryStructure.SimulateDeploy(pathWorkingDir, ProjectConf)
	}

	if Local {
		deployLocal := deployRemote.NewDeployCmd(f)
		return deployLocal.ExternalRun(f, ProjectConf, Env, Sync, Auto, SkipBuild)
	}

	msgs := []string{}
	logger.FInfoFlags(cmd.F.IOStreams.Out, "Running deploy command\n", cmd.F.Format, cmd.F.Out)
	msgs = append(msgs, "Running deploy command")
	ctx := context.Background()

	err := cmd.CheckToken(f)
	if err != nil {
		return err
	}

	settings, err := cmd.ReadSettings()
	if err != nil {
		return err
	}

	conf, err := cmd.GetAzionJsonContent(ProjectConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	//create credentials if they are not found on settings file
	if settings.S3AccessKey == "" || settings.S3SecretKey == "" {
		nameBucket := fmt.Sprintf("%s-%s", conf.Name, cmd.VersionID())
		storageClient := storage.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
		err := storageClient.CreateBucket(ctx, storage.RequestBucket{BucketCreateRequest: sdk.BucketCreateRequest{Name: nameBucket, EdgeAccess: "read_write"}})
		if err != nil {
			return err
		}

		// Get the current time
		now := time.Now()

		// Add one year to the current time
		oneYearLater := now.AddDate(1, 0, 0)

		request := new(storage.RequestCredentials)
		request.Name = nameBucket
		request.Capabilities = []string{"listAllBucketNames", "listBuckets", "listFiles", "readFiles", "writeFiles", "deleteFiles"}
		request.Bucket = &nameBucket
		request.ExpirationDate = &oneYearLater

		creds, err := storageClient.CreateCredentials(ctx, *request)
		if err != nil {
			return err
		}
		settings.S3AccessKey = creds.Data.GetAccessKey()
		settings.S3SecretKey = creds.Data.GetSecretKey()
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

	err = cmd.UploadFiles(f, conf, &msgs, localDir, settings.S3Bucket, cmd, settings)
	if err != nil {
		return err
	}

	id, err := cmd.CallScript(settings.Token, settings.S3AccessKey, settings.S3SecretKey, conf.Prefix, settings.S3Bucket, ProjectConf, cmd)
	if err != nil {
		return err
	}

	err = cmd.OpenBrowser(f, fmt.Sprintf("%s/create/deploy/%s", DeployURL, id), cmd)
	if err != nil {
		return err
	}

	err = cmd.CaptureLogs(id, settings.Token, cmd)
	if err != nil {
		return err
	}

	conf, err = cmd.GetAzionJsonContent(ProjectConf)
	if err != nil {
		logger.Debug("Failed to get Azion JSON content", zap.Error(err))
		return err
	}

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.DeploySuccessful, f.Format, f.Out)
	msgs = append(msgs, msg.DeploySuccessful)

	msgfOutputDomainSuccess := fmt.Sprintf(msg.DeployOutputDomainSuccess, conf.Domain.Url)
	logger.FInfoFlags(cmd.F.IOStreams.Out, msgfOutputDomainSuccess, f.Format, f.Out)
	msgs = append(msgs, msgfOutputDomainSuccess)

	logger.FInfoFlags(cmd.F.IOStreams.Out, msg.DeployPropagation, f.Format, f.Out)
	msgs = append(msgs, msg.DeployPropagation)

	return nil
}

func captureLogs(execId, token string, cmd *DeployCmd) error {
	logsURL := fmt.Sprintf("%s/api/script-runner/executions/%s/logs", DeployURL, execId)
	resultsURL := fmt.Sprintf("%s/api/script-runner/executions/%s/results", DeployURL, execId)

	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Suffix = " Deploying your project..."
	s.FinalMSG = "Deployed finished executing\n"
	s.Start() // Start the spinner
	defer s.Stop()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", logsURL, bytes.NewBuffer([]byte{}))
	if err != nil {
		logger.Debug("Error creating request", zap.Error(err))
		return err
	}

	// Set headers
	req.Header.Set("accept", "application/json; version=3")
	req.Header.Set("content-type", "application/json; version=3")
	req.Header.Set("Authorization", "Token "+token)

	// Send the request
	client := &http.Client{}

	for {
		resp, err := client.Do(req)
		if err != nil {
			logger.Debug("Error sending request", zap.Error(err))
			return err
		}
		defer resp.Body.Close()

		// Read the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		if err := cmd.Unmarshal(body, &Logs); err != nil {
			logger.Debug("Error unmarshalling response", zap.Error(err))
			return err
		}

		switch Logs.Status {
		case "queued", "running", "started", "pending finish":
			time.Sleep(7 * time.Second)
			continue
		case "succeeded":
			// Create a new HTTP request
			requestResults, err := http.NewRequest("GET", resultsURL, bytes.NewBuffer([]byte{}))
			if err != nil {
				logger.Debug("Error creating request", zap.Error(err))
				return err
			}

			// Set headers
			requestResults.Header.Set("accept", "application/json; version=3")
			requestResults.Header.Set("content-type", "application/json; version=3")
			requestResults.Header.Set("Authorization", "Token "+token)

			// Send the request
			clientResults := &http.Client{}

			respResults, err := clientResults.Do(requestResults)
			if err != nil {
				logger.Debug("Error sending request", zap.Error(err))
				return err
			}
			defer respResults.Body.Close()

			// Read the response
			body, err := io.ReadAll(respResults.Body)
			if err != nil {
				return err
			}

			if err := cmd.Unmarshal(body, &Result); err != nil {
				logger.Debug("Error unmarshalling response", zap.Error(err))
				return err
			}

			if Result.Result.Errors != nil {
				return fmt.Errorf(msg.ERRORCAPTURELOGS, Result.Result.Errors.Stack)
			}

			err = cmd.WriteAzionJsonContent(Result.Result.Azion, ProjectConf)
			if err != nil {
				return err
			}
		default:
			s.Stop()
			return msg.ErrorDeployRemote
		}
		s.Stop()
		break
	}

	return nil
}
