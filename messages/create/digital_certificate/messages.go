package digitalcertificate

var (
	Usage                  = "digital-certificate"
	CreateShortDescription = "Creates a Digital Certificate"
	CreateLongDescription  = "Creates a Digital Certificate by either uploading your own PEM-encoded certificate and private key, or by requesting a managed certificate (e.g. Let's Encrypt) issued by a Certificate Authority"
	FlagIn                 = "Path to a JSON file containing the attributes of the Digital Certificate being created; you can use - for reading from stdin"
	CreateFlagHelp         = "Displays more information about the create digital-certificate command"
	CreateOutputSuccess    = "Created Digital Certificate with ID %d"
	RequestOutputSuccess   = "Requested Digital Certificate with ID %d"

	FlagName            = "Digital Certificate's name"
	FlagCertificate     = "Path to the file containing the certificate (PEM format)"
	FlagPrivateKey      = "Path to the file containing the private key (PEM format)"
	FlagCertificateType = "Digital Certificate's type (e.g. edge_certificate, trusted_ca_certificate)"

	FlagAuthority        = "Certificate Authority that will issue the certificate (e.g. lets_encrypt). When provided, the certificate is requested instead of uploaded"
	FlagChallenge        = "Method used to solve the ACME challenge (dns or http)"
	FlagCommonName       = "Common Name (CN) of the certificate being requested"
	FlagAlternativeNames = "Comma-separated list of Subject Alternative Names (SANs) for the requested certificate"
	FlagKeyAlgorithm     = "Key algorithm used to generate the certificate (rsa_2048, rsa_4096, or ecc_384)"

	AskName            = "Enter the Digital Certificate's name:"
	AskCertificate     = "Enter the path to the certificate file:"
	AskPrivateKey      = "Enter the path to the private key file:"
	AskCertificateType = "Enter the Digital Certificate's type:"
	AskChallenge       = "Enter the ACME challenge method (dns or http):"
	AskCommonName      = "Enter the Common Name (CN) of the certificate:"
)
