package version

import "context"

// Versioned is the contract the seven resources with a build share
// (applications, connectors, custom pages, firewalls, functions, network lists
// and WAFs). It only describes what their clients in pkg/api/v6 already do, so
// the build steps drive any of them without repeating the code per resource.
type Versioned interface {
	// CreateVersion opens a new Draft. An empty from clones the latest ready
	// version, which is the API's default. It returns the new version id.
	CreateVersion(ctx context.Context, id int64, from, comment string) (string, error)
	// Build queues the build of a Draft. It returns as soon as the API accepts
	// the request; VersionState tells when the build settles.
	Build(ctx context.Context, id int64, versionID, comment string) error
	// CancelBuild asks the API to stop a build that has not settled yet.
	CancelBuild(ctx context.Context, id int64, versionID string) error
	// VersionState reads the build state of one version.
	VersionState(ctx context.Context, id int64, versionID string) (State, error)
}
