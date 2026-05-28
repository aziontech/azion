package crl

var (
	Usage                  = "crl"
	UpdateShortDescription = "Updates a Certificate Revocation List"
	UpdateLongDescription  = "Modifies a Certificate Revocation List's attributes, such as its name, issuer, or CRL content, based on the given ID"
	UpdateOutputSuccess    = "Updated Certificate Revocation List with ID %d"
	UpdateHelpFlag         = "Displays more information about the update crl command"

	FlagID           = "The Certificate Revocation List's id"
	UpdateFlagName   = "The Certificate Revocation List's name"
	UpdateFlagIssuer = "The Certificate Revocation List's issuer"
	UpdateFlagCRL    = "Path to the file containing the CRL content (PEM format)"
	UpdateFlagActive = "Whether the Certificate Revocation List is active or not"
	UpdateFlagFile   = "Given path and JSON file to automatically update the Certificate Revocation List attributes; you can use - for reading from stdin"

	UpdateAskCRLID = "Enter the ID of the Certificate Revocation List you wish to update:"
)
