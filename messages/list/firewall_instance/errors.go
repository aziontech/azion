package firewallinstance

import "errors"

var (
	ErrorGetAll            = "failed to list Firewall Function Instances: %w"
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
)
