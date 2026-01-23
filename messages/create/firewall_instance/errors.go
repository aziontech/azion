package firewallinstance

import "errors"

var (
	ErrorCreate            = errors.New("failed to create the Firewall Function Instance: %w")
	ErrorConvertFirewallId = errors.New("invalid Firewall ID. The value must be an integer")
	ErrorConvertFunctionID = errors.New("invalid Function ID. The value must be an integer")
	ErrorIsActiveFlag      = errors.New("invalid value for 'active' flag")
	ErrorArgsFlag          = errors.New("failed to read args file")
	ErrorParseArgs         = errors.New("failed to parse args JSON")
)
