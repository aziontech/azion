package digitalcertificate

import "errors"

// Used only by the v4 command tree.
var (
	ErrorGetDigitalCertificate       = "Failed to get the Digital Certificate: %s. Check your settings and try again. If the error persists, contact Azion support"
	ErrorConvertIdDigitalCertificate = errors.New("The Digital Certificate ID you provided is invalid. The value must be an integer. You may run the 'azion list digital-certificate' command to check your Digital Certificate ID")
)
