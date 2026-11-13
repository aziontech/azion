// Package deploy is `azion deploy` for API v6 accounts. It keeps the v4
// command's flags, adds the ones v6 needs, and runs the v6 pipeline in
// pkg/cmd/deploy_remote/v6.
//
// Without --local or --dry-run it runs the remote deploy (Remote): the same
// script runner script v4 uses, which runs the deploy on Azion's side. That
// lives in pkg/cmd/deploy, which imports this package, so it is given to the
// command rather than imported.
package deploy

import (
	"context"
	"time"

	"github.com/MakeNowJust/heredoc"
	msg "github.com/aziontech/azion-cli/messages/deploy"
	remotemsg "github.com/aziontech/azion-cli/messages/deploy-remote"
	deployRemote "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote"
	deployV6 "github.com/aziontech/azion-cli/pkg/cmd/deploy_remote/v6"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/poll"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"github.com/spf13/cobra"
)

// DeployCmd holds the flags of the v6 deploy.
type DeployCmd struct {
	F *cmdutil.Factory

	Path          string
	Auto          bool
	NoPrompt      bool
	SkipBuild     bool
	ProjectConf   string
	Local         bool
	DryRun        bool
	WriteBucket   bool
	Env           string
	SkipFramework bool
	Workers       int
	AliasEnv      bool
	ConfigFile    string

	Environment string
	NoActivate  bool
	Timeout     time.Duration

	// Remote runs the deploy on the script runner, as v4 does. Nil means
	// there is none, and a deploy without --local runs locally instead.
	Remote func(ctx context.Context, d *DeployCmd) error

	// Run replaces the pipeline in tests.
	run func(ctx context.Context, r *deployV6.Runner) error
}

func NewDeployCmd(f *cmdutil.Factory) *DeployCmd {
	return &DeployCmd{
		F: f, ProjectConf: "azion", Env: ".edge/.env", Timeout: poll.DefaultTimeout,
		run: func(ctx context.Context, r *deployV6.Runner) error { return r.Run(ctx) },
	}
}

func NewCobraCmd(d *DeployCmd) *cobra.Command {
	c := &cobra.Command{
		Use:           msg.DeployUsage,
		Short:         msg.DeployShortDescription,
		Long:          msg.DeployLongDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		Example: heredoc.Doc(`
       $ azion deploy --help
       $ azion deploy --environment Staging
       $ azion deploy --no-activate
       $ azion deploy --timeout 2m
       `),
		RunE: func(cmd *cobra.Command, args []string) error {
			return d.Run(cmd.Context())
		},
	}
	fl := c.Flags()
	fl.BoolP("help", "h", false, msg.DeployFlagHelp)
	fl.StringVar(&d.Path, "path", "", msg.EdgeApplicationDeployPathFlag)
	fl.BoolVar(&d.Auto, "auto", false, msg.DeployFlagAuto)
	fl.BoolVar(&d.NoPrompt, "no-prompt", false, msg.DeployFlagNoPrompt)
	fl.BoolVar(&d.SkipBuild, "skip-build", false, msg.DeployFlagSkipBuild)
	fl.StringVar(&d.ProjectConf, "config-dir", "azion", msg.EdgeApplicationDeployProjectConfFlag)
	fl.BoolVar(&d.Local, "local", false, msg.EdgeApplicationDeployLocal)
	fl.BoolVar(&d.DryRun, "dry-run", false, msg.EdgeApplicationDeployDryrun)
	fl.BoolVar(&d.WriteBucket, "writable-bucket", false, msg.WritableBucketFlag)
	fl.StringVar(&d.Env, "env", ".edge/.env", msg.EnvFlag)
	fl.BoolVar(&d.SkipFramework, "skip-framework-build", false, msg.SkipFrameworkBuild)
	fl.IntVar(&d.Workers, "workers", 0, msg.WorkersFlag)
	fl.BoolVar(&d.AliasEnv, "alias-env", false, msg.AliasEnvFlag)
	fl.StringVar(&d.ConfigFile, "config-file", "", msg.FlagConfigFile)
	// --env already names the .env file, so the Environment flag is spelled out.
	fl.StringVar(&d.Environment, "environment", "", remotemsg.V6FlagEnvironment)
	fl.BoolVar(&d.NoActivate, "no-activate", false, remotemsg.V6FlagNoActivate)
	fl.DurationVar(&d.Timeout, "timeout", poll.DefaultTimeout, remotemsg.V6FlagTimeout)
	return c
}

func NewCmd(f *cmdutil.Factory) *cobra.Command {
	return NewCobraCmd(NewDeployCmd(f))
}

// ExternalRun is the deploy init and link start when the user asks to deploy
// right away. v6 has no sync step, so sync is ignored.
func (d *DeployCmd) ExternalRun(_ *cmdutil.Factory, configPath string, _, local, skipFramework, aliasEnv bool) error {
	d.ProjectConf = configPath
	d.Local = local
	d.SkipFramework = skipFramework
	d.AliasEnv = aliasEnv
	return d.Run(context.Background())
}

func (d *DeployCmd) Run(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := vulcanPkg.CheckConfigFile(d.ConfigFile); err != nil {
		return err
	}
	if !d.Local && !d.DryRun {
		if d.Remote != nil {
			// The script runs the deploy with its defaults; it is given none
			// of the v6 pipeline's options.
			if d.Environment != "" || d.NoActivate || d.Timeout != poll.DefaultTimeout {
				return remotemsg.ErrorV6LocalOnlyFlags
			}
			return d.Remote(ctx, d)
		}
		logger.FInfoFlags(d.F.IOStreams.Out, remotemsg.V6RemoteNotAvailable, d.F.Format, d.F.Out)
	}

	shared := deployRemote.NewDeployCmd(d.F)
	shared.ProjectConf = d.ProjectConf
	shared.Env = d.Env
	shared.Path = d.Path
	shared.Auto = d.Auto
	shared.NoPrompt = d.NoPrompt
	shared.SkipBuild = d.SkipBuild
	shared.SkipFramework = d.SkipFramework
	shared.WriteBucket = d.WriteBucket
	shared.Workers = d.Workers
	shared.AliasEnv = d.AliasEnv
	shared.ConfigFile = d.ConfigFile

	runner := deployV6.NewRunner(d.F, shared, deployV6.Options{
		ConfigDir:     d.ProjectConf,
		Environment:   d.Environment,
		NoActivate:    d.NoActivate,
		Timeout:       d.Timeout,
		DryRun:        d.DryRun,
		SkipBuild:     d.SkipBuild,
		SkipFramework: d.SkipFramework,
	})
	return d.run(ctx, runner)
}
