package csr

import "errors"

var (
	ErrorCreateCSR   = errors.New("Failed to create the Certificate Signing Request: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorInvalidJSON = errors.New("Failed to parse the input file: it does not match the Certificate Signing Request schema. Check the file contents and try again")
)
