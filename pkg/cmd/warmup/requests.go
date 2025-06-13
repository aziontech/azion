package warmup

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2"
	msg "github.com/aziontech/azion-cli/messages/warmup"
	"github.com/aziontech/azion-cli/pkg/cmdutil"
	"github.com/aziontech/azion-cli/pkg/logger"
	"go.uber.org/zap"
)


type urlCache struct {
	sync.RWMutex
	visited map[string]bool
	failed  map[string]bool
}

func newUrlCache() *urlCache {
	return &urlCache{
		visited: make(map[string]bool),
		failed:  make(map[string]bool),
	}
}

func (c *urlCache) isVisited(url string) bool {
	c.RLock()
	defer c.RUnlock()
	return c.visited[url]
}

func (c *urlCache) markVisited(url string) {
	c.Lock()
	defer c.Unlock()
	c.visited[url] = true
}

func (c *urlCache) markFailed(url string) {
	c.Lock()
	defer c.Unlock()
	c.failed[url] = true
}

func (c *urlCache) failedCount() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.failed)
}

func (c *urlCache) failedURLs() []string {
	c.RLock()
	defer c.RUnlock()
	urls := make([]string, 0, len(c.failed))
	for url := range c.failed {
		urls = append(urls, url)
	}
	return urls
}
var blacklist = []string{
	".pdf", ".zip", ".rar", ".7z", ".tar", ".gz",
	".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm",
	".mp3", ".wav", ".ogg", ".m4a", ".aac",
	"mailto:", "tel:", "sms:", "ftp:", "file:",
	"javascript:", "data:", "#",
}


func formatLog(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}


func warmupCache(ctx context.Context, baseUrl string, maxUrls int, maxConcurrent int, timeout int, f *cmdutil.Factory) error {

	_, err := url.Parse(baseUrl)
	if err != nil {
		return msg.ErrorInvalidUrl
	}


	cache := newUrlCache()
	pendingUrls := make([]string, 0)
	pendingUrls = append(pendingUrls, baseUrl)

	totalProcessed := 0
	startTime := time.Now()

	logger.FInfo(f.IOStreams.Out, "\nInitializing cache warming...\n")
	logger.FInfo(f.IOStreams.Out, formatLog("Target site: %s\n", baseUrl))
	logger.FInfo(f.IOStreams.Out, formatLog("Configuration: %d concurrent requests, max %d URLs\n", maxConcurrent, maxUrls))
	logger.FInfo(f.IOStreams.Out, formatLog("Request timeout: %dms\n", timeout))
	logger.FInfo(f.IOStreams.Out, "\n")

	for len(pendingUrls) > 0 && totalProcessed < maxUrls {
		batchSize := min(len(pendingUrls), maxConcurrent)
		currentBatch := pendingUrls[:batchSize]
		pendingUrls = pendingUrls[batchSize:]

		logger.FInfo(f.IOStreams.Out, formatLog("Processing batch: %d URLs (%d remaining in queue)\n", batchSize, len(pendingUrls)))

		var wg sync.WaitGroup
		var mutex sync.Mutex
		var logMutex sync.Mutex
		newLinks := 0

		for _, currentUrl := range currentBatch {
			wg.Add(1)

			go func(url string) {
				defer wg.Done()

				links, err := processURL(url, timeout, cache, totalProcessed, maxUrls, baseUrl, f.IOStreams.Out, &logMutex)
				if err != nil {
					logger.Debug("Error processing URL", zap.String("url", url), zap.Error(err))
					cache.markFailed(url)
					return
				}

				var newValidLinks []string
				for _, link := range links {
					if !cache.isVisited(link) && !contains(pendingUrls, link) {
						newValidLinks = append(newValidLinks, link)
					}
				}

				if len(newValidLinks) > 0 {
					mutex.Lock()
					pendingUrls = append(pendingUrls, newValidLinks...)
					newLinks += len(newValidLinks)
					mutex.Unlock()
				}
			}(currentUrl)
		}

		wg.Wait()
		totalProcessed += batchSize

		if newLinks > 0 {
			logger.FInfo(f.IOStreams.Out, formatLog("Added to queue: %d new URLs\n", newLinks))
		}

		logger.FInfo(f.IOStreams.Out, formatLog("Status: %d processed | %d queued | %d failed\n\n", totalProcessed, len(pendingUrls), cache.failedCount()))

		time.Sleep(200 * time.Millisecond)
	}
	duration := time.Since(startTime)
	logger.FInfo(f.IOStreams.Out, "COMPLETED!\n")
	logger.FInfo(f.IOStreams.Out, formatLog("Processed: %d URLs\n", totalProcessed))
	logger.FInfo(f.IOStreams.Out, formatLog("Successful: %d\n", totalProcessed-cache.failedCount()))
	logger.FInfo(f.IOStreams.Out, formatLog("Failed: %d\n", cache.failedCount()))
	logger.FInfo(f.IOStreams.Out, formatLog("Total time: %.1fs\n", duration.Seconds()))
	
	if totalProcessed > 0 {
		logger.FInfo(f.IOStreams.Out, formatLog("Speed: %.1f URLs/s\n", float64(totalProcessed)/duration.Seconds()))
	}


	failedURLs := cache.failedURLs()
	if len(failedURLs) > 0 && len(failedURLs) <= 10 {
		logger.FInfo(f.IOStreams.Out, "\nFailed URLs:\n")
		for _, u := range failedURLs {
			logger.FInfo(f.IOStreams.Out, formatLog("  - %s\n", formatURL(u, baseUrl)))
		}
	}

	return nil
}

