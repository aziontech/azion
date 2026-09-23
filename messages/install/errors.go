package install

import "errors"

var (
	ErrorResolveHomeDir    = errors.New("failed to resolve home directory. Check your environment variables")
	ErrorSourceDirNotFound = errors.New("failed to locate bundled skills directory. CLI installation may be corrupted")
	ErrorCreateTargetDir   = errors.New("failed to create target directory: %s")
	ErrorRemoveExisting    = errors.New("failed to remove existing skill: %s")
	ErrorCopySkill         = errors.New("failed to copy skill '%s': %v")
	ErrorReadSourceDir     = errors.New("failed to read source skills directory")
)
