package wafexceptions

import "errors"

// Used only by the v4 command tree.
var (
	ErrorCreate       = errors.New("failed to create the WAF Exception: %w")
	ErrorConvertWafID = errors.New("invalid WAF ID. The value must be an integer")
	ErrorIsActiveFlag = errors.New("invalid value for 'active' flag")
)
