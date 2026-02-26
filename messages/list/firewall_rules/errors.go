package firewallrules

import "errors"

var (
	ErrorGetAll            = "failed to list Firewall Rules: %w"
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
)
