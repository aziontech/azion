package v6

import (
	"fmt"
	"io"
	"time"

	msg "github.com/aziontech/azion-cli/messages/deploy-remote"
)

// stepLabels names each step in the timing summary. Every step of
// V6DeployPipeline must have one; TestEveryStepHasATimingLabel holds that, so
// renaming a step cannot silently drop it from the summary.
var stepLabels = map[string]string{
	StepGenerateManifest:        "Generate manifest",
	StepBuild:                   "Build",
	StepReadManifest:            "Read manifest",
	StepBucket:                  "Bucket",
	StepBundlerInit:             "Bundler init",
	StepReadManifestAgain:       "Read manifest",
	StepCredentials:             "Bucket credentials",
	StepUploadStaticFiles:       "Upload static files",
	StepManifestFunctions:       "Functions",
	StepManifestNetworkLists:    "Network lists",
	StepManifestWAFs:            "WAFs",
	StepManifestConnectors:      "Connectors",
	StepVersionBuildDeps:        "Build dependencies",
	StepManifestApplication:     "Application",
	StepManifestFirewalls:       "Firewall",
	StepManifestCustomPages:     "Custom page",
	StepVersionBuildPublishable: "Build application, firewall, custom page",
	StepManifestEnvironments:    "Environments",
	StepManifestDeployments:     "Deployments",
	StepReleaseCreate:           "Release create",
	StepReleaseBuild:            "Release build",
	StepReleaseActivate:         "Release activate",
	StepManifestWorkloads:       "Workloads",
	StepVerifyDeploy:            "Verify",
	StepManifestPurge:           "Purge",
	StepOrphans:                 "Orphan check",
}

// TimingSummaryV6 is the v6 deploy's timing report. The time spent waiting
// for builds is shown apart from the time of each step, since in v6 most of a
// deploy is waiting.
type TimingSummaryV6 struct {
	Steps     []StepTiming
	BuildWait time.Duration
	Total     time.Duration
}

type StepTiming struct {
	Label    string
	Duration time.Duration
}

// Record is the pipeline's timing callback. Two steps with the same label
// (both manifest reads) add up.
func (t *TimingSummaryV6) Record(name string, d time.Duration) {
	label := stepLabels[name]
	if label == "" {
		label = name
	}
	for i := range t.Steps {
		if t.Steps[i].Label == label {
			t.Steps[i].Duration += d
			return
		}
	}
	t.Steps = append(t.Steps, StepTiming{Label: label, Duration: d})
}

func (t *TimingSummaryV6) Print(w io.Writer) {
	fmt.Fprint(w, msg.V6TimingHeader)
	for _, s := range t.Steps {
		fmt.Fprintf(w, "  %-42s %8s\n", s.Label, s.Duration.Round(time.Millisecond))
	}
	fmt.Fprintf(w, "  %-42s %8s\n", "Of which waiting for builds", t.BuildWait.Round(time.Millisecond))
	fmt.Fprintf(w, "  %-42s %8s\n", "Total", t.Total.Round(time.Millisecond))
}
