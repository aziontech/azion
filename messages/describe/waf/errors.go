package waf

import "errors"

var (
	ErrorGetWaf       = "Failed to get the WAF: %s. Check your settings and try again. If the error persists, contact Azion support"
	ErrorConvertIdWaf = errors.New("The WAF ID you provided is invalid. The value must be an integer. You may run the 'azion list waf' command to check your WAF ID")
)
