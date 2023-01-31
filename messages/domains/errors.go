package domains

import "errors"

var (
	ErrorNameFlag            = errors.New("Failed to read the field name")
	ErrorCnames              = errors.New("Failed to read the field cnames")
	ErrorCnameAccesOnly      = errors.New("Failed to read cname field CnameAccessOnly")
	ErrorCnamesCertificateID = errors.New("Failed to read the field CnamesCertificateID")
	ErrorEdgeApplicationID   = errors.New("Failed to read the field EdgeApplicationId")
	IsActive                 = errors.New("Failed to read the field IsActive")

	ErrorMandatoryCreateFlags    = errors.New("A mandatory flag is missing. You must provide --name, --edge-application-id flags when the --in flag is not provided. Run the command 'azioncli edge_functions create --help' to display more information and try again")
	ErrorCodeFlag                = errors.New("Failed to read the code file. Verify if the file name and its path are correct and the file content has a valid code format. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorArgsFlag                = errors.New("Failed to read the args file. Verify if the file name and its path are correct and the file's content has a valid JSON format. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorParseArgs               = errors.New("Failed to parse JSON args. Verify if the file's content has a valid JSON format. Run the command 'azioncli edge_functions <subcommand> --help' to display more information and try again")
	ErrorMissingDomainIdArgument = errors.New("A mandatory flag is missing. You must provide a domain_id as an argument or path to import the file. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorCreateDomain            = errors.New("Failed to create Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorGetDomains              = errors.New("Failed to get the Domains: %s. Check your settings and try again. If the error persists, contact Azion support")

	//used by more than one cmd
	DomainsFlagId                      = "Unique identifier of the Domain"
	DomainsFileWritten                 = "File successfully written to: %s\n"
	ErrorMissingApplicationIdArgument  = errors.New("A mandatory flag is missing. You must provide a domain_id as an argument or path to import the file. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorGetDomain                     = errors.New("Failed to get the domain: %s. Check your settings and try again. If the error persists, contact Azion support")
	ErrorMissingDomainIdArgumentDelete = errors.New("A mandatory flag is missing. You must provide a domain_id as an argument. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorFailToDeleteDomain            = errors.New("Failed to delete the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")

  ErrorMissingCnames                 = errors.New("If CnameAccessOnly is true, you need to inform at least one Cname.")
	ErrorActiveFlag                    = errors.New("Invalid --active flag provided. The flag must have 'true' or 'false' values. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorDigitalCertificateFlag        = errors.New("Invalid --digital-certificate-id flag provided. The flag must have an Integer or 'null' as a value. Run the command 'azioncli domains <subcommand> --help' to display more information and try again")
	ErrorUpdateDomain                  = errors.New("Failed to update the Domain: %s. Check your settings and try again. If the error persists, contact Azion support")
)
