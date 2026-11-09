package root

import (
	"bytes"
	"net/http"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/apiversion"
	"github.com/aziontech/azion-cli/pkg/cmdregistry"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/iostreams"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func tableFactory(v apiversion.Version) *cmdutil.Factory {
	return &cmdutil.Factory{
		HttpClient: &http.Client{Timeout: 50 * time.Second},
		IOStreams: &iostreams.IOStreams{
			Out: &bytes.Buffer{},
			Err: &bytes.Buffer{},
		},
		Config:     viper.New(),
		APIVersion: v,
	}
}

// TestTableNamesAreUniquePerVersion guards the reason the tables can hold both
// generations at once: several resources are registered twice under the same
// CLI name (cache-setting and rules-engine have a v3 and a v4 entry), and only
// the version filter keeps one of each. Cobra accepts a duplicate sub-command
// silently and then dispatches to whichever was added first, so a table entry
// that forgets its Versions would ship a command that quietly runs against the
// wrong API generation. This walks the whole resolved tree for each generation
// and fails on any name registered twice under the same parent.
func TestTableNamesAreUniquePerVersion(t *testing.T) {
	for _, v := range []apiversion.Version{apiversion.V3, apiversion.V4, apiversion.V6} {
		t.Run(string(v), func(t *testing.T) {
			f := tableFactory(v)
			for _, cmd := range cmdregistry.Children(children, f) {
				assertUniqueChildren(t, cmd, cmd.Name())
			}
		})
	}
}

func assertUniqueChildren(t *testing.T, cmd *cobra.Command, path string) {
	t.Helper()
	seen := make(map[string]bool, len(cmd.Commands()))
	for _, sub := range cmd.Commands() {
		name := sub.Name()
		if seen[name] {
			t.Errorf("%s registers %q more than once", path, name)
		}
		seen[name] = true
		assertUniqueChildren(t, sub, path+" "+name)
	}
}

// TestTopLevelNamesAreUniquePerVersion covers the root's own table, which is
// where the two generations' project commands (deploy, init, and the rest) sit
// side by side under the same names.
func TestTopLevelNamesAreUniquePerVersion(t *testing.T) {
	for _, v := range []apiversion.Version{apiversion.V3, apiversion.V4, apiversion.V6} {
		t.Run(string(v), func(t *testing.T) {
			seen := make(map[string]bool)
			for _, cmd := range cmdregistry.Children(children, tableFactory(v)) {
				if seen[cmd.Name()] {
					t.Errorf("root registers %q more than once on %s", cmd.Name(), v)
				}
				seen[cmd.Name()] = true
			}
		})
	}
}

// TestEveryGenerationGetsBothTrees is a smoke test on the resolved surface: a
// regression that dropped the version filter, or one that left a table empty,
// would show up here as a command count far from what each generation ships.
func TestEveryGenerationGetsBothTrees(t *testing.T) {
	counts := map[apiversion.Version]int{}
	for _, v := range []apiversion.Version{apiversion.V3, apiversion.V4} {
		counts[v] = len(cmdregistry.Children(children, tableFactory(v)))
	}
	if counts[apiversion.V3] != 22 {
		t.Errorf("v3 root: got %d top-level commands, want 22", counts[apiversion.V3])
	}
	if counts[apiversion.V4] != 25 {
		t.Errorf("v4 root: got %d top-level commands, want 25", counts[apiversion.V4])
	}
}

// TestV6GetsOnlyTheProjectAndGenericCommands pins what a v6 account sees: the
// project commands and the ones that touch no resource API. The CRUD verbs
// and every other command that calls the v4 API stay off the list until they
// have a v6 implementation; adding one is a deliberate edit to this test.
func TestV6GetsOnlyTheProjectAndGenericCommands(t *testing.T) {
	byName := func(v apiversion.Version) map[string]*cobra.Command {
		out := map[string]*cobra.Command{}
		for _, cmd := range cmdregistry.Children(children, tableFactory(v)) {
			out[cmd.Name()] = cmd
		}
		return out
	}
	v4, v6 := byName(apiversion.V4), byName(apiversion.V6)

	want := []string{
		"build", "completion", "create", "delete", "deploy", "dev", "init",
		"link", "login", "logout", "profiles", "reset", "version", "whoami",
	}
	got := make([]string, 0, len(v6))
	for name := range v6 {
		got = append(got, name)
	}
	sort.Strings(got)
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Errorf("v6 root:\n got  %v\n want %v", got, want)
	}

	// create and delete are there for profiles alone: no resource command.
	for _, verb := range []string{"create", "delete"} {
		var subs []string
		for _, c := range v6[verb].Commands() {
			subs = append(subs, c.Name())
		}
		if strings.Join(subs, ",") != "profile" {
			t.Errorf("v6 %s offers %v, want only [profile]", verb, subs)
		}
	}

	// deploy is the v6 implementation; it has the flags only the v6
	// pipeline understands.
	for _, flag := range []string{"environment", "no-activate", "timeout"} {
		if v6["deploy"].Flags().Lookup(flag) == nil {
			t.Errorf("v6 deploy has no --%s: it is not the v6 implementation", flag)
		}
		if v4["deploy"].Flags().Lookup(flag) != nil {
			t.Errorf("v4 deploy has --%s: v4 is being served the v6 implementation", flag)
		}
	}
}
