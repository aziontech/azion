package workers

import (
	"runtime"
	"testing"
)

func TestCalculateOptimal(t *testing.T) {
	tests := []struct {
		name     string
		workers  int
		expected int
	}{
		{
			name:     "user specified valid value",
			workers:  10,
			expected: 10,
		},
		{
			name:     "user specified value exceeds max",
			workers:  50,
			expected: MaxWorkers,
		},
		{
			name:     "user specified value below min",
			workers:  0, // 0 means auto-calculate
			expected: expectedAutoValue(),
		},
		{
			name:     "user specified value of 1",
			workers:  1,
			expected: 1,
		},
		{
			name:     "auto calculate with 0",
			workers:  0,
			expected: expectedAutoValue(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateOptimal(tt.workers)
			if result != tt.expected {
				t.Errorf("CalculateOptimal(%d) = %d, expected %d", tt.workers, result, tt.expected)
			}
		})
	}

	// Test bounds for auto-calculated values
	t.Run("auto-calculated value is within bounds", func(t *testing.T) {
		result := CalculateOptimal(0)
		if result < MinWorkers || result > MaxWorkers {
			t.Errorf("CalculateOptimal(0) = %d, expected between %d and %d", result, MinWorkers, MaxWorkers)
		}
	})
}

func expectedAutoValue() int {
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
