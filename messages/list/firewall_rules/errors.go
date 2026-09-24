package firewallrules

import "errors"

// Used only by the v4 command tree.
var (
	ErrorGetAll            = "failed to list Firewall Rules: %w"
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
)
