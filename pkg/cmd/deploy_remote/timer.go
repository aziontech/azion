package deploy

import (
	"fmt"
	"sync"
	"time"

	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)

// TimingSummary holds the timing information for each phase of the deploy process
type TimingSummary struct {
	mu                    sync.Mutex
	TotalDeployTime       time.Duration
	UploadStaticFilesTime time.Duration
	BucketCreateTime      time.Duration
	CredentialsTime       time.Duration
	ManifestCreateTime    time.Duration
	ReadManifestTime      time.Duration
	// Manifest resource timings
	ManifestFunctionsTime           time.Duration
	ManifestFunctionInstancesTime   time.Duration
	ManifestEdgeApplicationTime     time.Duration
	ManifestCacheSettingsTime       time.Duration
	ManifestConnectorsTime          time.Duration
	ManifestRulesEngineTime         time.Duration
	ManifestWorkloadsTime           time.Duration
	ManifestWorkloadDeploymentsTime time.Duration
	ManifestFirewallsTime           time.Duration
	ManifestPurgeTime               time.Duration
	APICallTimes                    []APICallTiming
}

// APICallTiming holds timing information for individual API calls
type APICallTiming struct {
	Name     string
	Duration time.Duration
}

// Timer is a simple utility for timing operations
type Timer struct {
	startTime time.Time
	name      string
}

// NewTimer creates a new timer and starts it
func NewTimer(name string) *Timer {
	return &Timer{
		startTime: time.Now(),
		name:      name,
	}
}

// Elapsed returns the time elapsed since the timer was started
func (t *Timer) Elapsed() time.Duration {
	return time.Since(t.startTime)
}

// Stop stops the timer and returns the elapsed time
func (t *Timer) Stop() time.Duration {
	elapsed := time.Since(t.startTime)
	logger.Debug("Timer stopped", zap.String("name", t.name), zap.Duration("elapsed", elapsed))
	return elapsed
}

// GlobalTimingSummary is the global timing summary for the deploy process
var GlobalTimingSummary *TimingSummary

// InitTimingSummary initializes the global timing summary
func InitTimingSummary() {
	GlobalTimingSummary = &TimingSummary{
		APICallTimes: make([]APICallTiming, 0),
	}
}

// AddAPICallTime adds an API call timing to the summary
func (ts *TimingSummary) AddAPICallTime(name string, duration time.Duration) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	ts.APICallTimes = append(ts.APICallTimes, APICallTiming{
		Name:     name,
		Duration: duration,
	})
}

// PrintSummary prints the timing summary in debug format
func (ts *TimingSummary) PrintSummary() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	logger.Debug("========== DEPLOY TIMING SUMMARY ==========")
	logger.Debug(fmt.Sprintf("Total Deploy Time: %v", ts.TotalDeployTime))

	// Only print non-zero non-manifest timings
	if ts.UploadStaticFilesTime > 0 {
		logger.Debug(fmt.Sprintf("Upload Static Files: %v", ts.UploadStaticFilesTime))
	}
	if ts.BucketCreateTime > 0 {
		logger.Debug(fmt.Sprintf("Bucket Creation: %v", ts.BucketCreateTime))
	}
	if ts.CredentialsTime > 0 {
		logger.Debug(fmt.Sprintf("Credentials: %v", ts.CredentialsTime))
	}
	if ts.ReadManifestTime > 0 {
		logger.Debug(fmt.Sprintf("Read Manifest: %v", ts.ReadManifestTime))
	}
	if ts.ManifestCreateTime > 0 {
		logger.Debug(fmt.Sprintf("Manifest Apply: %v", ts.ManifestCreateTime))
	}

	// Always print manifest resource breakdown
	logger.Debug("--- Manifest Resource Timings ---")
	logger.Debug(fmt.Sprintf("  Functions: %v", ts.ManifestFunctionsTime))
	logger.Debug(fmt.Sprintf("  Function Instances: %v", ts.ManifestFunctionInstancesTime))
	logger.Debug(fmt.Sprintf("  Application: %v", ts.ManifestEdgeApplicationTime))
	logger.Debug(fmt.Sprintf("  Cache Settings: %v", ts.ManifestCacheSettingsTime))
	logger.Debug(fmt.Sprintf("  Connectors: %v", ts.ManifestConnectorsTime))
	logger.Debug(fmt.Sprintf("  Rules Engine: %v", ts.ManifestRulesEngineTime))
	logger.Debug(fmt.Sprintf("  Workloads: %v", ts.ManifestWorkloadsTime))
	logger.Debug(fmt.Sprintf("  Workload Deployments: %v", ts.ManifestWorkloadDeploymentsTime))
	logger.Debug(fmt.Sprintf("  Firewalls: %v", ts.ManifestFirewallsTime))
	logger.Debug(fmt.Sprintf("  Purge: %v", ts.ManifestPurgeTime))

	if len(ts.APICallTimes) > 0 {
		logger.Debug("---------- Individual API Call Times ----------")
		var totalAPITime time.Duration
		for _, call := range ts.APICallTimes {
			logger.Debug(fmt.Sprintf("API Call [%s]: %v", call.Name, call.Duration))
			totalAPITime += call.Duration
		}
		logger.Debug(fmt.Sprintf("Total API Calls Time: %v", totalAPITime))
	}

	logger.Debug("============================================")
}

// HandleManifestTimingCallback handles timing callbacks from the manifest package
func HandleManifestTimingCallback(name string, duration time.Duration) {
	if GlobalTimingSummary == nil {
		return
	}
	GlobalTimingSummary.mu.Lock()
	defer GlobalTimingSummary.mu.Unlock()

	switch name {
	case "ManifestFunctions":
		GlobalTimingSummary.ManifestFunctionsTime = duration
	case "ManifestFunctionInstances":
		GlobalTimingSummary.ManifestFunctionInstancesTime = duration
	case "ManifestEdgeApplication":
		GlobalTimingSummary.ManifestEdgeApplicationTime = duration
	case "ManifestCacheSettings":
		GlobalTimingSummary.ManifestCacheSettingsTime = duration
	case "ManifestConnectors":
		GlobalTimingSummary.ManifestConnectorsTime = duration
	case "ManifestRulesEngine":
		GlobalTimingSummary.ManifestRulesEngineTime = duration
	case "ManifestWorkloads":
		GlobalTimingSummary.ManifestWorkloadsTime = duration
	case "ManifestWorkloadDeployments":
		GlobalTimingSummary.ManifestWorkloadDeploymentsTime = duration
	case "ManifestFirewalls":
		GlobalTimingSummary.ManifestFirewallsTime = duration
	case "ManifestPurge":
		GlobalTimingSummary.ManifestPurgeTime = duration
	default:
		// Unknown timing name, ignore
	}
}

// GetTotalTime returns the total deploy time
func (ts *TimingSummary) GetTotalTime() time.Duration {
	return ts.TotalDeployTime
}

// TimeAPICall is a helper function to time an API call and record it
func TimeAPICall(name string, fn func() error) error {
	start := time.Now()
	err := fn()
	duration := time.Since(start)
	if GlobalTimingSummary != nil {
		GlobalTimingSummary.AddAPICallTime(name, duration)
	}
	return err
}

// TimeAPICallWithResult is a helper function to time an API call that returns a result and record it
func TimeAPICallWithResult[T any](name string, fn func() (T, error)) (T, error) {
	start := time.Now()
	result, err := fn()
	duration := time.Since(start)
	if GlobalTimingSummary != nil {
		GlobalTimingSummary.AddAPICallTime(name, duration)
	}
	return result, err
}
