package firewallrules

import "errors"

var (
	ErrorGetFirewallRule   = "failed to describe the Firewall Rule: %s"
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
	ErrorConvertRuleId     = errors.New("invalid Rule ID. The value must be an integer")
)
