package waf

import "errors"

var (
	ErrorCreateWAF            = errors.New("Failed to create the WAF: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorActiveFlag           = errors.New("Invalid value for --active flag")
	ErrorEngineVersionFlag    = errors.New("Invalid value for --engine-version flag")
	ErrorRulesetsFlag         = errors.New("Invalid value for --rulesets flag. Must be comma-separated integers")
	ErrorThresholdsFlag       = errors.New("Invalid value for --thresholds flag. Format: threat1=sensitivity1,threat2=sensitivity2")
	ErrorThresholdThreat      = errors.New("Invalid threat value for threshold. Valid values: cross_site_scripting, directory_traversal, evading_tricks, file_upload, identified_attack, remote_file_inclusion, sql_injection, unwanted_access")
	ErrorThresholdSensitivity = errors.New("Invalid sensitivity value for threshold. Valid values: highest, high, medium, low, lowest")
)
