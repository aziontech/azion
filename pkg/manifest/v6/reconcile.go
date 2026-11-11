package v6

import (
	"context"
	"encoding/json"
	"fmt"

	msg "github.com/aziontech/azion-cli/messages/manifest"

	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	"github.com/aziontech/azion-cli/pkg/contracts"
)

// sub is one kind of sub-resource inside a Draft, as the client exposes it.
// Each field is a client method value, for example
// sub[sdk.CacheSettingRequest]{list: ..., create: c.CreateCacheSetting, ...}.
type sub[Req any] struct {
	// kind names the sub-resource in the deploy output ("cache setting").
	kind   string
	list   func(ctx context.Context, id int64, versionID string) ([]version.Item, error)
	create func(ctx context.Context, id int64, versionID string, req Req) (version.Item, error)
	update func(ctx context.Context, id int64, versionID string, subID int64, req Req) (version.Item, error)
	delete func(ctx context.Context, id int64, versionID string, subID int64) error
	// order is nil for kinds without an order.
	order func(ctx context.Context, id int64, versionID string, ids []int64) error
}

// page is a list response of the v6 SDK.
type page[T any] interface {
	GetResults() []T
	GetNext() string
}

// listAll follows the pages of a sub-resource list and returns every item.
// fetch is the client's List method for one page.
func listAll[T any, P interface {
	*T
	version.Item
}, L page[T]](ctx context.Context, fetch func(ctx context.Context, opts *contracts.ListOptions) (L, error)) ([]version.Item, error) {
	var all []version.Item
	for n := int64(1); ; n++ {
		res, err := fetch(ctx, &contracts.ListOptions{Page: n, PageSize: 100})
		if err != nil {
			return nil, err
		}
		results := res.GetResults()
		for i := range results {
			all = append(all, P(&results[i]))
		}
		if res.GetNext() == "" || len(results) == 0 {
			return all, nil
		}
	}
}

// reconcile makes the sub-resources of one kind inside a Draft match the
// manifest, by name: it updates the ones that exist, creates the missing ones
// and deletes the ones the manifest no longer has. Deleting is safe because the
// Draft is not live. When the kind is ordered, the order becomes the
// manifest's.
//
// It returns the id of each desired item by name, for the references later
// writes need. Each create, update and removal is reported through say, which
// names the resource that owns the Draft (Applier.within).
func reconcile[Req any](ctx context.Context, say func(string), s sub[Req], id int64, vid string, desired []Req) (map[string]int64, error) {
	existing, err := s.list(ctx, id, vid)
	if err != nil {
		return nil, err
	}
	byName := make(map[string]int64, len(existing))
	type named struct {
		name string
		id   int64
	}
	var dups []named
	for _, it := range existing {
		if _, dup := byName[it.GetName()]; dup {
			// A name twice in the Draft: keep the first, drop the rest.
			dups = append(dups, named{it.GetName(), it.GetId()})
			continue
		}
		byName[it.GetName()] = it.GetId()
	}

	ids := make(map[string]int64, len(desired))
	order := make([]int64, 0, len(desired))
	for _, req := range desired {
		name := nameOf(req)
		var item version.Item
		updated := false
		if subID, ok := byName[name]; ok {
			item, err = s.update(ctx, id, vid, subID, req)
			delete(byName, name)
			updated = true
		} else {
			item, err = s.create(ctx, id, vid, req)
		}
		if err != nil {
			return nil, err
		}
		if updated {
			say(fmt.Sprintf(msg.V6SubUpdated, s.kind, name, item.GetId()))
		} else {
			say(fmt.Sprintf(msg.V6SubCreated, s.kind, name, item.GetId()))
		}
		ids[name] = item.GetId()
		order = append(order, item.GetId())
	}

	for _, d := range dups {
		if err := s.delete(ctx, id, vid, d.id); err != nil {
			return nil, err
		}
		say(fmt.Sprintf(msg.V6SubDuplicateRemoved, s.kind, d.name, d.id))
	}
	for name, subID := range byName {
		if err := s.delete(ctx, id, vid, subID); err != nil {
			return nil, err
		}
		say(fmt.Sprintf(msg.V6SubRemoved, s.kind, name, subID))
	}

	if s.order != nil && len(order) > 1 {
		if err := s.order(ctx, id, vid, order); err != nil {
			return nil, err
		}
		say(fmt.Sprintf(msg.V6SubsOrdered, len(order), s.kind))
	}
	return ids, nil
}

// nameOf reads the name of any request body. Every resource and sub-resource
// in the manifest is identified by its name.
func nameOf(req any) string {
	var n struct {
		Name string `json:"name"`
	}
	data, err := json.Marshal(req)
	if err != nil {
		return ""
	}
	_ = json.Unmarshal(data, &n)
	return n.Name
}
