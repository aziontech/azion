package functioninstance

var (
	Usage                 = "function-instance"
	ShortDescription      = "Updates a Function Instance"
	LongDescription       = "Updates a Function Instance based on given attributes"
	FlagName              = "The Function Instance's name"
	FlagIsActive          = "Whether the Function Instance is active or not"
	FlagInstanceID        = "Unique identifier of the Function Instance"
	FlagFile              = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OutputSuccess         = "Updated Function Instance with ID %d"
	HelpFlag              = "Displays more information about the update function-instance command"
	AskInputName          = "Enter the new Function Instance's name:"
	AskInputApplicationID = "Enter the Application's ID this Function Instance will be associated with:"
	AskInputInstanceID    = "Enter the Function Instance's ID:"
	AskInputFunctionID    = "Enter the Function's ID:"
	FlagApplicationID     = "Unique identifier of the Application"
	FlagFunctionID        = "Unique identifier of the Function"
	FlagArgs              = "The Function Instance's arguments"
)
