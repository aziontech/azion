package crl

import "errors"

var (
	ErrorUpdateCRL    = errors.New("Failed to update the Certificate Revocation List: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorActiveFlag   = errors.New("Invalid value for --active flag")
	ErrorReadCRLFile  = errors.New("Failed to read the CRL file")
	ErrorConvertIdCRL = errors.New("The Certificate Revocation List ID you provided is invalid. The value must be an integer. You may run the 'azion list crl' command to check your Certificate Revocation List ID")
)
