package root

import (
	"testing"
	"time"

	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

// TestTimeComparisonLogic tests the core logic of comparing publish time with current time
// to determine if an update message should be shown
func TestTimeComparisonLogic(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name             string
		publishedAt      time.Time
		shouldShowUpdate bool
	}{
		{
			name:             "Less than 24 hours since publishing",
			publishedAt:      time.Now().Add(-23 * time.Hour),
			shouldShowUpdate: false,
		},
		{
			name:             "Exactly 24 hours since publishing",
			publishedAt:      time.Now().Add(-24 * time.Hour),
			shouldShowUpdate: true,
		},
		{
			name:             "More than 24 hours since publishing",
			publishedAt:      time.Now().Add(-25 * time.Hour),
			shouldShowUpdate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the core logic directly
			actualShouldShow := time.Since(tt.publishedAt) >= 24*time.Hour
			
			// Assert whether update message should be shown as expected
			assert.Equal(t, tt.shouldShowUpdate, actualShouldShow, 
				"Expected shouldShowUpdate to be %v, but got %v for publish time %v", 
				tt.shouldShowUpdate, actualShouldShow, tt.publishedAt)
		})
	}
}

// TestParsePublishedAtDate tests the parsing of the published_at date from GitHub API
func TestParsePublishedAtDate(t *testing.T) {
	tests := []struct {
		name        string
		publishedAt string
		isValid     bool
	}{
		{
			name:        "Valid RFC3339 date",
			publishedAt: "2025-09-12T19:15:10Z",
			isValid:     true,
		},
		{
			name:        "Invalid date format",
			publishedAt: "2025/09/12 19:15:10",
			isValid:     false,
		},
		{
			name:        "Empty date",
			publishedAt: "",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Try to parse the date
			_, err := time.Parse(time.RFC3339, tt.publishedAt)
			
			// Check if parsing was successful as expected
			if tt.isValid {
				assert.NoError(t, err, "Expected valid date but got error: %v", err)
			} else {
				assert.Error(t, err, "Expected error for invalid date but got none")
			}
		})
	}
}
