package firewallinstance

import "errors"

var (
	ErrorUpdate                            = errors.New("failed to update the Firewall Function Instance: %w")
	ErrorConvertFirewallId                 = errors.New("invalid Firewall ID. The value must be an integer")
	ErrorConvertFirewallFunctionInstanceId = errors.New("invalid Firewall Function Instance ID. The value must be an integer")
	ErrorIsActiveFlag                      = errors.New("invalid value for 'active' flag")
	ErrorArgsFlag                          = errors.New("failed to read args file")
	ErrorParseArgs                         = errors.New("failed to parse args JSON")
)
