package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	"github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	deploy "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	dryrun "github.com/aziontech/azion-cli/pkg/dry_run"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/token"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-go-sdk/storage"
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
	GetAzionJsonContent   func(pathConfig string) (*contracts.AzionApplicationOptions, error)
	WriteAzionJsonContent func(conf *contracts.AzionApplicationOptions, confConf string) error
	BuildCmd              func(f *cmdutil.Factory) *build.BuildCmd
	Open                  func(name string) (*os.File, error)
	FilepathWalk          func(root string, fn filepath.WalkFunc) error
	F                     *cmdutil.Factory
	Unmarshal             func(data []byte, v interface{}) error
	Interpreter           func() *manifestInt.ManifestInterpreter
	VersionID             func() string
	CallScript            func(token string, id string, secret string, prefix string, name, confDir string, cmd *DeployCmd) (string, error)
	OpenBrowser           func(f *cmdutil.Factory, urlConsoleDeploy string, cmd *DeployCmd) error
	CaptureLogs           func(execId string, token string, cmd *DeployCmd) error
	CheckToken            func(f *cmdutil.Factory) error
	ReadSettings          func() (token.Settings, error)
	UploadFiles           func(f *cmdutil.Factory, conf *contracts.AzionApplicationOptions, msgs *[]string, pathStatic, bucket string, cmd *DeployCmd, settings token.Settings) error
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
	AzionEdges  = []string{
		"Africa, Angola, Luanda (LAD)",
		"Asia, China, Hong Kong - (HKG)",
		"Europe, France, Paris (CDG)",
		"Europe, Germany, Frankfurt (FRA)",
		"Europe, Italy, Milan (MXP)",
		"Europe, London, United Kingdom (LYC)",
		"Europe, Spain, Madrid (MAD)",
		"Europe, Sweden, Stockholm (ARN)",
		"Latin America, Argentina, Buenos Aires 1 (AEP)",
		"Latin America, Argentina, Buenos Aires 2 (AEP)",
		"Latin America, Brazil, Aracaju 1 (AJU)",
		"Latin America, Brazil, Aracaju 2 (AJU)",
		"Latin America, Brazil, Barueri (CGH)",
		"Latin America, Brazil, Belem (BEL)",
		"Latin America, Brazil, Belo Horizonte 1 (PLU)",
		"Latin America, Brazil, Belo Horizonte 2 (PLU)",
		"Latin America, Brazil, Belo Horizonte 3 (PLU)",
		"Latin America, Brazil, Brasília (BSB)",
		"Latin America, Brazil, Campo Grande (CGR)",
		"Latin America, Brazil, Campinas (VCP)",
		"Latin America, Brazil, Cotia (CGH)",
		"Latin America, Brazil, Cuiabá 1 (CGB)",
		"Latin America, Brazil, Cuiabá 2 (CGB)",
		"Latin America, Brazil, Curitiba 1 (CWB)",
		"Latin America, Brazil, Curitiba 2 (CWB)",
		"Latin America, Brazil, Florianópolis 1 (FLN)",
		"Latin America, Brazil, Florianópolis 2 (FLN)",
		"Latin America, Brazil, Fortaleza 1 (FOR)",
		"Latin America, Brazil, Fortaleza 2 (FOR)",
		"Latin America, Brazil, Fortaleza 3 (FOR)",
		"Latin America, Brazil, Franca (FRC)",
		"Latin America, Brazil, Goiania (GYN)",
		"Latin America, Brazil, João Pessoa (JPA)",
		"Latin America, Brazil, Juazeiro do Norte (JDO)",
		"Latin America, Brazil, Linhares (VIX)",
		"Latin America, Brazil, Londrina (LDB)",
		"Latin America, Brazil, Macapá 1 (MCP)",
		"Latin America, Brazil, Macapá 2 (MCP)",
		"Latin America, Brazil, Maceió (MCZ)",
		"Latin America, Brazil, Manaus 1 (MAO)",
		"Latin America, Brazil, Manaus 2 (MAO)",
		"Latin America, Brazil, Natal (NAT)",
		"Latin America, Brazil, Osasco (CGH)",
		"Latin America, Brazil, Porto Alegre 1 (POA)",
		"Latin America, Brazil, Porto Alegre 2 (POA)",
		"Latin America, Brazil, Porto Alegre 3 (POA)",
		"Latin America, Brazil, Recife 1 (REC)",
		"Latin America, Brazil, Recife 2 (REC)",
		"Latin America, Brazil, Recife 3 (REC)",
		"Latin America, Brazil, Rio Branco (RBR)",
		"Latin America, Brazil, Rio de Janeiro 1 (SDU)",
		"Latin America, Brazil, Rio de Janeiro 2 (SDU)",
		"Latin America, Brazil, Rio de Janeiro 3 (SDU)",
		"Latin America, Brazil, Rio de Janeiro 4 (GIG)",
		"Latin America, Brazil, Rio de Janeiro 5 (GIG)",
		"Latin America, Brazil, Rio de Janeiro 6 (SDU)",
		"Latin America, Brazil, Salvador 1 (SSA)",
		"Latin America, Brazil, Salvador 2 (SSA)",
		"Latin America, Brazil, Salvador 4 (SSA)",
		"Latin America, Brazil, São Luis (SLZ)",
		"Latin America, Brazil, São Paulo 1 (CGH)",
		"Latin America, Brazil, São Paulo 2 (CGH)",
		"Latin America, Brazil, São Paulo 3 (CGH)",
		"Latin America, Brazil, São Paulo 4 (CGH)",
		"Latin America, Brazil, São Paulo 5 (CGH)",
		"Latin America, Brazil, Sorocaba 1 (SOD)",
		"Latin America, Brazil, Sorocaba 2 (SOD)",
		"Latin America, Chile, Santiago (SCL)",
		"Latin America, Brazil, Santos (SSZ)",
		"Latin America, Brazil, Vitória (VIX)",
		"Latin America, Colombia, Bogota (BOG)",
		"Latin America, Mexico, Queretaro (QRO)",
		"Latin America, Peru, Lima (LIM)",
		"Latin America, Brazil, Manaus 3 (MAO)",
		"North America, USA, Ashburn (IAD)",
		"North America, USA, Atlanta (ATL)",
		"North America, USA, Chicago (MDW)",
		"North America, USA, Dallas (DAL)",
		"North America, USA, Denver (DEN)",
		"North America, USA, Los Angeles (LAX)",
		"North America, USA, McAllen (MFE)",
		"North America, USA, Miami (MIA)",
		"North America, USA, New York (EWR)",
		"North America, USA, Orlando (MCO)",
		"North America, USA, Phoenix (PHX)",
		"North America, USA, Santa Clara (SJC)",
		"North America, USA, Seattle (SEA)",
		"Oceania, Sydney, Australia - (SYD)",
	}
)

