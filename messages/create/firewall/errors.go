package firewall

import "errors"

var (
	ErrorCreateFunction        = errors.New("Failed to create the Firewall: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorDebugFlag             = errors.New("Invalid value for --debug-rules flag")
	ErrorActiveFlag            = errors.New("Invalid value for --active flag")
	ErrorFunctionsEnabledFlag  = errors.New("Invalid value for --functions-enabled flag")
	ErrorNetworkProtectionFlag = errors.New("Invalid value for --network-protection flag")
	ErrorWafEnabledFlag        = errors.New("Invalid value for --waf-enabled flag")
)
