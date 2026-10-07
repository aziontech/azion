package sync

import (
	"context"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/sync"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/manifest"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
	edgesdk "github.com/aziontech/azionapi-v4-go-sdk-dev/azion-api"
)

// SyncState is what the v4 sync steps read and write.
type SyncState struct {
	F    *cmdutil.Factory
	Info contracts.SyncOpts
	Sync *SyncCmd

	// Manifest is produced by the read step and consumed by the ones that
	// reconcile against it.
	Manifest *contracts.ManifestV4
}

// V4SyncPipeline is the sync flow for the v4 API.
type V4SyncPipeline struct{}

func (V4SyncPipeline) Name() string { return "sync-v4" }

func (V4SyncPipeline) Steps() []pipeline.Step[SyncState] {
	return []pipeline.Step[SyncState]{
		{
			Name: "ReadManifest",
			Run:  func(_ context.Context, s *SyncState) error { return s.readManifest() },
		},
		{
			Name: "Cache",
			Run: func(_ context.Context, s *SyncState) error {
				if _, err := s.Sync.syncCache(s.Info, s.F, s.Manifest); err != nil {
					return fmt.Errorf(msg.ERRORSYNC, err.Error())
				}
				return nil
			},
		},
		{
			Name: "Rules",
			Run: func(_ context.Context, s *SyncState) error {
				if err := s.Sync.syncRules(s.Info, s.F, s.Manifest); err != nil {
					return fmt.Errorf(msg.ERRORSYNC, err.Error())
				}
				return nil
			},
		},
		{
			Name: "Env",
			Run: func(_ context.Context, s *SyncState) error {
				if err := s.Sync.syncEnv(s.F); err != nil {
					return fmt.Errorf(msg.ERRORSYNC, err.Error())
				}
				return nil
			},
		},
		{
			Name:    "IaC",
			Applies: func(s *SyncState) bool { return s.Sync.IaC },
			Run:     func(_ context.Context, s *SyncState) error { return s.writeIaC() },
		},
	}
}

// emptyManifest is what sync reconciles against when the project has no
// manifest to read, so that every remote resource counts as unknown locally.
func emptyManifest() *contracts.ManifestV4 {
	return &contracts.ManifestV4{
		Applications:        []contracts.Applications{},
		Workloads:           []contracts.WorkloadManifest{},
		WorkloadDeployments: []contracts.WorkloadDeployment{},
		Purge:               []contracts.PurgeManifest{},
		Storage:             []contracts.StorageManifest{},
		Functions:           []contracts.Function{},
		Connectors:          []edgesdk.ConnectorRequest{},
	}
}

// readManifest loads the project manifest, falling back to an empty one if it
// cannot be located or parsed. Neither failure stops the sync.
func (s *SyncState) readManifest() error {
	var msgs []string
	interpreter := manifest.NewManifestInterpreter()

	pathManifest, err := interpreter.ManifestPath()
	if err != nil {
		s.Manifest = emptyManifest()
		return nil
	}

	s.Manifest, err = interpreter.ReadManifest(pathManifest, s.F, &msgs)
	if err != nil {
		s.Manifest = emptyManifest()
	}
	return nil
}

// writeIaC converts the manifest into an azion.config file.
func (s *SyncState) writeIaC() error {
	synch, f := s.Sync, s.F
	if synch.IaCFormat != "mjs" && synch.IaCFormat != "cjs" && synch.IaCFormat != "js" && synch.IaCFormat != "ts" {
		return msg.INVALIDFORMAT
	}

	if err := synch.WriteManifest(s.Manifest, ""); err != nil {
		return err
	}
	defer os.Remove("manifesttoconvert.json")
	fileName := fmt.Sprintf("azion.config.%s", synch.IaCFormat)

	vul := vulcanPkg.NewVulcan()
	command := vul.Command("", "manifest transform --output %s --entry %s", f)
	return synch.CommandRunInteractive(f, fmt.Sprintf(command, fileName, "manifesttoconvert.json"))
}