func NewDeployCmd(f *cmdutil.Factory) *DeployCmd {
	return &DeployCmd{
		Io:                    f.IOStreams,
		GetWorkDir:            utils.GetWorkingDir,
		FileReader:            os.ReadFile,
		WriteFile:             os.WriteFile,
		BuildCmd:              build.NewBuildCmd,
		GetAzionJsonContent:   utils.GetAzionJsonContent,
		WriteAzionJsonContent: utils.WriteAzionJsonContent,
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
	deployCmd.Flags().BoolVar(&Local, "local", false, msg.EdgeApplicationDeployLocal)
	deployCmd.Flags().BoolVar(&DryRun, "dry-run", false, msg.EdgeApplicationDeployDryrun)
	deployCmd.Flags().StringVar(&Env, "env", ".edge/.env", msg.EnvFlag)
	return deployCmd
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeployCmd(f))
}

func (cmd *DeployCmd) ExternalRun(f *cmdutil.Factory, configPath string, sync, local bool) error {
	Local = local
	Sync = sync
	ProjectConf = configPath
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
		deployLocal := deploy.NewDeployCmd(f)
		err := deployLocal.ExternalRun(f, ProjectConf, Env, Sync, Auto, SkipBuild)
		if err != nil {
			return err
		}

		// Obter a configuração após o deploy local
		conf, err := cmd.GetAzionJsonContent(ProjectConf)
		if err != nil {
			logger.Debug("Failed to get Azion JSON content", zap.Error(err))
			return err
		}

		// Verificar a propagação do domínio
		err = checkDomainPropagation(conf.Domain.Url)
		if err != nil {
			logger.Debug("Error checking domain propagation", zap.Error(err))
			return err
		}

		return nil
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
	if settings.S3AccessKey == "" || settings.S3SecreKey == "" {
		nameBucket := fmt.Sprintf("%s-%s", conf.Name, cmd.VersionID())
		storageClient := storage.NewClient(f.HttpClient, f.Config.GetString("storage_url"), f.Config.GetString("token"))
		err := storageClient.CreateBucket(ctx, storage.RequestBucket{BucketCreate: sdk.BucketCreate{Name: nameBucket, EdgeAccess: sdk.READ_ONLY}})
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

	err = cmd.UploadFiles(f, conf, &msgs, localDir, settings.S3Bucket, cmd, settings)
	if err != nil {
		return err
	}

	id, err := cmd.CallScript(settings.Token, settings.S3AccessKey, settings.S3SecreKey, conf.Prefix, settings.S3Bucket, ProjectConf, cmd)
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
	if !cmd.F.Debug {
		s.Start() // Start the spinner
	}
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
	logTime := time.Now()
	lastLog := ""
	// Custom layout for parsing the timestamp
	layout := "2006-01-02 15:04:05.000" // Layout must match your timestamp format exactly

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
			for _, event := range Logs.Logs {
				if event.Content != "" && event.Content != lastLog {
					// Parse the timestamp
					parsedTimestamp, err := time.Parse(layout, event.Timestamp)
					if err != nil {
						return err
					}
					if logTime.After(parsedTimestamp) {
						continue
					}
					lastLog = event.Content
					logger.Debug(event.Content)
					logTime = parsedTimestamp
				}
			}
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
				return errors.New(Result.Result.Errors.Stack) //TODO: add mensagem que deu ruim e é para verificar se criou algo na conta
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

func checkDomainPropagation(url string) error {
	maxAttempts := 180
	remainingEdges := make([]string, len(AzionEdges))
	copy(remainingEdges, AzionEdges)

	currentIndex := 0
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Suffix = fmt.Sprintf(" Propagating to %s", remainingEdges[currentIndex])
	s.Start()

	startTime := time.Now()

	for i := 0; i < maxAttempts; i++ {
		resp, err := http.Get(url)
		if err != nil {
			logger.Debug("HTTP request error", zap.Error(err))
			time.Sleep(time.Second)
			continue
		}

		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}

		logger.Debug(
			"Checking propagation",
			zap.Int("attempt", i+1),
			zap.Int("status", resp.StatusCode),
			zap.Duration("elapsed", time.Since(startTime)),
		)

		if resp.StatusCode != 404 {
			s.Stop()

			fmt.Print("\r\033[K")

			for _, edge := range remainingEdges {
				s.Suffix = fmt.Sprintf(" Propagating to %s", edge)
				s.Start()
				time.Sleep(100 * time.Millisecond)
				s.Stop()
			}

			fmt.Print("\r\033[K")
			fmt.Println("🎉 Domain propagation completed successfully!")

			err = open.Run(url)
			if err != nil {
				return fmt.Errorf("error opening browser: %w", err)
			}
			return nil
		}

		select {
		case <-ticker.C:
			currentIndex = (currentIndex + 1) % len(remainingEdges)
			s.Suffix = fmt.Sprintf(" Propagating to %s", remainingEdges[currentIndex])
		}
	}

	s.Stop()
	elapsedMinutes := time.Since(startTime).Minutes()
	return fmt.Errorf("timeout waiting for domain propagation after %.1f minutes", elapsedMinutes)
}
