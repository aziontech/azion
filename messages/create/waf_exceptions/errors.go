package wafexceptions

import "errors"

var (
	ErrorCreate         = errors.New("failed to create the WAF Exception: %w")
	ErrorConvertWafID   = errors.New("invalid WAF ID. The value must be an integer")
	ErrorConvertRuleID  = errors.New("invalid Rule ID. The value must be an integer")
	ErrorIsActiveFlag   = errors.New("invalid value for 'active' flag")
	ErrorConditionsFlag = errors.New("failed to parse conditions JSON")
)
