package firewallrules

import "errors"

// Used only by the v4 command tree.
var (
	ErrorUpdate            = errors.New("failed to update the Firewall Rule: %w")
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
	ErrorConvertRuleId     = errors.New("invalid Rule ID. The value must be an integer")
)
