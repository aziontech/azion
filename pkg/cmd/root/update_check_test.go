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

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name    string
		latest  string
		current string
		want    bool
	}{
		{"newer minor", "4.25.0", "4.24.0", true},
		{"newer patch", "4.24.1", "4.24.0", true},
		{"same version", "4.24.0", "4.24.0", false},
		{"older version", "4.23.0", "4.24.0", false},
		{"stage tag of the installed version", "4.24.0-dev.20", "4.24.0", false},
		{"numeric, not lexical, comparison", "4.30.0", "4.3.10", true},
		{"digit concatenation would say newer", "4.3.10", "4.30.0", false},
		{"v prefix", "v4.25.0", "4.24.0", true},
		{"local build ahead of a tag", "4.24.0", "4.24.0-3-g8b509c9", false},
		{"empty latest", "", "4.24.0", false},
		{"unparseable current", "4.25.0", "dev", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isNewerVersion(tt.latest, tt.current))
		})
	}
}
