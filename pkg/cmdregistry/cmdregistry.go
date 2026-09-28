// Package cmdregistry expresses which sub-commands a dispatcher offers, per API
// generation, as data rather than as a second copy of the dispatcher.

// A table lives next to the dispatcher it feeds (pkg/cmd/create/table.go, and so
// on) so that adding a resource stays a one-line edit in the package that owns
// the verb.
package cmdregistry

import (
	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/spf13/cobra"
)

// Builder constructs a command. It is the signature every NewCmd already has.
type Builder func(f *cmdutil.Factory) *cobra.Command

// Entry is one sub-command and the generations it is offered to.
type Entry struct {
	// Versions limits the entry to those generations. An empty Versions means
	// the command is version-agnostic and is offered to every generation.
	Versions []apiversion.Version
	// New constructs the command.
	New Builder
}

// Supports reports whether the entry is offered to generation v.
func (e Entry) Supports(v apiversion.Version) bool {
	if len(e.Versions) == 0 {
		return true
	}
	for _, ev := range e.Versions {
		if ev == v {
			return true
		}
	}
	return false
}

// Any registers a command that behaves the same on every generation.
func Any(new Builder) Entry {
	return Entry{New: new}
}

// V3Only registers a command that exists only on the v3 generation.
func V3Only(new Builder) Entry {
	return Entry{Versions: []apiversion.Version{apiversion.V3}, New: new}
}

// V4Only registers a command that exists only on the v4 generation.
func V4Only(new Builder) Entry {
	return Entry{Versions: []apiversion.Version{apiversion.V4}, New: new}
}

// Fallback is the generation assumed when a factory carries no resolved
// version. The root command always resolves one, so this covers callers that
// build a command tree on their own — documentation generation and unit tests —
// which saw the v4 tree before these tables existed.
const Fallback = apiversion.V4

func versionOf(f *cmdutil.Factory) apiversion.Version {
	if f.APIVersion == "" {
		return Fallback
	}
	return f.APIVersion
}

// Children builds the entries that the factory's generation is entitled to, in
// table order. Cobra sorts sub-commands for display, so the order here only
// affects the order of construction.
func Children(entries []Entry, f *cmdutil.Factory) []*cobra.Command {
	v := versionOf(f)
	cmds := make([]*cobra.Command, 0, len(entries))
	for _, e := range entries {
		if e.Supports(v) {
			cmds = append(cmds, e.New(f))
		}
	}
	return cmds
}

// Example picks the example block for the factory's generation.
func Example(examples map[apiversion.Version]string, f *cmdutil.Factory) string {
	return examples[versionOf(f)]
}
