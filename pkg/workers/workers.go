package workers

import "runtime"

const (
	MaxWorkers     = 20
	MinWorkers     = 5
	WorkersPerCore = 3
)

// CalculateOptimal calculates the optimal number of workers based on CPU cores
// If workers is 0 (not set), it auto-calculates based on CPU cores
// The result is clamped between MinWorkers and MaxWorkers
func CalculateOptimal(workers int) int {
	if workers > 0 {
		// User specified a value, clamp it to reasonable bounds
		if workers > MaxWorkers {
			return MaxWorkers
		}
		if workers < 1 {
			return MinWorkers
		}
		return workers
	}

	cpuCores := runtime.NumCPU()
	optimal := cpuCores * WorkersPerCore

	if optimal > MaxWorkers {
		optimal = MaxWorkers
	}
	if optimal < MinWorkers {
		optimal = MinWorkers
	}

	return optimal
}
