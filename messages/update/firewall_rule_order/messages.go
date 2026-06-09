package firewall_rule_order

import "errors"

var (
	Usage              = "firewall-rule-order"
	ShortDescription   = "Orders the rules in Rules Engine of a Firewall"
	LongDescription    = "Defines the execution order of the rules in Rules Engine for a given Firewall by sending the ordered list of rule IDs"
	FlagFirewallID     = "Unique identifier for a Firewall"
	FlagRuleIDs        = "Comma-separated list of rule IDs in the desired execution order (e.g. 123,456,789)"
	OutputSuccess      = "Ordered Rules Engine of Firewall with ID %d"
	HelpFlag           = "Displays more information about the azion update firewall-rule-order subcommand"
	AskInputFirewallID = "Enter the ID of the Firewall whose rules will be ordered:"
	AskInputRuleIDs    = "Enter the rule IDs in the desired execution order (comma-separated):"

	ErrorConvertFirewallID = errors.New("Invalid --firewall-id flag provided. The value must be an integer. Run the command 'azion update firewall-rule-order --help' to display more information and try again")
	ErrorConvertRuleIDs    = errors.New("Invalid --rule-ids flag provided. The value must be a comma-separated list of integers. Run the command 'azion update firewall-rule-order --help' to display more information and try again")
	ErrorOrder             = errors.New("Failed to order the rules in Rules Engine of the Firewall: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
