package waf

import "errors"

var (
	ErrorConvertId       = errors.New("The WAF ID you provided is invalid. The value must be an integer. You may run the 'azion list waf' command to check your WAF ID")
	ErrorFailToDeleteWaf = "Failed to delete the WAF: %s. Check your settings and try again. If the error persists, contact Azion support"
)
