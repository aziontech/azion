package firewallinstance

import "errors"

var (
	ErrorGetFirewallFunctionInstance       = "failed to describe the Firewall Function Instance: %s"
	ErrorConvertFirewallId                 = errors.New("invalid Firewall ID. The value must be an integer")
	ErrorConvertFirewallFunctionInstanceId = errors.New("invalid Firewall Function Instance ID. The value must be an integer")
)
