package firewallrules

import "errors"

// Used only by the v4 command tree.
var (
	ErrorGetFirewallRule   = "failed to describe the Firewall Rule: %s"
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
	ErrorConvertRuleId     = errors.New("invalid Rule ID. The value must be an integer")
)
