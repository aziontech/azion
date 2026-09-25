package link

import "errors"

// Used by both the v3 and the v4 command trees.
var (
	ErrorDeps             = errors.New("Failed to install project dependencies")
	ErrorReadingGitignore = errors.New("Failed to read your .gitignore file")
	ErrorWritingGitignore = errors.New("Failed to write to your .gitignore file")
	ErrorReadingWorkflow  = errors.New("Failed to check for GitHub Actions workflow file")
	ErrorWritingWorkflow  = errors.New("Failed to write GitHub Actions workflow file")
)
