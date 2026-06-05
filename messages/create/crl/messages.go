package crl

var (
	Usage                  = "crl"
	CreateShortDescription = "Creates a Certificate Revocation List (CRL)"
	CreateLongDescription  = "Creates a Certificate Revocation List (CRL) by uploading a PEM-encoded CRL issued by a Certificate Authority"
	CreateOutputSuccess    = "Created Certificate Revocation List with ID %d"
	CreateFlagHelp         = "Displays more information about the create crl command"

	FlagName   = "Name that identifies the Certificate Revocation List"
	FlagIssuer = "Issuer of the Certificate Revocation List"
	FlagCRL    = "Path to the file containing the CRL content (PEM format)"
	FlagActive = "Whether the Certificate Revocation List is active or not"
	FlagIn     = "Path to a JSON file containing the attributes of the Certificate Revocation List being created; you can use - for reading from stdin"

	AskName   = "Enter the Certificate Revocation List's name:"
	AskIssuer = "Enter the Certificate Revocation List's issuer:"
	AskCRL    = "Enter the path to the CRL file (PEM format):"
)
