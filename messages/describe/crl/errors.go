package crl

import "errors"

// Used only by the v4 command tree.
var (
	ErrorGetCRL       = "Failed to get the Certificate Revocation List: %s. Check your settings and try again. If the error persists, contact Azion support"
	ErrorConvertIdCRL = errors.New("The Certificate Revocation List ID you provided is invalid. The value must be an integer. You may run the 'azion list crl' command to check your Certificate Revocation List ID")
)
