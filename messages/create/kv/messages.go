package kv

// Used by both the v3 and the v4 command trees.
var (
	Usage               = "kv"
	ShortDescription    = "Creates a new namespace"
	LongDescription     = "Creates a namespace based on given value"
	FlagNamespace       = "The namespace's value"
	HelpFlag            = "Displays more information about the create kv command"
	AskNamespace        = "Enter the namespace:"
	CreateOutputSuccess = "Namespace '%s' created successfully"
)
