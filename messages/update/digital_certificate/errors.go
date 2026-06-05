package digitalcertificate

import "errors"

var (
	ErrorUpdateDigitalCertificate    = errors.New("Failed to update the Digital Certificate: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorActiveFlag                  = errors.New("Invalid value for --active flag")
	ErrorReadCertificateFile         = errors.New("Failed to read the certificate file")
	ErrorReadPrivateKeyFile          = errors.New("Failed to read the private key file")
	ErrorConvertIdDigitalCertificate = errors.New("The Digital Certificate ID you provided is invalid. The value must be an integer. You may run the 'azion list digital-certificate' command to check your Digital Certificate ID")
)
