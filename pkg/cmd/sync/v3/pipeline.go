package sync

import (
	"context"
	"fmt"
	"os"

	msg "github.com/aziontech/azion-cli/messages/sync"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/pkg/pipeline"
	vulcanPkg "github.com/aziontech/azion-cli/pkg/vulcan"
)

// SyncState is what the v3 sync steps read and write.
//
// RemoteCaches and RemoteOrigins are the reason the order matters: the rules
// step resolves its behaviour targets against what the two steps before it
// found remotely.
type SyncState struct {
	F    *cmdutil.Factory
	Info contracts.SyncOptsV3
	Sync *SyncCmd

	Manifest      *contracts.Manifest
	RemoteCaches  map[string]contracts.AzionJsonDataCacheSettings
	RemoteOrigins map[string]contracts.AzionJsonDataOrigin
}

// V3SyncPipeline is the sync flow for the v3 API.
//
// Unlike v4 it does not read the project manifest: it starts from an empty one
// and fills it from the remote resources it finds.
type V3SyncPipeline struct{}

func (V3SyncPipeline) Name() string { return "sync-v3" }

func (V3SyncPipeline) Steps() []pipeline.Step[SyncState] {
	return []pipeline.Step[SyncState]{
		{
			Name: "Cache",
			Run: func(_ context.Context, s *SyncState) error {
				remote, err := s.Sync.syncCache(s.Info, s.F, s.Manifest)
				if err != nil {
					return fmt.Errorf(msg.ERRORSYNC, err.Error())
				}
				s.RemoteCaches = remote
				return nil
			},
		},
		{
			Name: "Origin",
			Run: func(_ context.Context, s *SyncState) error {
				remote, err := s.Sync.syncOrigin(s.Info, s.F, s.Manifest)
				if err != nil {
					return fmt.Errorf(msg.ERRORSYNC, err.Error())
				}
				s.RemoteOrigins = remote
				return nil
			},
		},
		{
			Name: "Rules",
			Run: func(_ context.Context, s *SyncState) error {
				err := s.Sync.syncRules(s.Info, s.F, s.Manifest, s.RemoteCaches, s.RemoteOrigins)
				if err != nil {
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

	vul := vulcanPkg.NewVulcanV3()
	command := vul.Command("", "manifest -o %s transform %s", f)
	return synch.CommandRunInteractive(f, fmt.Sprintf(command, fileName, "manifesttoconvert.json"))
}
