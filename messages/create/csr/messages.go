package csr

var (
	Usage                  = "csr"
	CreateShortDescription = "Creates a Certificate Signing Request (CSR)"
	CreateLongDescription  = "Creates a Certificate Signing Request (CSR) to be submitted to a Certificate Authority. The generated CSR content is returned upon creation"
	CreateOutputSuccess    = "Created Certificate Signing Request with ID %d"
	CreateFlagHelp         = "Displays more information about the create csr command"

	FlagName              = "Name that identifies the Certificate Signing Request"
	FlagCommonName        = "Common Name (CN) of the certificate subject"
	FlagCountry           = "Country code of the certificate subject (e.g. US, BR)"
	FlagState             = "State or province name of the certificate subject"
	FlagLocality          = "City or locality name of the certificate subject"
	FlagOrganization      = "Organization name of the certificate subject"
	FlagOrganizationUnity = "Organizational unit name of the certificate subject"
	FlagEmail             = "Contact email address of the certificate subject"
	FlagAlternativeNames  = "Comma-separated list of Subject Alternative Names (SANs)"
	FlagCertificateType   = "Certificate type (e.g. edge_certificate, trusted_ca_certificate)"
	FlagKeyAlgorithm      = "Key algorithm used to generate the certificate (rsa_2048, rsa_4096, or ecc_384)"
	FlagIn                = "Path to a JSON file containing the attributes of the Certificate Signing Request being created; you can use - for reading from stdin"

	AskName              = "Enter the Certificate Signing Request's name:"
	AskCommonName        = "Enter the Common Name (CN) of the certificate subject:"
	AskCountry           = "Enter the country code of the certificate subject (e.g. US, BR):"
	AskState             = "Enter the state or province name of the certificate subject:"
	AskLocality          = "Enter the city or locality name of the certificate subject:"
	AskOrganization      = "Enter the organization name of the certificate subject:"
	AskOrganizationUnity = "Enter the organizational unit name of the certificate subject:"
	AskEmail             = "Enter the contact email address of the certificate subject:"
)
