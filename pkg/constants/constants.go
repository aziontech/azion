package constants

// These variables have their values set at link time through the Makefile
var (
	ApiURL        string
	AuthURL       string
	StorageApiURL string
)

const (
	PathAzionJson        = "/azion/azion.json"
	NpxVulcan     string = "npx --yes edge-functions@1.7.0 "
)
