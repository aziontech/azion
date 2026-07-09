package crl

import "errors"

var (
	ErrorConvertId       = errors.New("The Certificate Revocation List ID you provided is invalid. The value must be an integer. You may run the 'azion list crl' command to check your Certificate Revocation List ID")
	ErrorFailToDeleteCRL = "Failed to delete the Certificate Revocation List: %s. Check your settings and try again. If the error persists, contact Azion support"
)
