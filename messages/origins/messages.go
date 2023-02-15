package origins

var (
  // [ origins ]
	OriginsUsage                          = "origins"
	OriginsShortDescription               = "Origins is where data is fetched when the cache is not available."
  OriginsLongDescription                = "Origins is the original source of data in content delivery systems (CDN) where data is fetched when cache is not available. Data is stored at origin and can be retrieved by clients through cache servers distributed around the world. The fetch is done from the cache first, and if the data is not available, it is fetched from origin and saved in the cache for future use. This allows for fast data delivery."
	OriginsFlagHelp                       = "Displays more information about the origins command"

  // [ list ]
	OriginsListUsage                      = "list [flags]"
	OriginsListShortDescription           = "Displays yours origins"
	OriginsListLongDescription            = "Displays all your origins references to your edges"
	OriginsListHelpFlag                   = "Displays more information about the list subcommand"
	OriginsListFlagEdgeApplicationID      = "Is a unique identifier for the edge application that references the origins to direct data requests correctly."

	// [ describe ]
	OriginsDescribeUsage                  = "describe --application-id <domain_id> --origin-id <origin_id> [flags]"
	OriginsDescribeShortDescription       = "Returns the origin data"
	OriginsDescribeLongDescription        = "Displays information about the origin via a given ID to show the application’s attributes in detail"
	OriginsDescribeFlagApplicationID      = "Is a unique identifier for the edge application that references the origins to direct data requests correctly."
	OriginsDescribeFlagOriginID           = `is a unique identifier that identifies an "origins" in a list of results returned by the API. The "GetOrigin" function uses the "Origin Id" to search for the desired "origins" and returns an error if it is not found.`
	OriginsDescribeFlagOut                = "Exports the output to the given <file_path/file_name.ext>"
	OriginsDescribeFlagFormat             = "Changes the output format passing the json value to the flag"
	OriginsDescribeHelpFlag               = "Displays more information about the describe command"

	// [ create ]
	OriginsCreateUsage                    = "create [flags]"
	OriginsCreateShortDescription         = "Makes a new origin"
	OriginsCreateLongDescription          = "Makes a Origin based on given attributes to be used in Edge Applications"
	OriginsCreateFlagEdgeApplicationId    = "The Edge Application's unique identifier"
	OriginsCreateFlagName                 = "The Origin name"
  OriginsCreateFlagOriginType           = "Origin Type is a field used to identify the source of a record. It is typically used to differentiate between records that were created manually or automatically. For example, a record may have an Origin Type of 'Manual' if it was created by a user, or 'Automatic' if it was created by a system. Origin Type can also be used to differentiate between records that were imported from an external source, such as a CSV file, or created within the system."
  OriginsCreateFlagAddresses            = "Addresses linked to origins"
  OriginsCreateFlagOriginProtocolPolicy = "Is an origin protocol policy that specifies how Amazon CloudFront should respond to requests for content. This policy specifies whether CloudFront should use the origin protocol (HTTP or HTTPS) to get content from an origin server, or whether it should use the origin protocol regardless of the protocol used for the request."
  OriginsCreateFlagHostHeader           = "The HostHeader is an HTTP header field that specifies the hostname of the server being accessed. It is used to identify which website or application is being accessed. The HostHeader is sent by the browser to the web server and is used to determine which website or application to load."
  OriginsCreateFlagOriginPath           = "OriginPath is a file path that is used to identify the source of a file or directory. It is typically used to track the original location of a file or directory before it was moved or copied to a new location. OriginPath is often used in backup and restore operations to ensure that the original file or directory is not overwritten or lost."
  OriginsCreateFlagHmacAuthentication   = "HmacAuthentication is a type of authentication that uses a cryptographic hash function to verify the integrity of a message. It is a form of message authentication code (MAC) that uses a shared secret key between two parties to generate a signature for a message. The signature is then used to verify that the message has not been tampered with or altered in any way. HmacAuthentication is used to ensure that the message is authentic and has not been modified in transit. It is commonly used in web applications and APIs to secure data transmission"
  OriginsCreateFlagHmacRegionName       = "HmacRegionName is a field used in the Amazon Web Services (AWS) API to identify the region in which a particular request is being made. It is used to ensure that requests are routed to the correct region and that the correct authentication credentials are used."
  OriginsCreateFlagHmacAccessKey        = "HmacAccessKey is a type of authentication key used to access an API or other secure system. It is a combination of a secret key and a cryptographic hash algorithm, such as SHA-256, to generate a unique signature for each request. The signature is then used to verify the authenticity of the request and ensure that it has not been tampered with. HmacAccessKey is often used in combination with other authentication methods, such as OAuth or API keys, to provide an additional layer of security."
  OriginsCreateFlagHmacSecretKey        = "HmacSecretKey is a type of cryptographic key used in the HMAC (Hash-based Message Authentication Code) algorithm. It is a secret key that is used to generate a cryptographic hash of a message, which is then used to verify the authenticity and integrity of the message. The key is typically a string of random characters that is known only to the sender and receiver of the message."
	OriginsCreateFlagIn                   = " Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OriginsCreateOutputSuccess            = "Created origin with ID %d\n"
	OriginsCreateHelpFlag                 = "Displays more information about the create subcommand"

  // [ update ]
	OriginsUpdateUsage                    = "update [flags]"
	OriginsUpdateShortDescription         = "Modifies an Origin"
	OriginsUpdateLongDescription          = "Modifies an Origin based on its ID to update its name, activity status, and other attributes"
  OriginsUpdateFlagOriginKey             = "The Origins unique identifier"
	OriginsUpdateFlagEdgeApplicationId    = "The Edge Application's unique identifier"
	OriginsUpdateFlagName                 = "The Origin name"
  OriginsUpdateFlagOriginType           = "Origin Type is a field used to identify the source of a record. It is typically used to differentiate between records that were created manually or automatically. For example, a record may have an Origin Type of 'Manual' if it was created by a user, or 'Automatic' if it was created by a system. Origin Type can also be used to differentiate between records that were imported from an external source, such as a CSV file, or created within the system."
  OriginsUpdateFlagAddresses            = "Addresses linked to origins"
  OriginsUpdateFlagOriginProtocolPolicy = "Is an origin protocol policy that specifies how Amazon CloudFront should respond to requests for content. This policy specifies whether CloudFront should use the origin protocol (HTTP or HTTPS) to get content from an origin server, or whether it should use the origin protocol regardless of the protocol used for the request."
  OriginsUpdateFlagHostHeader           = "The HostHeader is an HTTP header field that specifies the hostname of the server being accessed. It is used to identify which website or application is being accessed. The HostHeader is sent by the browser to the web server and is used to determine which website or application to load."
  OriginsUpdateFlagOriginPath           = "OriginPath is a file path that is used to identify the source of a file or directory. It is typically used to track the original location of a file or directory before it was moved or copied to a new location. OriginPath is often used in backup and restore operations to ensure that the original file or directory is not overwritten or lost."
  OriginsUpdateFlagHmacAuthentication   = "HmacAuthentication is a type of authentication that uses a cryptographic hash function to verify the integrity of a message. It is a form of message authentication code (MAC) that uses a shared secret key between two parties to generate a signature for a message. The signature is then used to verify that the message has not been tampered with or altered in any way. HmacAuthentication is used to ensure that the message is authentic and has not been modified in transit. It is commonly used in web applications and APIs to secure data transmission"
  OriginsUpdateFlagHmacRegionName       = "HmacRegionName is a field used in the Amazon Web Services (AWS) API to identify the region in which a particular request is being made. It is used to ensure that requests are routed to the correct region and that the correct authentication credentials are used."
  OriginsUpdateFlagHmacAccessKey        = "HmacAccessKey is a type of authentication key used to access an API or other secure system. It is a combination of a secret key and a cryptographic hash algorithm, such as SHA-256, to generate a unique signature for each request. The signature is then used to verify the authenticity of the request and ensure that it has not been tampered with. HmacAccessKey is often used in combination with other authentication methods, such as OAuth or API keys, to provide an additional layer of security."
  OriginsUpdateFlagHmacSecretKey        = "HmacSecretKey is a type of cryptographic key used in the HMAC (Hash-based Message Authentication Code) algorithm. It is a secret key that is used to generate a cryptographic hash of a message, which is then used to verify the authenticity and integrity of the message. The key is typically a string of random characters that is known only to the sender and receiver of the message."
	OriginsUpdateFlagIn                   = "Path to a JSON file containing the attributes that will be updated; you can use - for reading from stdin"
	OriginsUpdateOutputSuccess            = "Update origin with ID %s\n"
	OriginsUpdateHelpFlag                 = "Displays more information about the create subcommand"

	// [ delete ] 
	OriginsDeleteUsage                    = "delete [flags]"
	OriginsDeleteShortDescription         = "Removes an Origin"
	OriginsDeleteLongDescription          = "Removes an Origin from the Edge Applications library based on its given ID"
  OriginsDeleteOutputSuccess            = "Origin %s was successfully deleted\n"
  OriginsDeleteFlagApplicationID        = "The Edge Application's unique identifier"
  OriginsDeleteFlagOriginKey            = "The Origin key unique identifier"
  OriginsDeleteHelpFlag                 = "Displays more information about the delete subcommand"

  // [ general ]
	OriginsFileWritten                    = "File successfully written to: %s\n"
)
