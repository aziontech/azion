package csr

import "errors"

var (
	ErrorGetCSR       = "Failed to get the Certificate Signing Request: %s. Check your settings and try again. If the error persists, contact Azion support"
	ErrorConvertIdCSR = errors.New("The Certificate Signing Request ID you provided is invalid. The value must be an integer")
)
