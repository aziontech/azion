package firewallinstance

import "errors"

// Used only by the v4 command tree.
var (
	ErrorFailToDeletInstance               = errors.New("failed to delete the Firewall Function Instance: %w")
	ErrorConvertFirewallId                 = errors.New("invalid Firewall ID. The value must be an integer")
	ErrorConvertFirewallFunctionInstanceId = errors.New("invalid Firewall Function Instance ID. The value must be an integer")
)
