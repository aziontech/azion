package csr

import "errors"

var (
	ErrorConvertId       = errors.New("The Certificate Signing Request ID you provided is invalid. The value must be an integer")
	ErrorFailToDeleteCSR = "Failed to delete the Certificate Signing Request: %s. Check your settings and try again. If the error persists, contact Azion support"
)
