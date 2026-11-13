package deploy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	msg "github.com/aziontech/azion-cli/messages/deploy"
	v6deploy "github.com/aziontech/azion-cli/pkg/cmd/deploy/v6"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/utils"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// NewV6Cmd is `azion deploy` for API v6 accounts: the v6 pipeline with --local
// or --dry-run, and otherwise the same remote deploy as v4, on the same script
// runner script.
func NewV6Cmd(f *cmdutil.Factory) *cobra.Command {
	return v6deploy.NewCobraCmd(newV6Deploy(f))
}

// newV6Deploy is the v6 deploy with the remote deploy wired in. Init and link
// get it through NewExternalDeployer.
func newV6Deploy(f *cmdutil.Factory) *v6deploy.DeployCmd {
	d := v6deploy.NewDeployCmd(f)
	d.Remote = func(_ context.Context, d *v6deploy.DeployCmd) error {
		return remoteDeployV6(f, d).Run(f)
	}
	return d
}

// remoteDeployV6 is the v4 remote deploy for a v6 project. The steps are the
// same (credentials, upload, script, logs); only azion.json differs, so it is
// read and saved in the v6 format.
func remoteDeployV6(f *cmdutil.Factory, d *v6deploy.DeployCmd) *DeployCmd {
	cmd := NewDeployCmd(f)
	cmd.ProjectConf = d.ProjectConf
	cmd.Env = d.Env
	cmd.Path = d.Path
	cmd.Auto = d.Auto
	cmd.NoPrompt = d.NoPrompt
	cmd.SkipBuild = d.SkipBuild
	cmd.SkipFramework = d.SkipFramework
	cmd.WriteBucket = d.WriteBucket
	cmd.Workers = d.Workers
	cmd.AliasEnv = d.AliasEnv
	cmd.ConfigFile = d.ConfigFile

	cmd.GetAzionJsonContent = azionJsonForRemoteV6
	cmd.SaveResult = saveResultV6
	cmd.DeployedURL = deployedURLV6
	return cmd
}

// azionJsonForRemoteV6 gives the remote deploy what it reads from azion.json:
// the project name and prefix. A v6 project's azion.json cannot be read as v4
// ("workloads" is a list), so it is read as v6. Without azion.json it fails as
// the v4 remote deploy does.
func azionJsonForRemoteV6(confPath string) (*contracts.AzionApplicationOptions, error) {
	wd, err := utils.GetWorkingDir()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(path.Join(wd, confPath, "azion.json")); err != nil {
		logger.Debug("Error reading stats of azion.json file", zap.Error(err))
		return nil, utils.ErrorOpeningAzionJsonFile
	}
	conf, err := utils.GetAzionJsonContentV6(confPath)
	if err != nil {
		return nil, err
	}
	return &contracts.AzionApplicationOptions{
		Name: conf.Name, Bucket: conf.Bucket, Preset: conf.Preset, Env: conf.Env,
		Prefix: conf.Prefix, RotatePrefix: conf.RotatePrefix, NotFirstRun: conf.NotFirstRun,
	}, nil
}

// resultV6 is the script runner's result of a v6 deploy. azion.json is kept as
// the deploy wrote it, in the v6 format.
type resultV6 struct {
	Result struct {
		Azion       json.RawMessage         `json:"azion,omitempty"`
		Errors      *contracts.ErrorDetails `json:"errors,omitempty"`
		AzionConfig *contracts.AzionConfig  `json:"azion_config,omitempty"`
	} `json:"result"`
}

// saveResultV6 stores the result of a v6 remote deploy: azion.json exactly as
// the deploy on Azion's side wrote it, and its azion.config.
func saveResultV6(body []byte, cmd *DeployCmd) error {
	var res resultV6
	if err := cmd.Unmarshal(body, &res); err != nil {
		logger.Debug("Error unmarshalling response", zap.Error(err))
		return err
	}
	if res.Result.Errors != nil {
		return fmt.Errorf(msg.ERRORCAPTURELOGS, res.Result.Errors.Stack)
	}

	if azion := bytes.TrimSpace(res.Result.Azion); len(azion) > 0 && !bytes.Equal(azion, []byte("null")) {
		var conf contracts.AzionApplicationOptionsV6
		if err := json.Unmarshal(azion, &conf); err != nil {
			logger.Debug("Error unmarshalling the azion.json of the remote deploy", zap.Error(err))
			return utils.ErrorUnmarshalAzionJsonFile
		}
		if err := utils.WriteAzionJsonContentV6(&conf, cmd.ProjectConf); err != nil {
			return err
		}
	}

	if res.Result.AzionConfig == nil {
		return nil
	}
	return cmd.WriteAzionConfig(res.Result.AzionConfig)
}

// deployedURLV6 is the URL of the first Workload in a v6 azion.json.
func deployedURLV6(cmd *DeployCmd) (string, error) {
	conf, err := utils.GetAzionJsonContentV6(cmd.ProjectConf)
	if err != nil {
		return "", err
	}
	for _, w := range conf.Workloads {
		if w.Url != "" {
			return w.Url, nil
		}
	}
	return "", nil
}
