package firewallrules

import "errors"

var (
	ErrorFailToDeleteRule  = errors.New("failed to delete the Firewall Rule: %w")
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
	ErrorConvertRuleId     = errors.New("invalid Rule ID. The value must be an integer")
)
