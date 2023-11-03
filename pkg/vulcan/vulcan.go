package vulcan

import "fmt"

const (
	versionVulcan        = "2.0.0"
	installEdgeFunctions = "npx --yes %s edge-functions@%s %s"
)

func Command(flags, params string) string {
	return fmt.Sprintf(installEdgeFunctions, flags, versionVulcan, params)
}
