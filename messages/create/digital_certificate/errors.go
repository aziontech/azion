package digitalcertificate

import "errors"

var (
	ErrorCreateDigitalCertificate  = errors.New("Failed to create the Digital Certificate: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorRequestDigitalCertificate = errors.New("Failed to request the Digital Certificate: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorCertificateTypeFlag       = errors.New("Invalid value for --certificate-type flag")
	ErrorReadCertificateFile       = errors.New("Failed to read the certificate file")
	ErrorReadPrivateKeyFile        = errors.New("Failed to read the private key file")
	ErrorInvalidJSONFile           = errors.New("Failed to parse the input file: it does not match either a Digital Certificate or a Certificate Request schema. Check the file contents and try again")
)
