package wafexceptions

import "errors"

var (
	ErrorGetWafException    = "failed to describe the WAF Exception: %s"
	ErrorConvertWafID       = errors.New("invalid WAF ID. The value must be an integer")
	ErrorConvertExceptionID = errors.New("invalid WAF Exception ID. The value must be an integer")
)
