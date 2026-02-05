package firewall

import "errors"

var (
	ErrorUpdateFirewall        = errors.New("Failed to update the Firewall: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorDebugFlag             = errors.New("Invalid value for --debug-rules flag")
	ErrorActiveFlag            = errors.New("Invalid value for --active flag")
	ErrorFunctionsEnabledFlag  = errors.New("Invalid value for --functions-enabled flag")
	ErrorNetworkProtectionFlag = errors.New("Invalid value for --network-protection flag")
	ErrorWafEnabledFlag        = errors.New("Invalid value for --waf-enabled flag")
	ErrorConvertIdFirewall     = errors.New("The firewall ID you provided is invalid. The value must be an integer. You may run the 'azion list firewall' command to check your firewall ID")
)
