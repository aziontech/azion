package deploy

import (
	"fmt"
	"sync"
	"time"

	"github.com/aziontech/azion-cli/pkg/logger"
)

// TimingSummary holds the timing information for each phase of the deploy process
type TimingSummary struct {
	mu               sync.Mutex
	TotalDeployTime  time.Duration
	BucketCreateTime time.Duration
	CredentialsTime  time.Duration
	UploadFilesTime  time.Duration
	ScriptCallTime   time.Duration
	LogsCaptureTime  time.Duration
	APICallTimes     []APICallTiming
}

// APICallTiming holds timing information for individual API calls
type APICallTiming struct {
	Name     string
	Duration time.Duration
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

	// Only print non-zero timings
	if ts.BucketCreateTime > 0 {
		logger.Debug(fmt.Sprintf("Bucket Creation: %v", ts.BucketCreateTime))
	}
	if ts.CredentialsTime > 0 {
		logger.Debug(fmt.Sprintf("Credentials Creation: %v", ts.CredentialsTime))
	}
	if ts.UploadFilesTime > 0 {
		logger.Debug(fmt.Sprintf("Upload Files: %v", ts.UploadFilesTime))
	}
	if ts.ScriptCallTime > 0 {
		logger.Debug(fmt.Sprintf("Script Call: %v", ts.ScriptCallTime))
	}
	if ts.LogsCaptureTime > 0 {
		logger.Debug(fmt.Sprintf("Logs Capture: %v", ts.LogsCaptureTime))
	}

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
