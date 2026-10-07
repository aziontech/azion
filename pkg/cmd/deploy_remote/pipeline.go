package deploy

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/logger"
	manifestInt "github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	"github.com/aziontech/azion-cli/pkg/token"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	"go.uber.org/zap"
)

// DeployState is what the v4 deploy steps read and write.
//
// The sequence this replaces threaded all of it through locals in one ~160-line
// function. Two fields are written by one step and read by later ones —
// PathManifest and Manifest — which is the dependency the step order encodes.
type DeployState struct {
	Cmd  *DeployCmd
	F    *cmdutil.Factory
	Ctx  context.Context
	Conf *contracts.AzionApplicationOptions
	Msgs *[]string

	Clients     *Clients
	Interpreter *manifestInt.ManifestInterpreter

	// PathManifest is set by the manifest-path step and read by both reads.
	PathManifest string
	// Manifest is set by the first read and replaced by the second, after the
	// build has had a chance to regenerate it.
	Manifest *contracts.ManifestV4

	// Creds is obtained by the credentials step and used by the upload step.
	Creds token.S3Credentials

	// uploadDecided memoises shouldUpload, which both the credentials and the
	// upload step ask.
	uploadDecided bool
	upload        bool

	// OldPrefix and NewPrefix are computed once, before the pipeline runs, and
	// consumed by the build and bundler-init steps.
	OldPrefix, NewPrefix string
}

// V4DeployPipeline is the deploy flow for the v4 API.
//
// Steps do not time themselves: the runner reports each step's duration to
// HandleDeployTimingCallback, which maps the step name onto the field the
// summary prints.
type V4DeployPipeline struct{}

func (V4DeployPipeline) Name() string { return "deploy-v4" }

func (V4DeployPipeline) Steps() []pipeline.Step[DeployState] {
	return []pipeline.Step[DeployState]{
		{
			Name: "Build",
			Applies: func(s *DeployState) bool {
				return !s.Cmd.SkipBuild && s.Conf.NotFirstRun
			},
			Run: func(_ context.Context, s *DeployState) error { return s.build() },
		},
		{
			Name: "ManifestPath",
			Run: func(_ context.Context, s *DeployState) error {
				path, err := s.Interpreter.ManifestPath()
				if err != nil {
					return err
				}
				s.PathManifest = path
				return nil
			},
		},
		{
			Name: "Application",
			Run: func(_ context.Context, s *DeployState) error {
				return s.Cmd.doApplication(s.Clients.Application, context.Background(), s.Conf, s.Msgs)
			},
		},
		{
			Name: "ReadManifest",
			Run: func(_ context.Context, s *DeployState) error {
				m, err := s.Interpreter.ReadManifest(s.PathManifest, s.F, s.Msgs)
				if err != nil {
					return err
				}
				s.Manifest = m
				return nil
			},
		},
		{
			Name: "Bucket",
			Applies: func(s *DeployState) bool {
				if len(s.Manifest.Storage) == 0 {
					// Kept here so the reason for skipping still reaches the
					// debug log, as it did when this was an if/else.
					logger.Debug(msg.SkipBucket)
					return false
				}
				return true
			},
			Run: func(_ context.Context, s *DeployState) error {
				return s.Cmd.doBucket(s.Clients.Bucket, s.Ctx, s.Conf, s.Msgs, s.Manifest.Storage)
			},
		},
		{
			Name: "BundlerInit",
			Applies: func(s *DeployState) bool {
				return !s.Conf.NotFirstRun && (!s.Cmd.SkipBuild || !s.Cmd.SkipFramework)
			},
			Run: func(_ context.Context, s *DeployState) error { return s.bundlerInit() },
		},
		{
			// The manifest is read again because the build above may have
			// regenerated it. The timing accumulates onto the first read.
			Name: "ReadManifestAgain",
			Run: func(_ context.Context, s *DeployState) error {
				m, err := s.Interpreter.ReadManifest(s.PathManifest, s.F, s.Msgs)
				if err != nil {
					return err
				}
				s.Manifest = m
				return nil
			},
		},
		{
			// Obtaining the bucket credentials is its own step so that it keeps
			// its own entry in the timing summary, as it had when both lived in
			// one block.
			Name:    "Credentials",
			Applies: shouldUpload,
			Run: func(_ context.Context, s *DeployState) error {
				creds, err := s.Cmd.GetOrCreateCredentials(s.Ctx, s.Conf.Bucket, s.F.GetActiveProfile())
				if err != nil {
					return err
				}
				s.Creds = creds
				return nil
			},
		},
		{
			Name:    "UploadStaticFiles",
			Applies: shouldUpload,
			Run: func(_ context.Context, s *DeployState) error {
				for _, storage := range s.Manifest.Storage {
					err := s.Cmd.uploadFilesWithCreds(s.F, s.Conf, s.Msgs, storage.Dir, s.Conf.Bucket, s.Creds)
					if err != nil {
						return err
					}
				}
				return nil
			},
		},
		{
			Name: "RulesEngine",
			Applies: func(s *DeployState) bool {
				return len(s.Conf.RulesEngine.Rules) == 0 && !s.Conf.NotFirstRun
			},
			Run: func(_ context.Context, s *DeployState) error {
				return s.Cmd.doRulesDeploy(s.Ctx, s.Conf, s.Clients.Application, s.Msgs)
			},
		},
		{
			Name: "Manifest",
			Run: func(_ context.Context, s *DeployState) error {
				// Set before applying the manifest, exactly where the sequence
				// set it: a run that reaches this point is no longer a first run.
				s.Conf.NotFirstRun = true

				return s.Interpreter.CreateResources(s.Conf, s.Manifest, s.F, s.Cmd.ProjectConf, s.Msgs)
			},
		},
		{
			Name: "Workload",
			Applies: func(s *DeployState) bool {
				return len(s.Manifest.Workloads) == 0 || s.Manifest.Workloads[0].Name == ""
			},
			Run: func(_ context.Context, s *DeployState) error {
				return s.Cmd.doWorkload(s.Clients.Workload, s.Ctx, s.Conf, s.Msgs)
			},
		},
	}
}

