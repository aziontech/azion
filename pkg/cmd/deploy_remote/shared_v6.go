package deploy

// ============================================================================
// Shared with the v6 deploy (pkg/cmd/deploy_remote/v6). Building the project,
// the bucket and the static-file upload are the same in v4 and v6 — storage is
// not versioned in v6 — so v6 reuses them through this file instead of copying
// them. They speak the v4 azion.json contract internally; Project carries the
// handful of fields they read and write, and the v6 side copies them back into
// its own contract.
// ============================================================================

import (
	"context"

	"github.com/aziontech/azion-cli/pkg/api/storage"
	"github.com/aziontech/azion-cli/pkg/cmd/build"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/registry"
	"github.com/aziontech/azion-cli/pkg/token"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
)

// Project is what the shared steps read and write about the project.
type Project struct {
	Name         string
	Bucket       string
	Prefix       string
	Preset       string
	RotatePrefix *bool
	NotFirstRun  bool
}

func (p *Project) v4() *contracts.AzionApplicationOptions {
	return &contracts.AzionApplicationOptions{
		Name: p.Name, Bucket: p.Bucket, Prefix: p.Prefix, Preset: p.Preset,
		RotatePrefix: p.RotatePrefix, NotFirstRun: p.NotFirstRun,
	}
}

func (p *Project) from(conf *contracts.AzionApplicationOptions) {
	p.Bucket, p.Prefix = conf.Bucket, conf.Prefix
}

// ensureTiming makes the v4 timing summary usable: the shared code records
// into it, and a v6 run does not otherwise start it.
func ensureTiming() {
	if GlobalTimingSummary == nil {
		InitTimingSummary()
	}
}

// ForV6 prepares a DeployCmd for use by the v6 deploy. The shared helpers save
// azion.json in the v4 format when they finish; for a v6 project that would
// overwrite its own file, so the v6 deploy, which saves after each step, turns
// that off here.
//
// It also switches the bundler to the v6 one (vulcan.NewVulcanV6), so the
// manifest comes out in the v6 format.
func (cmd *DeployCmd) ForV6() *DeployCmd {
	cmd.WriteAzionJsonContent = func(*contracts.AzionApplicationOptions, string) error { return nil }
	cmd.NewVulcan = vulcanPkg.NewVulcanV6
	ensureTiming()
	return cmd
}

// useV6Build makes the build step run the v6 bundler. The build command reads
// azion.json in the v4 format for the preset and writes it back when it
// finishes; on a v6 project that would fail on a first run and overwrite the
// file afterwards, so it is given the preset in memory and nothing to write.
func (cmd *DeployCmd) useV6Build(p *Project) {
	cmd.BuildCmd = func(f *cmdutil.Factory) *build.BuildCmd {
		b := build.NewBuildCmd(f)
		b.NewVulcan = vulcanPkg.NewVulcanV6
		b.GetAzionJsonContent = func(string) (*contracts.AzionApplicationOptions, error) {
			return &contracts.AzionApplicationOptions{Name: p.Name, Preset: p.Preset}, nil
		}
		b.WriteAzionJsonContent = func(*contracts.AzionApplicationOptions, string) error { return nil }
		return b
	}
}

// ManifestMissing reports whether .edge/manifest.json is absent.
func ManifestMissing() bool { return manifestMissing() }

// GenerateManifest writes the manifest without building, like init does.
func (cmd *DeployCmd) GenerateManifest(p *Project, msgs *[]string) error {
	s := &DeployState{Cmd: cmd, F: cmd.F, Conf: p.v4(), Msgs: msgs}
	return s.generateManifest()
}

// Build rotates the prefix and builds a project that was deployed before.
func (cmd *DeployCmd) Build(p *Project, oldPrefix, newPrefix string, msgs *[]string) error {
	cmd.useV6Build(p)
	conf := p.v4()
	s := &DeployState{Cmd: cmd, F: cmd.F, Conf: conf, Msgs: msgs, OldPrefix: oldPrefix, NewPrefix: newPrefix}
	if err := s.build(); err != nil {
		return err
	}
	p.from(conf)
	return nil
}

// BundlerInit initialises the bundler on a first run and builds.
func (cmd *DeployCmd) BundlerInit(p *Project, newPrefix string, msgs *[]string) error {
	cmd.useV6Build(p)
	conf := p.v4()
	s := &DeployState{Cmd: cmd, F: cmd.F, Conf: conf, Msgs: msgs, NewPrefix: newPrefix}
	if err := s.bundlerInit(); err != nil {
		return err
	}
	p.from(conf)
	return nil
}

// EnsureBucket creates the project's bucket when it has none.
func (cmd *DeployCmd) EnsureBucket(ctx context.Context, client *storage.Client, p *Project, msgs *[]string, manifest []contracts.StorageManifest) error {
	conf := p.v4()
	if err := cmd.doBucket(client, ctx, conf, msgs, manifest); err != nil {
		return err
	}
	p.from(conf)
	return nil
}

// StorageClient is the client EnsureBucket takes.
func StorageClient(cmd *DeployCmd) *storage.Client { return registry.Storage(cmd.F) }

// ShouldUpload reports whether this run uploads static files, with the same
// rules as the v4 deploy.
func (cmd *DeployCmd) ShouldUpload() bool {
	return shouldUpload(&DeployState{Cmd: cmd})
}

// UploadStatic uploads each storage directory of the manifest to the bucket.
func (cmd *DeployCmd) UploadStatic(p *Project, msgs *[]string, dirs []contracts.StorageManifest, creds token.S3Credentials) error {
	conf := p.v4()
	for _, st := range dirs {
		if err := cmd.uploadFilesWithCreds(cmd.F, conf, msgs, st.Dir, p.Bucket, creds); err != nil {
			return err
		}
	}
	return nil
}
