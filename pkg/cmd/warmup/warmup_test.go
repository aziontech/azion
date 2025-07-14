package warmup

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/httpmock"
	"github.com/aziontech/azion-cli/pkg/logger"
	"github.com/aziontech/azion-cli/pkg/testutils"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewCmd(t *testing.T) {
	logger.New(zapcore.DebugLevel)
	mock := &httpmock.Registry{}
	f, _, _ := testutils.NewFactory(mock)
	cmd := NewCmd(f)

	assert.Equal(t, cmd.Use, "warmup")
	assert.Equal(t, cmd.Short, "Preload URLs into edge cache")
	assert.NotEmpty(t, cmd.Example)
}

func TestWarmupCmd_Run(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tests := []struct {
		name      string
		mock      func(*WarmupCmd)
		wantErr   bool
		wantErrIs error
	}{
		{
			name: "successfully warmup cache",
			mock: func(w *WarmupCmd) {
				w.WarmupCache = func(ctx context.Context, baseUrl string, maxUrls int, maxConcurrent int, timeout int, f *cmdutil.Factory) error {
					return nil
				}
				w.AskForUrl = func() (string, error) {
					return "https://example.com", nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &httpmock.Registry{}
			f, stdout, _ := testutils.NewFactory(mock)

			cmd := &cobra.Command{}
			baseUrl = "https://example.com"

			warmupCmd := &WarmupCmd{
				Io: f.IOStreams,
			}

			tt.mock(warmupCmd)

			err := warmupCmd.Run(context.Background(), cmd, f)
			if tt.wantErr {
				assert.Error(t, err)
				if tt.wantErrIs != nil {
					assert.ErrorIs(t, err, tt.wantErrIs)
				}
			} else {
				assert.NoError(t, err)
				assert.Contains(t, stdout.String(), "Cache warming completed successfully!")
			}
		})
	}
}

func TestWarmupCache(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.html")
	err := os.WriteFile(testFile, []byte(`<html><body><a href="/page1">Page 1</a></body></html>`), 0644)
	assert.NoError(t, err)

	mock := &httpmock.Registry{}
	f, _, _ := testutils.NewFactory(mock)

	err = warmupCache(context.Background(), "https://example.com", 10, 2, 30, f)
	assert.NoError(t, err)
}

func TestExtractLinks(t *testing.T) {
	logger.New(zapcore.DebugLevel)

	html := `<html><head>
		<link rel="stylesheet" href="/styles/main.css">
		<script src="/js/app.js"></script>
		<link rel="canonical" href="https://example.com/canonical-page">
		<meta http-equiv="refresh" content="5; url=/redirect-page">
		<link rel="icon" href="/favicon.ico">
		<style>
			.hero { background-image: url('/images/hero-bg.jpg'); }
		</style>
	</head><body>
		<a href="/page1">Page 1</a>
		<a href="https://example.com/page2">Page 2</a>
		<a href="mailto:test@example.com">Email</a>
		<a href="/test.pdf">PDF</a>
		<a href="./relative-page">Relative</a>
		<a href="../parent-page">Parent</a>
		<img src="/images/logo.png" alt="Logo">
		<img src="https://example.com/images/banner.jpg" alt="Banner">
		<form action="/search" method="post">
			<input type="submit" value="Search">
		</form>
	</body></html>`

	mock := &httpmock.Registry{}
	f, _, _ := testutils.NewFactory(mock)

	var logMutex sync.Mutex
	links := extractLinks(html, "https://example.com", f.IOStreams.Out, &logMutex)

	// Verify traditional links
	assert.Contains(t, links, "https://example.com/page1")
	assert.Contains(t, links, "https://example.com/page2")

	// Verify CSS/JS resources
	assert.Contains(t, links, "https://example.com/styles/main.css")
	assert.Contains(t, links, "https://example.com/js/app.js")

	// Verify forms
	assert.Contains(t, links, "https://example.com/search")

	// Verify canonical
	assert.Contains(t, links, "https://example.com/canonical-page")

	// Verify meta refresh
	assert.Contains(t, links, "https://example.com/redirect-page")

	// Verify relative URLs
	assert.Contains(t, links, "https://example.com/relative-page")
	assert.Contains(t, links, "https://example.com/parent-page")

	// Verify images (new functionality)
	assert.Contains(t, links, "https://example.com/images/logo.png")
	assert.Contains(t, links, "https://example.com/images/banner.jpg")

	// Verify icons
	assert.Contains(t, links, "https://example.com/favicon.ico")

	// Verify background images
	assert.Contains(t, links, "https://example.com/images/hero-bg.jpg")

	// Verify blacklist still works
	assert.NotContains(t, links, "mailto:test@example.com")
	assert.NotContains(t, links, "https://example.com/test.pdf")
}

func TestIsBlacklisted(t *testing.T) {
	tests := []struct {
		url         string
		blacklisted bool
	}{
		// Normal pages - should pass
		{"https://example.com/page.html", false},
		{"https://example.com/normal-page", false},

		// Images - should now pass (no longer in blacklist)
		{"https://example.com/image.jpg", false},
		{"https://example.com/photo.png", false},
		{"https://example.com/icon.svg", false},

		// CSS and JS - should pass
		{"https://example.com/style.css", false},
		{"https://example.com/script.js", false},

		// Documents - should be blocked
		{"https://example.com/file.pdf", true},
		{"https://example.com/doc.docx", true},

		// Videos/audios - should be blocked
		{"https://example.com/video.mp4", true},
		{"https://example.com/audio.mp3", true},

		// Special links - should be blocked
		{"mailto:test@example.com", true},
		{"tel:+1234567890", true},
		{"javascript:void(0)", true},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := isBlacklisted(tt.url)
			assert.Equal(t, tt.blacklisted, result)
		})
	}
}