// build rotates the prefix and runs the build command.
func (s *DeployState) build() error {
	cmd, f := s.Cmd, s.F
	if !cmd.SkipFramework {
		s.Conf.Prefix = s.NewPrefix
		if s.Conf.RotatePrefix == nil || *s.Conf.RotatePrefix {
			replacements := map[string]string{s.OldPrefix: s.Conf.Prefix}
			replacementsJSON, jsonErr := json.Marshal(replacements)
			if jsonErr != nil {
				return fmt.Errorf("failed to marshal replacements: %w", jsonErr)
			}
			cmdStr := fmt.Sprintf("config replace --replacements '%s'", string(replacementsJSON))
			vul := vulcanPkg.NewVulcan()
			command := vul.Command("", cmdStr, cmd.F)
			logger.Debug("Running the following command", zap.Any("Command", command))
			if err := cmd.commandRunInteractive(cmd.F, command); err != nil {
				return err
			}
		}
	}

	buildCmd := cmd.BuildCmd(f)
	err := buildCmd.ExternalRun(
		&contracts.BuildInfo{Preset: s.Conf.Preset, AliasEnv: cmd.AliasEnv},
		cmd.ProjectConf, s.Msgs, cmd.SkipFramework)
	if err != nil {
		logger.Debug("Error while running build command called by deploy command", zap.Error(err))
		return err
	}
	return nil
}

// bundlerInit initialises the bundler on a first run and builds.
func (s *DeployState) bundlerInit() error {
	cmd, f := s.Cmd, s.F
	s.Conf.Prefix = s.NewPrefix
	if err := cmd.callBundlerInit(s.Conf); err != nil {
		return err
	}
	buildCmd := cmd.BuildCmd(f)
	err := buildCmd.ExternalRun(
		&contracts.BuildInfo{AliasEnv: cmd.AliasEnv}, cmd.ProjectConf, s.Msgs, cmd.SkipFramework)
	if err != nil {
		logger.Debug("Error while running build command called by deploy command", zap.Error(err))
		return err
	}
	return nil
}

// shouldUpload reports whether this run uploads static files.
//
// Both the credentials and the upload step ask, so the answer is computed once
// and remembered. Without that, a skipped upload would log its reason twice
// where the if/else-if chain this replaces logged it once.
//
// Each skip reason keeps its own message: "no static files" and "the build was
// skipped" are different diagnoses.
func shouldUpload(s *DeployState) bool {
	if s.uploadDecided {
		return s.upload
	}
	s.uploadDecided = true

	if _, err := os.Stat(PathStatic); os.IsNotExist(err) {
		logger.Debug(msg.SkipUpload)
		s.upload = false
		return false
	}
	if s.Cmd.SkipBuild || s.Cmd.SkipFramework {
		logger.Debug(msg.SkipUploadBuild)
		s.upload = false
		return false
	}
	s.upload = true
	return true
}
