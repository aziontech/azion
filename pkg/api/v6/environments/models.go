package environments

import (
	"github.com/aziontech/azion-cli/pkg/api/v6/version"
	sdk "github.com/aziontech/azionapi-v6-go-sdk/azion-api"
)

// Deployment policies an Environment (and its Deployments) can have.
const (
	PolicySingleVersion = "single_version"
	PolicyVersionedURLs = "versioned_urls"
)

// Environment is what the deploy needs from one.
type Environment struct {
	ID               string
	Name             string
	DeploymentPolicy string
	VersionID        string
	VersionState     version.State
}

func environmentFrom(e *sdk.Environment) Environment {
	return Environment{
		ID: e.GetId(), Name: e.GetName(), DeploymentPolicy: e.GetDeploymentPolicy(),
		VersionID: e.GetVersionId(), VersionState: version.State(e.GetVersionState()),
	}
}