func TestFormatURL(t *testing.T) {
	tests := []struct {
		fullURL  string
		baseURL  string
		expected string
	}{
		{"https://example.com/page1", "https://example.com", "/page1"},
		{"https://example.com/", "https://example.com", "/"},
		{"https://other.com/page", "https://example.com", "https://other.com/page"},
	}

	for _, tt := range tests {
		t.Run(tt.fullURL, func(t *testing.T) {
			result := formatURL(tt.fullURL, tt.baseURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		link     string
		baseURL  string
		expected string
	}{
		// Absolute URLs
		{"/page1", "https://example.com", "https://example.com/page1"},
		{"https://example.com/page2", "https://example.com", "https://example.com/page2"},

		// Relative URLs
		{"./page3", "https://example.com/dir/", "https://example.com/dir/page3"},
		{"../page4", "https://example.com/dir/subdir/", "https://example.com/dir/page4"},
		{"page5.html", "https://example.com/dir/", "https://example.com/dir/page5.html"},

		// Protocol-relative
		{"//example.com/page6", "https://example.com", "https://example.com/page6"},

		// URLs with anchors (should be removed)
		{"/page7#section", "https://example.com", "https://example.com/page7"},

		// External URLs (should be ignored)
		{"https://other.com/page", "https://example.com", ""},
		{"//other.com/page", "https://example.com", ""},

		// Empty or invalid URLs
		{"", "https://example.com", ""},
		{"#", "https://example.com", ""},
		{"#section", "https://example.com", ""},
	}

	for _, tt := range tests {
		t.Run(tt.link, func(t *testing.T) {
			result := normalizeURL(tt.link, tt.baseURL)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessFoundLink(t *testing.T) {
	foundLinks := make(map[string]bool)

	// Valid link
	processFoundLink("/valid-page", "https://example.com", foundLinks)
	assert.Contains(t, foundLinks, "https://example.com/valid-page")

	// Blacklisted link
	processFoundLink("/document.pdf", "https://example.com", foundLinks)
	assert.NotContains(t, foundLinks, "https://example.com/document.pdf")

	// External link
	processFoundLink("https://other.com/page", "https://example.com", foundLinks)
	assert.NotContains(t, foundLinks, "https://other.com/page")
}
