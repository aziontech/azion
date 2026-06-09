package rules_engine_order

import "errors"

var (
	Usage                 = "rules-engine-order"
	ShortDescription      = "Orders the rules in Rules Engine of an Application"
	LongDescription       = "Defines the execution order of the rules in Rules Engine for a given Application and phase by sending the ordered list of rule IDs"
	FlagApplicationID     = "Unique identifier for an Application"
	FlagPhase             = "The phase is either 'request' or 'response'"
	FlagRuleIDs           = "Comma-separated list of rule IDs in the desired execution order (e.g. 123,456,789)"
	OutputSuccess         = "Ordered Rules Engine of Application with ID %d"
	HelpFlag              = "Displays more information about the azion update rules-engine-order subcommand"
	AskInputApplicationID = "Enter the ID of the Application whose Rules Engine will be ordered:"
	AskInputPhase         = "Enter the Rule Engine's phase (request/response):"
	AskInputRuleIDs       = "Enter the rule IDs in the desired execution order (comma-separated):"

	ErrorConvertApplicationID = errors.New("Invalid --application-id flag provided. The value must be an integer. Run the command 'azion update rules-engine-order --help' to display more information and try again")
	ErrorConvertRuleIDs       = errors.New("Invalid --rule-ids flag provided. The value must be a comma-separated list of integers. Run the command 'azion update rules-engine-order --help' to display more information and try again")
	ErrorInvalidPhase         = errors.New("Invalid phase value provided. The value must be 'request' or 'response'. Run the command 'azion update rules-engine-order --help' to display more information and try again")
	ErrorOrder                = errors.New("Failed to order the rules in Rules Engine: %s. Check your settings and try again. If the error persists, contact Azion support.")
)