func processURL(currentURL string, timeoutMs int, cache *urlCache, processed int, maxUrls int, baseURL string, out io.Writer, logMutex *sync.Mutex) ([]string, error) {
	if cache.isVisited(currentURL) || processed >= maxUrls {
		return nil, nil
	}

	cache.markVisited(currentURL)

	shortURL := formatURL(currentURL, baseURL)
	
	logMutex.Lock()
	logger.FInfo(out, formatLog("[%d/%d] %s\n", processed+1, maxUrls, shortURL))
	logMutex.Unlock()

	client := &http.Client{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}

	req, err := http.NewRequest(http.MethodGet, currentURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	links := extractLinks(string(body), baseURL, out, logMutex)
	return links, nil
}

func extractLinks(html, baseURL string, out io.Writer, logMutex *sync.Mutex) []string {
	foundLinks := make(map[string]bool)
	
	linkRegex := regexp.MustCompile(`<a[^>]+href=["']([^"']+)["'][^>]*>`)
	matches := linkRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	formRegex := regexp.MustCompile(`<form[^>]+action=["']([^"']+)["'][^>]*>`)
	matches = formRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	cssRegex := regexp.MustCompile(`<link[^>]+href=["']([^"']+\.css[^"']*)["'][^>]*>`)
	matches = cssRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	jsRegex := regexp.MustCompile(`<script[^>]+src=["']([^"']+\.js[^"']*)["'][^>]*>`)
	matches = jsRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	metaRegex := regexp.MustCompile(`<meta[^>]+http-equiv=["']refresh["'][^>]+content=["'][^;]*;\s*url=([^"']+)["'][^>]*>`)
	matches = metaRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}

	canonicalRegex := regexp.MustCompile(`<link[^>]+rel=["']canonical["'][^>]+href=["']([^"']+)["'][^>]*>`)
	matches = canonicalRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	imgRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches = imgRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	iconRegex := regexp.MustCompile(`<link[^>]+rel=["'](?:icon|shortcut icon|apple-touch-icon)["'][^>]+href=["']([^"']+)["'][^>]*>`)
	matches = iconRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	fontRegex := regexp.MustCompile(`@font-face[^}]*url\(["']?([^"')]+)["']?\)`)
	matches = fontRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	bgRegex := regexp.MustCompile(`background(?:-image)?:\s*url\(["']?([^"')]+)["']?\)`)
	matches = bgRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	links := make([]string, 0, len(foundLinks))
	for link := range foundLinks {
		links = append(links, link)
	}
	
	if len(links) > 0 {
		logMutex.Lock()
		logger.FInfo(out, formatLog("  Found links: %d\n", len(links)))
		
		examples := min(5, len(links))
		if examples > 0 {
			logger.FInfo(out, "  Examples:\n")
			for i := 0; i < examples; i++ {
				logger.FInfo(out, formatLog("    - %s\n", formatURL(links[i], baseURL)))
			}
		}
		logMutex.Unlock()
	}
	
	return links
}

func processFoundLink(link, baseURL string, foundLinks map[string]bool) {
	if isBlacklisted(link) {
		return
	}
	
	normalizedLink := normalizeURL(link, baseURL)
	if normalizedLink != "" {
		foundLinks[normalizedLink] = true
	}
}

func isBlacklisted(url string) bool {
	urlLower := strings.ToLower(url)
	for _, pattern := range blacklist {
		if strings.Contains(urlLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func formatURL(fullURL string, baseURL string) string {
	result := strings.Replace(fullURL, baseURL, "", 1)
	if result == "" {
		return "/"
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func askForUrl() (string, error) {
	var baseUrl string
	prompt := &survey.Input{
		Message: msg.AskForUrl,
	}

	err := survey.AskOne(prompt, &baseUrl, survey.WithKeepFilter(true), survey.WithStdio(os.Stdin, os.Stderr, os.Stdout))
	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(baseUrl, "http://") && !strings.HasPrefix(baseUrl, "https://") {
		baseUrl = "https://" + baseUrl
	}

	return baseUrl, nil
}

func normalizeURL(link, baseURL string) string {
	link = strings.TrimSpace(link)
	if link == "" || link == "#" {
		return ""
	}
	
	if strings.Contains(link, "#") {
		link = strings.Split(link, "#")[0]
		if link == "" {
			return ""
		}
	}
	
	baseURLParsed, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	
	// Use url.ResolveReference for proper URL resolution
	linkURL, err := url.Parse(link)
	if err != nil {
		return ""
	}
	
	resolvedURL := baseURLParsed.ResolveReference(linkURL)
	
	// Only allow same-host URLs
	if resolvedURL.Host != baseURLParsed.Host {
		return ""
	}
	
	// Clean up the URL
	resolvedURL.Fragment = ""
	
	// Keep only important query parameters for cache warming
	paramsToKeep := []string{"id", "category", "product", "page", "search", "q", "query", "filter", "sort", "lang", "locale"}
	query := resolvedURL.Query()
	newQuery := make(url.Values)
	
	for _, param := range paramsToKeep {
		if value := query.Get(param); value != "" {
			newQuery.Set(param, value)
		}
	}
	
	resolvedURL.RawQuery = newQuery.Encode()
	
	return resolvedURL.String()
} 