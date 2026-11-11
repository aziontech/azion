package build

import (
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
)

// The build reads azion.json for the preset and writes it back when it
// finishes. A v6 deploy writes azion.json in the v6 format, which the v4
// contract cannot read ("workloads" is a list there), so a build of a project
// deployed on v6 failed before it started, and writing the v4 contract back
// would have replaced what the deploy recorded.
//
// These pick the format by the file itself: a v4 file, or none, takes the v4
// path exactly as before, and a v6 file is read for its name and preset and
// written back with only the preset changed.

func getAzionJsonContent(confPath string) (*contracts.AzionApplicationOptions, error) {
	if !utils.IsAzionJsonV6(confPath) {
		return utils.GetAzionJsonContent(confPath)
	}
	conf, err := utils.GetAzionJsonContentV6(confPath)
	if err != nil {
		return nil, err
	}
	return &contracts.AzionApplicationOptions{Name: conf.Name, Preset: conf.Preset}, nil
}

func writeAzionJsonContent(in *contracts.AzionApplicationOptions, confPath string) error {
	if !utils.IsAzionJsonV6(confPath) {
		return utils.WriteAzionJsonContent(in, confPath)
	}
	conf, err := utils.GetAzionJsonContentV6(confPath)
	if err != nil {
		return err
	}
	conf.Preset = in.Preset
	return utils.WriteAzionJsonContentV6(conf, confPath)
}
