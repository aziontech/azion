package csr

// Used only by the v4 command tree.
var (
	Usage            = "csr"
	ShortDescription = "Returns the Certificate Signing Request data"
	LongDescription  = "Displays information about a Certificate Signing Request via a given ID, including the generated CSR content"
	HelpFlag         = "Displays more information about the describe csr command"

	FlagId        = "Unique identifier of the Certificate Signing Request"
	AskInputCSRID = "Enter the Certificate Signing Request's ID:"
)
