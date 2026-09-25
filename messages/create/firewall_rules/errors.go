package firewallrules

import "errors"

// Used only by the v4 command tree.
var (
	ErrorCreate            = errors.New("failed to create the Firewall Rule: %w")
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
)
