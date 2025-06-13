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

// Cache para URLs já visitadas e que falharam
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

// Função auxiliar para formatar mensagens de log
func formatLog(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Realiza o warmup do cache a partir da URL base
func warmupCache(ctx context.Context, baseUrl string, maxUrls int, maxConcurrent int, timeout int, f *cmdutil.Factory) error {
	// Validar URL
	_, err := url.Parse(baseUrl)
	if err != nil {
		return msg.ErrorInvalidUrl
	}

	// Inicializar cache de URLs
	cache := newUrlCache()
	pendingUrls := make([]string, 0)
	pendingUrls = append(pendingUrls, baseUrl)

	// Estatísticas
	totalProcessed := 0
	startTime := time.Now()

	// Cabeçalho inicial
	logger.FInfo(f.IOStreams.Out, "\nInitializing cache warming...\n")
	logger.FInfo(f.IOStreams.Out, formatLog("Target site: %s\n", baseUrl))
	logger.FInfo(f.IOStreams.Out, formatLog("Configuration: %d concurrent requests, max %d URLs\n", maxConcurrent, maxUrls))
	logger.FInfo(f.IOStreams.Out, formatLog("Request timeout: %dms\n", timeout))
	logger.FInfo(f.IOStreams.Out, "\n")

	// Iniciar o processamento em lotes
	for len(pendingUrls) > 0 && totalProcessed < maxUrls {
		// Limita o tamanho do lote ao número de URLs pendentes ou ao máximo de concorrência
		batchSize := min(len(pendingUrls), maxConcurrent)
		currentBatch := pendingUrls[:batchSize]
		pendingUrls = pendingUrls[batchSize:]

		logger.FInfo(f.IOStreams.Out, formatLog("Processing batch: %d URLs (%d remaining in queue)\n", batchSize, len(pendingUrls)))

		// Criar wait group para sincronizar o processamento em paralelo
		var wg sync.WaitGroup
		var mutex sync.Mutex // Para proteger a lista de pendingUrls durante atualizações concorrentes
		var logMutex sync.Mutex // Para sincronizar as mensagens de log
		newLinks := 0

		// Processar cada URL do lote em paralelo
		for _, currentUrl := range currentBatch {
			wg.Add(1)

			go func(url string) {
				defer wg.Done()

				// Processar URL e obter novos links
				links, err := processURL(url, timeout, cache, totalProcessed, maxUrls, baseUrl, f.IOStreams.Out, &logMutex)
				if err != nil {
					logger.Debug("Error processing URL", zap.String("url", url), zap.Error(err))
					cache.markFailed(url)
					return
				}

				// Filtrar apenas os links não visitados
				var newValidLinks []string
				for _, link := range links {
					if !cache.isVisited(link) && !contains(pendingUrls, link) {
						newValidLinks = append(newValidLinks, link)
					}
				}

				// Adicionar novos links à lista de pendentes
				if len(newValidLinks) > 0 {
					mutex.Lock()
					pendingUrls = append(pendingUrls, newValidLinks...)
					newLinks += len(newValidLinks)
					mutex.Unlock()
				}
			}(currentUrl)
		}

		// Aguardar todas as goroutines terminarem
		wg.Wait()
		totalProcessed += batchSize

		if newLinks > 0 {
			logger.FInfo(f.IOStreams.Out, formatLog("Added to queue: %d new URLs\n", newLinks))
		}

		// Status após cada lote
		logger.FInfo(f.IOStreams.Out, formatLog("Status: %d processed | %d queued | %d failed\n\n", totalProcessed, len(pendingUrls), cache.failedCount()))

		// Pequeno delay entre lotes para não sobrecarregar
		time.Sleep(200 * time.Millisecond)
	}

	// Estatísticas finais
	duration := time.Since(startTime)
	logger.FInfo(f.IOStreams.Out, "COMPLETED!\n")
	logger.FInfo(f.IOStreams.Out, formatLog("Processed: %d URLs\n", totalProcessed))
	logger.FInfo(f.IOStreams.Out, formatLog("Successful: %d\n", totalProcessed-cache.failedCount()))
	logger.FInfo(f.IOStreams.Out, formatLog("Failed: %d\n", cache.failedCount()))
	logger.FInfo(f.IOStreams.Out, formatLog("Total time: %.1fs\n", duration.Seconds()))
	
	if totalProcessed > 0 {
		logger.FInfo(f.IOStreams.Out, formatLog("Speed: %.1f URLs/s\n", float64(totalProcessed)/duration.Seconds()))
	}

	// Mostrar URLs que falharam (limitado a 10)
	failedURLs := cache.failedURLs()
	if len(failedURLs) > 0 && len(failedURLs) <= 10 {
		logger.FInfo(f.IOStreams.Out, "\nFailed URLs:\n")
		for _, u := range failedURLs {
			logger.FInfo(f.IOStreams.Out, formatLog("  - %s\n", formatURL(u, baseUrl)))
		}
	}

	return nil
}

// Processar uma URL e extrair links
func processURL(currentURL string, timeoutMs int, cache *urlCache, processed int, maxUrls int, baseURL string, out io.Writer, logMutex *sync.Mutex) ([]string, error) {
	if cache.isVisited(currentURL) || processed >= maxUrls {
		return nil, nil
	}

	cache.markVisited(currentURL)

	shortURL := formatURL(currentURL, baseURL)
	
	// Sincronizar a saída de log
	logMutex.Lock()
	logger.FInfo(out, formatLog("[%d/%d] %s\n", processed+1, maxUrls, shortURL))
	logMutex.Unlock()

	// Fazer requisição com timeout
	client := &http.Client{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}

	req, err := http.NewRequest(http.MethodGet, currentURL, nil)
	if err != nil {
		return nil, err
	}

	// Adicionar headers
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

	// Ler o conteúdo da página
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extrair e processar links
	links := extractLinks(string(body), baseURL, out, logMutex)
	return links, nil
}

// Extrai links de uma página HTML usando regex mais abrangente
func extractLinks(html, baseURL string, out io.Writer, logMutex *sync.Mutex) []string {
	foundLinks := make(map[string]bool)
	
	// 1. Links tradicionais <a href>
	linkRegex := regexp.MustCompile(`<a[^>]+href=["']([^"']+)["'][^>]*>`)
	matches := linkRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	// 2. Formulários <form action>
	formRegex := regexp.MustCompile(`<form[^>]+action=["']([^"']+)["'][^>]*>`)
	matches = formRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	// 3. Recursos CSS e JS (importantes para cache)
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
	
	// 4. Meta refresh redirects
	metaRegex := regexp.MustCompile(`<meta[^>]+http-equiv=["']refresh["'][^>]+content=["'][^;]*;\s*url=([^"']+)["'][^>]*>`)
	matches = metaRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	// 5. Canonical URLs
	canonicalRegex := regexp.MustCompile(`<link[^>]+rel=["']canonical["'][^>]+href=["']([^"']+)["'][^>]*>`)
	matches = canonicalRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	// 6. Imagens (importantes para cache!)
	imgRegex := regexp.MustCompile(`<img[^>]+src=["']([^"']+)["'][^>]*>`)
	matches = imgRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	// 7. Favicons e ícones
	iconRegex := regexp.MustCompile(`<link[^>]+rel=["'](?:icon|shortcut icon|apple-touch-icon)["'][^>]+href=["']([^"']+)["'][^>]*>`)
	matches = iconRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	// 8. Fontes web
	fontRegex := regexp.MustCompile(`@font-face[^}]*url\(["']?([^"')]+)["']?\)`)
	matches = fontRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	// 9. Background images em CSS inline
	bgRegex := regexp.MustCompile(`background(?:-image)?:\s*url\(["']?([^"')]+)["']?\)`)
	matches = bgRegex.FindAllStringSubmatch(html, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			processFoundLink(match[1], baseURL, foundLinks)
		}
	}
	
	// Converter map para slice
	links := make([]string, 0, len(foundLinks))
	for link := range foundLinks {
		links = append(links, link)
	}
	
	// Log dos links encontrados (sincronizado)
	if len(links) > 0 {
		logMutex.Lock()
		logger.FInfo(out, formatLog("  Found links: %d\n", len(links)))
		
		// Mostrar exemplos (até 5, um por linha, sem emoji)
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

// Processa um link encontrado e adiciona ao mapa se válido
func processFoundLink(link, baseURL string, foundLinks map[string]bool) {
	// Pular links blacklistados
	if isBlacklisted(link) {
		return
	}
	
	// Normalizar URL
	normalizedLink := normalizeURL(link, baseURL)
	if normalizedLink != "" {
		foundLinks[normalizedLink] = true
	}
}

// Verifica se uma URL está na blacklist
func isBlacklisted(url string) bool {
	urlLower := strings.ToLower(url)
	for _, pattern := range blacklist {
		if strings.Contains(urlLower, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

// Formata uma URL para exibição
func formatURL(fullURL string, baseURL string) string {
	result := strings.Replace(fullURL, baseURL, "", 1)
	if result == "" {
		return "/"
	}
	return result
}

// Verifica se uma slice contém um elemento
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Função min simples (Go 1.21+)
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Solicita a URL base ao usuário
func askForUrl() (string, error) {
	var baseUrl string
	prompt := &survey.Input{
		Message: msg.AskForUrl,
	}

	err := survey.AskOne(prompt, &baseUrl, survey.WithKeepFilter(true), survey.WithStdio(os.Stdin, os.Stderr, os.Stdout))
	if err != nil {
		return "", err
	}

	// Garantir que a URL começa com http:// ou https://
	if !strings.HasPrefix(baseUrl, "http://") && !strings.HasPrefix(baseUrl, "https://") {
		baseUrl = "https://" + baseUrl
	}

	return baseUrl, nil
}

// Normaliza uma URL relativa para absoluta e remove parâmetros desnecessários (versão melhorada)
func normalizeURL(link, baseURL string) string {
	// Pular links vazios ou inválidos
	link = strings.TrimSpace(link)
	if link == "" || link == "#" {
		return ""
	}
	
	// Remover âncoras no final
	if strings.Contains(link, "#") {
		link = strings.Split(link, "#")[0]
		if link == "" {
			return ""
		}
	}
	
	var fullURL string
	
	// Parse da URL base para ter contexto
	baseURLParsed, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}
	
	// Converter para URL absoluta
	if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
		// URL absoluta - verificar se é do mesmo domínio
		linkParsed, err := url.Parse(link)
		if err != nil {
			return ""
		}
		if linkParsed.Host != baseURLParsed.Host {
			return "" // URL externa, ignorar
		}
		fullURL = link
	} else if strings.HasPrefix(link, "//") {
		// Protocol-relative URL
		fullURL = baseURLParsed.Scheme + ":" + link
		linkParsed, err := url.Parse(fullURL)
		if err != nil {
			return ""
		}
		if linkParsed.Host != baseURLParsed.Host {
			return "" // URL externa, ignorar
		}
	} else if strings.HasPrefix(link, "/") {
		// URL absoluta no mesmo domínio
		fullURL = baseURLParsed.Scheme + "://" + baseURLParsed.Host + link
	} else if strings.HasPrefix(link, "./") {
		// URL relativa atual
		link = strings.TrimPrefix(link, "./")
		basePath := strings.TrimSuffix(baseURLParsed.Path, "/")
		fullURL = baseURLParsed.Scheme + "://" + baseURLParsed.Host + basePath + "/" + link
	} else if strings.HasPrefix(link, "../") {
		// URL relativa pai
		basePath := baseURLParsed.Path
		for strings.HasPrefix(link, "../") {
			link = strings.TrimPrefix(link, "../")
			// Subir um nível no path
			if basePath == "/" || basePath == "" {
				basePath = "/"
			} else {
				parts := strings.Split(strings.Trim(basePath, "/"), "/")
				if len(parts) > 1 {
					basePath = "/" + strings.Join(parts[:len(parts)-1], "/")
				} else {
					basePath = "/"
				}
			}
		}
		if basePath == "/" {
			fullURL = baseURLParsed.Scheme + "://" + baseURLParsed.Host + "/" + link
		} else {
			fullURL = baseURLParsed.Scheme + "://" + baseURLParsed.Host + basePath + "/" + link
		}
	} else {
		// URL relativa simples (page.html, etc.)
		basePath := strings.TrimSuffix(baseURLParsed.Path, "/")
		if basePath == "" {
			basePath = "/"
		}
		// Se o basePath não termina com /, adicionar
		if !strings.HasSuffix(basePath, "/") {
			// Se basePath tem extensão, pegar só o diretório
			if strings.Contains(basePath, ".") {
				parts := strings.Split(basePath, "/")
				if len(parts) > 1 {
					basePath = strings.Join(parts[:len(parts)-1], "/")
				} else {
					basePath = "/"
				}
			}
			if basePath != "/" {
				basePath += "/"
			}
		}
		fullURL = baseURLParsed.Scheme + "://" + baseURLParsed.Host + basePath + link
	}
	
	// Parse da URL final para validação e limpeza
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return ""
	}
	
	// Verificar se o host é válido
	if parsedURL.Host == "" {
		return ""
	}
	
	// Remover fragmento
	parsedURL.Fragment = ""
	
	// Manter parâmetros importantes para cache warming
	paramsToKeep := []string{"id", "category", "product", "page", "search", "q", "query", "filter", "sort", "lang", "locale"}
	query := parsedURL.Query()
	newQuery := make(url.Values)
	
	for _, param := range paramsToKeep {
		if value := query.Get(param); value != "" {
			newQuery.Set(param, value)
		}
	}
	
	parsedURL.RawQuery = newQuery.Encode()
	
	// Retornar URL limpa
	result := parsedURL.String()
	
	// Verificação adicional para evitar URLs malformadas
	if strings.Contains(result, "xn--") || strings.Contains(result, "@") {
		return ""
	}
	
	return result
} 