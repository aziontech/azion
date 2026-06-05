package crl

import "errors"

var (
	ErrorCreateCRL   = errors.New("Failed to create the Certificate Revocation List: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorReadCRLFile = errors.New("Failed to read the CRL file")
	ErrorActiveFlag  = errors.New("Invalid value for --active flag")
	ErrorInvalidJSON = errors.New("Failed to parse the input file: it does not match the Certificate Revocation List schema. Check the file contents and try again")
)
