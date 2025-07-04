package warmup

var (
	Usage            = "warmup"
	ShortDescription = "Preload URLs into edge cache"
	LongDescription  = "Performs cache warming by systematically crawling a URL and preloading its content into the edge cache"
	FlagHelp         = "Displays more information about the warmup command"
	FlagUrl          = "Base URL to start warming the cache"
	FlagMaxUrls      = "Maximum number of URLs to process (default: 1500)"
	FlagMaxConcurrent = "Maximum number of concurrent requests (default: 2)"
	FlagTimeout      = "Timeout in milliseconds for each request (default: 8000)"
	WarmupSuccessful = "Cache warming completed successfully!"
	AskForUrl        = "Enter the base URL to warm cache (e.g., https://example.com):"
	
	// Log messages
	LogInitializing   = "\nInitializing cache warming"
	LogSite           = "Target site"
	LogConfiguration  = "Configuration"
	LogTimeout        = "Request timeout"
	LogBatchProcess   = "\nProcessing batch"
	LogAddedToQueue   = "\nAdded to queue"
	LogStatus         = "Status"
	LogCompleted      = "\nCOMPLETED"
	LogProcessed      = "URLs processed"
	LogSuccessful     = "Successful"
	LogFailed         = "Failed"
	LogTotalTime      = "Total time"
	LogSpeed          = "Speed"
	LogFailedUrls     = "Failed URLs"
	LogFoundLinks     = "\nFound links"
	LogLinkExamples   = "Examples"
) 