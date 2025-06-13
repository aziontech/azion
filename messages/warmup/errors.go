package warmup

import "errors"

var (
	ErrorInvalidUrl        = errors.New("Invalid URL provided. URL must be a valid HTTP/HTTPS URL")
	ErrorMaxUrlsExceeded   = errors.New("Maximum number of URLs exceeded")
	ErrorRequestTimeout    = errors.New("Request timed out")
	ErrorProcessingFailed  = errors.New("Failed to process URL")
) 