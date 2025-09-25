package functioninstance

var (
	Usage                 = "function-instance"
	ShortDescription      = "Creates a new Function Instance"
	LongDescription       = "Creates a Function Instance based on given attributes"
	FlagName              = "The Function Instance's name"
	FlagIsActive          = "Whether the Function Instance is active or not"
	FlagFile              = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess         = "Created Function Instance with ID %d"
	HelpFlag              = "Displays more information about the create function-instance command"
	AskInputName          = "Enter the new Function Instance's name:"
	AskInputApplicationID = "Enter the Application's ID this Function Instance will be associated with:"
	AskInputFunctionID    = "Enter the Function's ID:"
	FlagApplicationID     = "Unique identifier of the Application"
	FlagFunctionID        = "Unique identifier of the Function"
	FlagArgs              = "The Function Instance's arguments"
)
