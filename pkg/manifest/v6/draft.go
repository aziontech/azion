// Package v6 applies a v6 manifest. Where the v4 applier writes straight to
// live resources, this one writes into Drafts: for each resource it finds or
// opens a Draft version, writes the whole configuration from the manifest into
// it, reconciles the sub-resources inside it and checks that nobody else took
// the Draft meanwhile. Nothing it does reaches traffic; the build and release
// steps in pkg/cmd/deploy_remote/v6 do that.
package v6

import (
	"context"
	"errors"
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/manifest"
	apiApplications "github.com/aziontech/azion-cli/pkg/api/v6/applications"
	apiConnectors "github.com/aziontech/azion-cli/pkg/api/v6/connectors"
	apiCustomPages "github.com/aziontech/azion-cli/pkg/api/v6/custompages"
	apiFirewalls "github.com/aziontech/azion-cli/pkg/api/v6/firewalls"
	apiFunctions "github.com/aziontech/azion-cli/pkg/api/v6/functions"
	apiNetworkLists "github.com/aziontech/azion-cli/pkg/api/v6/networklists"
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	apiWafs "github.com/aziontech/azion-cli/pkg/api/v6/wafs"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
)

// DraftComment is recorded on every version the CLI opens, so the Console
// shows where it came from.
const DraftComment = "azion-cli deploy"

// Pending is a Draft this run wrote and that still needs a build.
type Pending struct {
	Kind      string
	Name      string
	ID        int64
	VersionID string
	Client    version.Versioned
	// Record is the azion.json list the resource lives in, so the build step
	// can store the state it reaches.
	Record *[]contracts.AzionJsonV6Resource
}

// SetState stores the build state reached in azion.json.
func (p Pending) SetState(s version.State) {
	if r := find(p.Record, p.Name); r != nil {
		r.VersionState = string(s)
	}
}

// resourceClient is what the applier needs from the client of one versioned
// resource kind. Every client in pkg/api/v6 has these methods; the interface
// lets one draft function serve all seven.
type resourceClient[Req any] interface {
	version.Versioned
	Create(ctx context.Context, req Req) (version.Object, error)
	FindByName(ctx context.Context, name string) (version.Object, bool, error)
	GetVersion(ctx context.Context, id int64, versionID string) (version.Meta, error)
	UpdateVersion(ctx context.Context, id int64, versionID string, req Req) (version.Meta, error)
}

var (
	_ resourceClient[sdk.ApplicationRequest] = (*apiApplications.Client)(nil)
	_ resourceClient[sdk.FunctionsRequest]   = (*apiFunctions.Client)(nil)
	_ resourceClient[sdk.ConnectorRequest]   = (*apiConnectors.Client)(nil)
	_ resourceClient[sdk.FirewallRequest]    = (*apiFirewalls.Client)(nil)
	_ resourceClient[sdk.CustomPageRequest]  = (*apiCustomPages.Client)(nil)
	_ resourceClient[sdk.NetworkListRequest] = (*apiNetworkLists.Client)(nil)
	_ resourceClient[sdk.WAFRequest]         = (*apiWafs.Client)(nil)
)

// editable reports whether a version can still be written. A build that
// failed or was canceled leaves the version editable, and reusing it keeps
// failed runs from piling up Drafts (the API allows 20 per resource).
func editable(s version.State) bool {
	return s == version.Draft || s == version.Canceled || s == version.Error
}

// draft finds or opens the Draft that this run writes req into, following the
// rules in doc/V6-DEPLOY-PLAN.md section 5, and writes it. Sub-resources are
// left to the caller, which then calls verify.
func draft[Req any](ctx context.Context, a *Applier, kind string, r resourceClient[Req], list *[]contracts.AzionJsonV6Resource, name string, req Req) (Pending, error) {
	rec := find(list, name)

	var id int64
	var vid string
	created := false

	switch {
	case rec != nil && rec.ID != 0:
		id = rec.ID
	default:
		m, found, err := r.FindByName(ctx, name)
		if err != nil {
			return Pending{}, err
		}
		if found {
			id = m.GetId()
			a.say(fmt.Sprintf(msg.V6DraftAdopted, kind, name, id))
		}
	}

	if id != 0 && rec != nil && rec.ID == id && rec.VersionID != "" {
		m, err := r.GetVersion(ctx, id, rec.VersionID)
		switch {
		case err == nil && editable(version.State(m.GetVersionState())):
			vid = rec.VersionID
			a.say(fmt.Sprintf(msg.V6DraftReused, kind, name, vid))
		case err != nil && !errors.Is(err, utils.ErrorNotFound404):
			return Pending{}, err
		}
	}

	if id != 0 && vid == "" {
		v, err := r.CreateVersion(ctx, id, "", DraftComment)
		switch {
		case err == nil:
			vid = v
			a.say(fmt.Sprintf(msg.V6DraftOpened, kind, name, vid))
		case errors.Is(err, utils.ErrorNotFound404):
			a.say(fmt.Sprintf(msg.V6ResourceRecreate, kind, name, id))
			id = 0
		default:
			return Pending{}, err
		}
	}

	if id == 0 {
		m, err := r.Create(ctx, req)
		if err != nil {
			return Pending{}, err
		}
		id, vid, created = m.GetId(), m.GetVersionId(), true
		a.say(fmt.Sprintf(msg.V6DraftCreated, kind, name, id, vid))
	}

	if !created {
		if _, err := r.UpdateVersion(ctx, id, vid, req); err != nil {
			return Pending{}, err
		}
		a.say(fmt.Sprintf(msg.V6DraftWritten, kind, name, vid))
	}

	upsert(list, contracts.AzionJsonV6Resource{ID: id, Name: name, VersionID: vid, VersionState: string(version.Draft)})
	a.writeConfig()

	return Pending{Kind: kind, Name: name, ID: id, VersionID: vid, Client: r, Record: list}, nil
}

// verify is the first check of section 6: after everything was written, the
// version is still the expected editable Draft. Another process starting its
// build in the middle of the run would otherwise have the CLI write over it.
func verify[Req any](ctx context.Context, r resourceClient[Req], p Pending) error {
	m, err := r.GetVersion(ctx, p.ID, p.VersionID)
	if err != nil {
		return err
	}
	state := version.State(m.GetVersionState())
	if (m.GetVersionId() != "" && m.GetVersionId() != p.VersionID) || !editable(state) {
		return fmt.Errorf("%s %q, version %s is %s: %w", p.Kind, p.Name, p.VersionID, state, msg.ErrorV6DraftTaken)
	}
	p.SetState(state)
	return nil
}

func find(list *[]contracts.AzionJsonV6Resource, name string) *contracts.AzionJsonV6Resource {
	for i := range *list {
		if (*list)[i].Name == name {
			return &(*list)[i]
		}
	}
	return nil
}

func upsert(list *[]contracts.AzionJsonV6Resource, rec contracts.AzionJsonV6Resource) {
	if r := find(list, rec.Name); r != nil {
		*r = rec
		return
	}
	*list = append(*list, rec)
}
