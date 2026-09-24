package init

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorUnmarshalAzionFile           = errors.New("Failed to unmarshal the azion.json file. Verify if the file format is JSON or fix its content according to the JSON format specification at https://www.json.org/json-en.html")
	ErrorFailedCreatingAzionDirectory = errors.New("Failed to create the azion directory. The public's parent directory is read-only and/or isn't accessible. Change the permissions of the parent directory to read and write and/or give access to it")
	ErrorDeps                         = errors.New("Failed to install project dependencies")
	ErrorWorkingDir                   = errors.New("Failed to change current working directory")
	ErrorGetProjectInfo               = errors.New("Failed to get project preset")
)
