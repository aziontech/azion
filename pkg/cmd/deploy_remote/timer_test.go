package deploy

import (
	"reflect"
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/manifest"
)

// durations reads every exported time.Duration field of the timing summary
// without copying the struct, which holds a mutex.
func durations(ts *TimingSummary) map[string]time.Duration {
	v := reflect.ValueOf(ts).Elem()
	t := v.Type()
	out := make(map[string]time.Duration, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() || f.Type != reflect.TypeOf(time.Duration(0)) {
			continue
		}
		out[f.Name] = time.Duration(v.Field(i).Int())
	}
	return out
}

// records reports whether the timing summary has a field for this step name.
//
// HandleManifestTimingCallback ends in "default: // Unknown timing name,
// ignore", so an unrecognised name is dropped in silence. This calls it and
// looks for the effect.
func records(t *testing.T, step string) bool {
	t.Helper()
	InitTimingSummary()
	before := durations(GlobalTimingSummary)
	HandleManifestTimingCallback(step, 7*time.Second)
	after := durations(GlobalTimingSummary)
	for name, d := range after {
		if d != before[name] {
			return true
		}
	}
	return false
}

// TestStepNamesAreRecordedByTheTimingSummary ties the v4 pipeline's step names
// to the switch in HandleManifestTimingCallback.
//
// The two are coupled only by string equality: renaming a step compiles fine on
// both sides and silently stops recording that step's duration, because the
// switch ignores names it does not know. Nothing but this test notices.
func TestStepNamesAreRecordedByTheTimingSummary(t *testing.T) {
	// The steps the summary has a field for. The trailing cleanup step is
	// deliberately absent: the sequence this pipeline replaced never timed
	// orphan removal either.
	want := map[string]bool{
		"ManifestFunctions":           true,
		"ManifestFunctionInstances":   true,
		"ManifestEdgeApplication":     true,
		"ManifestCacheSettings":       true,
		"ManifestConnectors":          true,
		"ManifestRulesEngine":         true,
		"ManifestWorkloads":           true,
		"ManifestWorkloadDeployments": true,
		"ManifestFirewalls":           true,
		"ManifestPurge":               true,
	}

	timed := map[string]bool{}
	for _, step := range (manifest.V4Pipeline{}).Steps() {
		if records(t, step.Name) {
			timed[step.Name] = true
		}
		delete(want, step.Name)
	}

	// Every name the summary knows about must still be produced by a step.
	for name := range want {
		t.Errorf("the timing summary records %q but no pipeline step emits it", name)
	}

	for name := range timed {
		delete(want, name)
	}
	if len(timed) != 10 {
		t.Errorf("%d pipeline steps are recorded by the timing summary, want 10: %v", len(timed), timed)
	}
}
