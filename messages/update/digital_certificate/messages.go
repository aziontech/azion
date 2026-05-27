package digitalcertificate

var (
	Usage                          = "digital-certificate"
	UpdateShortDescription         = "Updates a Digital Certificate"
	UpdateLongDescription          = "Modifies a Digital Certificate's attributes, such as its name, PEM contents, or managed-certificate settings, based on the given ID"
	FlagID                         = "The Digital Certificate's id"
	UpdateFlagName                 = "The Digital Certificate's name"
	UpdateFlagActive               = "Whether the Digital Certificate is active or not"
	UpdateFlagCertificate          = "Path to the file containing the certificate (PEM format)"
	UpdateFlagPrivateKey           = "Path to the file containing the private key (PEM format)"
	UpdateFlagCertificateType      = "The Digital Certificate's type (e.g. edge_certificate, trusted_ca_certificate)"
	UpdateFlagAuthority            = "Certificate Authority that issues the managed certificate (e.g. lets_encrypt)"
	UpdateFlagChallenge            = "Method used to solve the ACME challenge (dns or http)"
	UpdateFlagKeyAlgorithm         = "Key algorithm used to generate the certificate (rsa_2048, rsa_4096, or ecc_384)"
	UpdateFlagFile                 = "Given path and JSON file to automatically update the Digital Certificate attributes; you can use - for reading from stdin"
	UpdateOutputSuccess            = "Updated Digital Certificate with ID %d"
	UpdateHelpFlag                 = "Displays more information about the update digital-certificate command"
	UpdateAskDigitalCertificateID  = "Enter the ID of the Digital Certificate you wish to update:"
)
