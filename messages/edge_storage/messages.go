package edge_storage

const (
	USAGE                           = "edge-storage"
	USAGE_BUCKET                    = "bucket"
	SHORT_DESCRIPTION               = "manage storage api Buckets and Objects directly through the command line"
	SHORT_DESCRIPTION_CREATE_BUCKET = "manage storage api Buckets and Objects directly through the command line"
	LONG_DESCRIPTION                = "allows users to perform a wide range of operations, from creating and deleting Buckets to adding, removing and manipulating Objects, all quickly and conveniently, without the need for graphical interfaces or complex manual interactions."
	LONG_DESCRIPTION_CREATE_BUCKET  = "allows users to perform a wide range of operations, from creating and deleting Buckets to adding, removing and manipulating Objects, all quickly and conveniently, without the need for graphical interfaces or complex manual interactions."
	EXAMPLE                         = "$ azion create edge-storage\n$ azion create edge-storage --help"
	EXAMPLE_CREATE_BUCKET           = "$ azion create edge-storage bucket --name 'zorosola' --edge-access 'read_only'"
	FLAG_HELP                       = "Displays more information about the edge-storage command"
	FLAG_HELP_CREATE_BUCKET         = "Displays more information about the create edge-storege bucket command"
	FLAG_NAME_CREATE_BUCKET         = "The name of the Edge Storage Bucket"
	FLAG_EDGE_ACCESS_CREATE_BUCKET  = "Edge access is the access level of the bucket"
	SUCCESS_CREATE_BUCKET           = "Created Bucket"
	ASK_NAME_CREATE_BUCKET          = "Enter your name Bucket"
	ASK_EDGE_ACCESSS_CREATE_BUCKET  = "Enter your Edge Access Bucket"
)
